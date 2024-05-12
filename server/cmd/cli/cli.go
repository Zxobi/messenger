package main

import (
	"bufio"
	"github.com/dvid-messanger/cmd/cli/builder"
	"github.com/dvid-messanger/cmd/cli/builder/cmdbuilders"
	"github.com/dvid-messanger/cmd/cli/printer"
	"github.com/dvid-messanger/cmd/cli/printer/formatter"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = "localhost:20203"

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	upstreamBuilder := builder.NewBuilder()
	cmdbuilders.AddAuthBuilders(upstreamBuilder)
	cmdbuilders.AddUserBuilders(upstreamBuilder)
	cmdbuilders.AddChatBuilders(upstreamBuilder)
	cmdbuilders.AddInfoInitBuilders(upstreamBuilder)
	cmdbuilders.AddSystemBuilders(upstreamBuilder)

	downstreamPrinter := printer.NewPrinter()
	formatter.AddAuthFormatters(downstreamPrinter)
	formatter.AddUserFormatters(downstreamPrinter)
	formatter.AddChatFormatters(downstreamPrinter)
	formatter.AddInfoFormatters(downstreamPrinter)
	formatter.AddSystemFormatters(downstreamPrinter)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial_err:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read_err: ", err)
				return
			}
			downstreamPrinter.Print(message)
		}
	}()

	input := make(chan []byte)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			cmd := upstreamBuilder.Build(scanner.Text())
			if cmd == nil {
				continue
			}

			input <- cmd
		}
	}()

	for {
		select {
		case <-done:
			return
		case msg := <-input:
			c.WriteMessage(websocket.BinaryMessage, msg)
		case <-interrupt:
			log.Println("interrupt")

			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("close_err:", err)
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
