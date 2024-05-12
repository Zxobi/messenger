package cmdbuilders

import (
	"github.com/dvid-messanger/cmd/cli/builder"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

func AddInfoInitBuilders(b *builder.Builder) {
	b.AddBuilder("init", frontendv1.UpstreamType_U_INFO_INIT, BuildInfoInit)
}

func BuildInfoInit(_ []string) proto.Message {
	return &frontendv1.UpstreamInfoInit{}
}
