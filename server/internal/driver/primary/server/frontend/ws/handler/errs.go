package handler

import (
	"github.com/dvid-messanger/internal/pkg/proto"
	"github.com/dvid-messanger/protos/gen/frontend"
)

func mustMakeUnauthorized(dt frontendv1.DownstreamType) []byte {
	msg, _ := proto.MarshalDownstream[*frontendv1.Upstream](
		nil,
		dt,
		&frontendv1.DownstreamError{Code: frontendv1.ErrorCode_UNAUTHORIZED, Desc: "unauthorized"},
	)

	return msg
}
