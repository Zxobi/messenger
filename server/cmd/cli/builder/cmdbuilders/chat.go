package cmdbuilders

import (
	"encoding/base64"
	"fmt"
	"github.com/dvid-messanger/cmd/cli/builder"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddChatBuilders(b *builder.Builder) {
	b.AddBuilder("chats", frontendv1.UpstreamType_U_GET_USER_CHATS, BuildGetUserChats)
	b.AddBuilder("cchat", frontendv1.UpstreamType_U_CREATE_CHAT, BuildCreateChat)
	b.AddBuilder("msgs", frontendv1.UpstreamType_U_CHAT_MESSAGES, BuildChatMessages)
	b.AddBuilder("msg", frontendv1.UpstreamType_U_SEND_MESSAGE, BuildSendMessage)
}

func BuildGetUserChats(_ []string) proto.Message {
	return &frontendv1.UpstreamGetUserChats{}
}

func BuildChatMessages(args []string) proto.Message {
	if len(args) < 1 {
		fmt.Println("usage: msgs [cid]")
		return nil
	}
	cid, err := base64.StdEncoding.DecodeString(args[0])
	if err != nil {
		fmt.Println("bad cid")
		return nil
	}

	return &frontendv1.UpstreamChatMessages{
		Cid: cid,
	}
}

func BuildCreateChat(args []string) proto.Message {
	if len(args) < 1 {
		fmt.Println("usage: cchat [uid]")
		return nil
	}
	uid, err := base64.StdEncoding.DecodeString(args[0])
	if err != nil {
		fmt.Println("bad uid")
		return nil
	}

	return &frontendv1.UpstreamCreateChat{
		Uid: uid,
	}
}

func BuildSendMessage(args []string) proto.Message {
	if len(args) < 2 {
		fmt.Println("usage: msg [cid] [text]")
		return nil
	}
	cid, err := base64.StdEncoding.DecodeString(args[0])
	if err != nil {
		fmt.Println("bad cid")
		return nil
	}

	return &frontendv1.UpstreamSendMessage{
		Cid:  cid,
		Text: args[1],
	}
}
