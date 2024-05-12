package cmdbuilders

import (
	"github.com/dvid-messanger/cmd/cli/builder"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddAuthBuilders(b *builder.Builder) {
	b.AddBuilder("login", frontendv1.UpstreamType_U_LOGIN, BuildLogin)
	b.AddBuilder("logout", frontendv1.UpstreamType_U_LOGOUT, BuildLogout)
}

func BuildLogin(args []string) proto.Message {
	return &frontendv1.UpstreamLogin{
		Email:    args[0],
		Password: args[1],
	}
}

func BuildLogout(_ []string) proto.Message {
	return &frontendv1.UpstreamLogout{}
}
