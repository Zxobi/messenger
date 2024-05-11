package proto

import (
	"fmt"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
	"reflect"
)

func Unmarshal[T proto.Message](msg []byte, base T) error {
	const op = "proto.Unmarshal"

	if err := proto.Unmarshal(msg, base); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func MarshalDownstream[T proto.Message](msg T, msgType frontendv1.DownstreamType, dErr *frontendv1.DownstreamError) ([]byte, error) {
	const op = "proto.MarshalDownstream"

	downstream := &frontendv1.Downstream{Type: msgType, Error: dErr}

	if !reflect.ValueOf(msg).IsZero() {
		marshalled, err := proto.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		downstream.Payload = marshalled
	}

	res, err := proto.Marshal(downstream)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
