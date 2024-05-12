package cmdbuilders

import (
	"encoding/base64"
	"fmt"
	"github.com/dvid-messanger/cmd/cli/builder"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddUserBuilders(b *builder.Builder) {
	b.AddBuilder("cur", frontendv1.UpstreamType_U_CUR_USER, BuildCurUser)
	b.AddBuilder("user", frontendv1.UpstreamType_U_GET_USER, BuildGetUser)
	b.AddBuilder("users", frontendv1.UpstreamType_U_GET_USERS, BuildGetUsers)
	b.AddBuilder("reg", frontendv1.UpstreamType_U_REG_USER, BuildRegUser)
}

func BuildCurUser(_ []string) proto.Message {
	return &frontendv1.UpstreamCurUser{}
}

func BuildGetUser(args []string) proto.Message {
	if len(args) < 1 {
		fmt.Println("usage: reg [uid]")
		return nil
	}
	uid, err := base64.StdEncoding.DecodeString(args[0])
	if err != nil {
		fmt.Println("bad uid")
		return nil
	}

	return &frontendv1.UpstreamGetUser{
		Uid: uid,
	}
}

func BuildGetUsers(_ []string) proto.Message {
	return &frontendv1.UpstreamGetUsers{}
}

func BuildRegUser(args []string) proto.Message {
	if len(args) < 2 {
		fmt.Println("usage: reg [email] [password]")
	}
	return &frontendv1.UpstreamRegUser{
		Email:    args[0],
		Password: args[1],
	}
}
