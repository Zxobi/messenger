package formatter

import (
	"github.com/dvid-messanger/cmd/cli/printer"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddAuthFormatters(printer *printer.Printer) {
	printer.AddFormatter(frontendv1.DownstreamType_D_LOGIN, &frontendv1.DownstreamLogin{}, FormatLogin)
	printer.AddFormatter(frontendv1.DownstreamType_D_LOGOUT, &frontendv1.DownstreamLogout{}, FormatLogout)
}

func FormatLogin(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamLogin)

	return "\t{ token=" + downstream.GetToken() + " }"
}

func FormatLogout(_ proto.Message) string {
	return "\tlogged out"
}
