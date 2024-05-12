package cmdbuilders

import (
	"fmt"
	"github.com/dvid-messanger/cmd/cli/builder"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddSystemBuilders(b *builder.Builder) {
	b.AddBuilder("echo", frontendv1.UpstreamType_U_ECHO, BuildEcho)
}

func BuildEcho(args []string) proto.Message {
	if len(args) < 1 {
		fmt.Println("usage: echo [content]")
		return nil
	}
	return &frontendv1.UpstreamEcho{
		Content: args[1],
	}
}
