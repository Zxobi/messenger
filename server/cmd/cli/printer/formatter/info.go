package formatter

import (
	"github.com/dvid-messanger/cmd/cli/printer"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddInfoFormatters(printer *printer.Printer) {
	printer.AddFormatter(frontendv1.DownstreamType_D_INFO_INIT, &frontendv1.DownstreamInfoInit{}, FormatInfoInit)
}

func FormatInfoInit(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamInfoInit)

	chats := FormatChats(downstream.GetChats())
	user := FormatUser(downstream.GetUser())

	return "\tuser:\n" + user + "\n\tchats:[\n" + chats + "\n\t]"
}
