package formatter

import (
	"encoding/base64"
	"github.com/dvid-messanger/cmd/cli/printer"
	"github.com/dvid-messanger/pkg/cutil"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	protocolv1 "github.com/dvid-messanger/protos/gen/protocol"
	"github.com/golang/protobuf/proto"
	"strings"
)

func AddChatFormatters(printer *printer.Printer) {
	printer.AddFormatter(frontendv1.DownstreamType_D_GET_CHAT, &frontendv1.DownstreamGetChat{}, FormatGetChat)
	printer.AddFormatter(frontendv1.DownstreamType_D_GET_USER_CHATS, &frontendv1.DownstreamGetUserChats{}, FormatGetUserChats)
	printer.AddFormatter(frontendv1.DownstreamType_D_CREATE_CHAT, &frontendv1.DownstreamCreateChat{}, FormatCreateChat)

	printer.AddFormatter(frontendv1.DownstreamType_D_CHAT_MESSAGES, &frontendv1.DownstreamChatMessages{}, FormatChatMessages)

	printer.AddFormatter(frontendv1.DownstreamType_D_SEND_MESSAGE, &frontendv1.DownstreamSendMessage{}, FormatSendMessage)
	printer.AddFormatter(frontendv1.DownstreamType_D_NEW_MESSAGE, &frontendv1.DownstreamNewMessage{}, FormatNewMessage)

}

func FormatGetChat(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamGetChat)

	return FormatChat(downstream.GetChat())
}

func FormatGetUserChats(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamGetUserChats)

	return FormatChats(downstream.GetChats())
}

func FormatCreateChat(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamCreateChat)

	return FormatChat(downstream.GetChat())
}

func FormatChatMessages(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamChatMessages)

	return FormatMessages(downstream.GetMessages())
}

func FormatSendMessage(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamSendMessage)

	return FormatMessage(downstream.GetMessage())
}

func FormatNewMessage(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamNewMessage)

	return FormatMessage(downstream.GetMessage())
}

func FormatChat(chat *protocolv1.Chat) string {
	members := strings.Join(cutil.Map(chat.GetChatMembers(), func(member *protocolv1.ChatMember) string {
		return "\t\t\t{ id=" + base64.StdEncoding.EncodeToString(member.GetUid()) + " }"
	}), ",\n")
	return "\t{ id=" + base64.StdEncoding.EncodeToString(chat.GetId()) + ", type=" + chat.GetType().String() +
		", members=[\n" + members + "\n\t]}"
}

func FormatChats(chats []*protocolv1.Chat) string {
	return strings.Join(cutil.Map(chats, func(chat *protocolv1.Chat) string {
		return FormatChat(chat)
	}), ",\n")
}

func FormatMessage(msg *protocolv1.ChatMessage) string {
	return "\t{ id=" + base64.StdEncoding.EncodeToString(msg.GetId()) +
		", from=" + base64.StdEncoding.EncodeToString(msg.GetId()) +
		", text=\"" + msg.Text + "\" }"
}

func FormatMessages(messages []*protocolv1.ChatMessage) string {
	return strings.Join(cutil.Map(messages, func(msg *protocolv1.ChatMessage) string {
		return FormatMessage(msg)
	}), ",\n")
}
