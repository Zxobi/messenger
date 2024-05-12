package formatter

import (
	"encoding/base64"
	"fmt"
	"github.com/dvid-messanger/cmd/cli/printer"
	"github.com/dvid-messanger/pkg/cutil"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	protocolv1 "github.com/dvid-messanger/protos/gen/protocol"
	"github.com/golang/protobuf/proto"
	"strings"
)

func AddUserFormatters(printer *printer.Printer) {
	printer.AddFormatter(frontendv1.DownstreamType_D_GET_USER, &frontendv1.DownstreamGetUser{}, FormatGetUser)
	printer.AddFormatter(frontendv1.DownstreamType_D_GET_USERS, &frontendv1.DownstreamGetUsers{}, FormatGetUsers)
	printer.AddFormatter(frontendv1.DownstreamType_D_CUR_USER, &frontendv1.DownstreamCurUser{}, FormatCurUser)

	printer.AddFormatter(frontendv1.DownstreamType_D_REG_USER, &frontendv1.DownstreamRegUser{}, FormatRegUser)
}

func FormatGetUser(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamGetUser)

	return FormatUser(downstream.GetUser())
}

func FormatGetUsers(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamGetUsers)

	return FormatUsers(downstream.GetUsers())
}

func FormatCurUser(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamCurUser)

	return FormatUser(downstream.GetUser())
}

func FormatRegUser(payload proto.Message) string {
	downstream := payload.(*frontendv1.DownstreamRegUser)

	return FormatUser(downstream.GetUser())
}

func FormatUser(user *protocolv1.User) string {
	return fmt.Sprintf("\t{ id=%s, email=%s, bio=%s}", base64.StdEncoding.EncodeToString(user.Id), user.Email, user.Bio)
}

func FormatUsers(users []*protocolv1.User) string {
	return strings.Join(cutil.Map(users, func(user *protocolv1.User) string {
		return FormatUser(user)
	}), ",\n")
}
