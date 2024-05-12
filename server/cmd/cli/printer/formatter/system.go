package formatter

import (
	"github.com/dvid-messanger/cmd/cli/printer"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddSystemFormatters(printer *printer.Printer) {
	printer.AddFormatter(frontendv1.DownstreamType_D_ECHO, &frontendv1.DownstreamEcho{}, FormatEcho)
}

func FormatEcho(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamEcho)

	return "\t{ content=\"" + downstream.GetContent() + "\" }"
}
