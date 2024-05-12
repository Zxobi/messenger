package printer

import (
	"fmt"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
	"reflect"
)

const prefix = "recv: "

type FormatterFunc = func(downstream proto.Message) string

type Formatter struct {
	msg proto.Message
	fun FormatterFunc
}

type Printer struct {
	formatters map[frontendv1.DownstreamType]*Formatter
}

func NewPrinter() *Printer {
	return &Printer{
		formatters: make(map[frontendv1.DownstreamType]*Formatter),
	}
}

func (p *Printer) AddFormatter(dt frontendv1.DownstreamType, msg proto.Message, fun FormatterFunc) {
	p.formatters[dt] = &Formatter{msg: msg, fun: fun}
}

func (p *Printer) Print(data []byte) {
	downstream := &frontendv1.Downstream{}
	if err := proto.Unmarshal(data, downstream); err != nil {
		fmt.Println(prefix + "error during unmarshalling " + err.Error())
	}

	if downstream.GetError() != nil && downstream.Error.GetCode() != frontendv1.ErrorCode_NO_ERROR {
		fmt.Println(prefix + formatError(downstream))
		return
	}

	formatter := p.formatters[downstream.Type]
	if formatter == nil {
		fmt.Println(prefix + downstream.GetType().String())
		return
	}

	payload := reflect.New(reflect.TypeOf(formatter.msg).Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(downstream.Payload, payload); err != nil {
		fmt.Println(prefix + downstream.GetType().String() + " error during payload unmarshalling " + err.Error())
	}

	fmt.Println(prefix + downstream.GetType().String() + " {\n" + formatter.fun(payload) + "\n}")
}

func formatError(downstream *frontendv1.Downstream) string {
	return fmt.Sprintf("error: code=%d desc=\"%s\"", downstream.GetError().Code, downstream.GetError().Desc)
}
