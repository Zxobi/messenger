package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/dvid-messanger/internal/lib/cutils"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	protocolv1 "github.com/dvid-messanger/protos/gen/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var addr = "localhost:20206"

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			printDownstream(message)
		}
	}()

	input := make(chan []byte)
	go readInput(input)

	for {
		select {
		case <-done:
			return
		case msg := <-input:
			c.WriteMessage(websocket.BinaryMessage, msg)
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func readInput(c chan<- []byte) {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()

		input = scanner.Text()
		fields := strings.Fields(input)

		if len(fields) < 1 {
			continue
		}

		switch fields[0] {
		case "echo":
			if !requireArgs("echo", fields[1:], 1) {
				continue
			}
			payload, _ := proto.Marshal(&frontendv1.UpstreamEcho{Content: fields[1]})
			msg, _ := proto.Marshal(&frontendv1.Upstream{
				Type:    frontendv1.UpstreamType_U_ECHO,
				Payload: payload,
			})
			c <- msg
		case "login":
			if !requireArgs("login", fields[1:], 2) {
				continue
			}
			payload, _ := proto.Marshal(&frontendv1.UpstreamLogin{
				Email:    fields[1],
				Password: fields[2],
			})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_LOGIN)
		case "reg":
			if !requireArgs("reg", fields[1:], 2) {
				continue
			}
			payload, _ := proto.Marshal(&frontendv1.UpstreamRegUser{
				Email:    fields[1],
				Password: fields[2],
			})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_REG_USER)
		case "cur":
			payload, _ := proto.Marshal(&frontendv1.UpstreamCurUser{})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_CUR_USER)
		case "users":
			payload, _ := proto.Marshal(&frontendv1.UpstreamGetUsers{})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_GET_USERS)
		case "chats":
			payload, _ := proto.Marshal(&frontendv1.UpstreamGetUserChats{})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_GET_USER_CHATS)
		case "newchat":
			if !requireArgs("newchat", fields[1:], 1) {
				continue
			}
			uid, err := base64.StdEncoding.DecodeString(fields[1])
			if err != nil {
				fmt.Println("newchat: error: " + err.Error())
				continue
			}
			payload, _ := proto.Marshal(&frontendv1.UpstreamCreateChat{Uid: uid})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_CREATE_CHAT)
		case "msg":
			if !requireArgs("msg", fields[1:], 2) {
				continue
			}
			cid, err := base64.StdEncoding.DecodeString(fields[1])
			if err != nil {
				fmt.Println("msg: error: " + err.Error())
				continue
			}
			payload, _ := proto.Marshal(&frontendv1.UpstreamSendMessage{Cid: cid, Text: strings.Join(fields[2:], " ")})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_SEND_MESSAGE)
		case "msgs":
			if !requireArgs("msgs", fields[1:], 1) {
				continue
			}
			cid, err := base64.StdEncoding.DecodeString(fields[1])
			if err != nil {
				fmt.Println("msgs: error: " + err.Error())
				continue
			}
			payload, _ := proto.Marshal(&frontendv1.UpstreamChatMessages{Cid: cid})
			c <- mustMakeUpstream(payload, frontendv1.UpstreamType_U_CHAT_MESSAGES)

		default:
			fmt.Println("unknown command")
		}

	}
}

func requireArgs(cmd string, args []string, rLen int) bool {
	if len(args) < rLen {
		fmt.Println(cmd + ": require " + strconv.Itoa(rLen) + " args")
		return false
	}

	return true
}

func mustMakeUpstream(payload []byte, upstreamType frontendv1.UpstreamType) []byte {
	msg, _ := proto.Marshal(&frontendv1.Upstream{
		Type:    upstreamType,
		Payload: payload,
	})

	return msg
}

func printDownstream(b []byte) {
	downstream := &frontendv1.Downstream{}
	err := proto.Unmarshal(b, downstream)
	if err != nil {
		panic(err)
	}

	if downstream.GetError() != nil {
		fmt.Println("recv: error: " + downstream.Type.String() + " - " + downstream.Error.GetCode().String() + " " + downstream.Error.GetDesc())
		return
	}

	switch downstream.Type {
	case frontendv1.DownstreamType_D_LOGIN:
		payload := &frontendv1.DownstreamLogin{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("login: " + payload.String())
	case frontendv1.DownstreamType_D_NEW_MESSAGE:
		payload := &frontendv1.DownstreamNewMessage{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("new_msg: " + msgString(payload.GetMessage()))
	case frontendv1.DownstreamType_D_REG_USER:
		payload := &frontendv1.DownstreamRegUser{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("reg: " + userString(payload.GetUser()))
	case frontendv1.DownstreamType_D_SEND_MESSAGE:
		payload := &frontendv1.DownstreamSendMessage{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("send_msg: " + msgString(payload.GetMessage()))
	case frontendv1.DownstreamType_D_GET_USERS:
		payload := &frontendv1.DownstreamGetUsers{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("users: " + strings.Join(cutils.Map(payload.GetUsers(), func(user *protocolv1.User) string {
			return userString(user)
		}), " "))
	case frontendv1.DownstreamType_D_CUR_USER:
		payload := &frontendv1.DownstreamCurUser{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("cur: " + userString(payload.GetUser()))
	case frontendv1.DownstreamType_D_GET_USER_CHATS:
		payload := &frontendv1.DownstreamGetUserChats{}
		proto.Unmarshal(downstream.Payload, payload)
		fmt.Println("chats: " + strings.Join(cutils.Map(payload.Chats, func(chat *protocolv1.Chat) string {
			return base64.StdEncoding.EncodeToString(chat.GetId())
		}), " "))

	default:
		fmt.Println("recv: " + downstream.Type.String() + " payload " + base64.StdEncoding.EncodeToString(downstream.Payload))
	}
}

func msgString(msg *protocolv1.ChatMessage) string {
	mid := base64.StdEncoding.EncodeToString(msg.GetId())
	cid := base64.StdEncoding.EncodeToString(msg.GetCid())
	uid := base64.StdEncoding.EncodeToString(msg.GetUid())
	sent := time.UnixMilli(msg.GetTimestamp()).Format(time.UnixDate)
	return fmt.Sprintf("msg(mid %s cid %s uid %s time \"%s\" text \"%s\")", mid, cid, uid, sent, msg.GetText())
}

func userString(user *protocolv1.User) string {
	uid := base64.StdEncoding.EncodeToString(user.GetId())
	return fmt.Sprintf("user(uid %s email %s)", uid, user.GetEmail())
}
