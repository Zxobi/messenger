package builder

import (
	"fmt"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
	"strings"
)

type UpstreamBuilderFunc = func(args []string) proto.Message

type UpstreamBuilder struct {
	ut  frontendv1.UpstreamType
	fun UpstreamBuilderFunc
}

type Builder struct {
	upstreamBuilders map[string]*UpstreamBuilder
}

func NewBuilder() *Builder {
	return &Builder{
		upstreamBuilders: make(map[string]*UpstreamBuilder),
	}
}

func (b *Builder) AddBuilder(upstream string, ut frontendv1.UpstreamType, fun UpstreamBuilderFunc) {
	b.upstreamBuilders[upstream] = &UpstreamBuilder{ut: ut, fun: fun}
}

func (b *Builder) Build(input string) []byte {
	args := strings.Fields(input)

	if len(args) < 1 {
		return nil
	}

	builder := b.upstreamBuilders[args[0]]
	if builder == nil {
		fmt.Println("unknown command")
		return nil
	}

	payload := builder.fun(args[1:])
	if payload == nil {
		fmt.Println("failed to make command")
		return nil
	}

	mPayload, err := proto.Marshal(payload)
	if err != nil {
		fmt.Println("failed to marshal payload " + err.Error())
		return nil
	}

	res, err := proto.Marshal(&frontendv1.Upstream{Type: builder.ut, Payload: mPayload})
	if err != nil {
		fmt.Println("failed to marshal upstream " + err.Error())
		return nil
	}

	return res
}
