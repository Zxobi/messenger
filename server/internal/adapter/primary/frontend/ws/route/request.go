package route

import (
	"fmt"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/golang/protobuf/proto"
)

type UpstreamRequest struct {
	ClientId []byte
	AuthUid  []byte
	Payload  proto.Message
}

type UpstreamResponse struct {
	ErrCode frontendv1.ErrorCode
	ErrDesc string
	Payload proto.Message
}

var ErrResponseInternal = &UpstreamResponse{
	ErrCode: frontendv1.ErrorCode_INTERNAL,
	ErrDesc: "internal error",
}

var ErrResponseTimeout = &UpstreamResponse{
	ErrCode: frontendv1.ErrorCode_TIMEOUT,
	ErrDesc: "request timeout",
}

var ErrResponseUnauthorized = &UpstreamResponse{
	ErrCode: frontendv1.ErrorCode_UNAUTHORIZED,
	ErrDesc: "unauthorized",
}

func MarshalResponse(response *UpstreamResponse, dt frontendv1.DownstreamType) ([]byte, error) {
	const op = "request.MakeResponse"

	var payload []byte
	var dErr frontendv1.DownstreamError
	if response.ErrCode == 0 {
		if response.Payload != nil {
			marshalled, err := proto.Marshal(response.Payload)
			if err != nil {
				return nil, err
			}
			payload = marshalled
		}
		return nil, nil
	} else {
		dErr.Code = response.ErrCode
		dErr.Desc = response.ErrDesc
	}

	downstream := &frontendv1.Downstream{
		Type:    dt,
		Payload: payload,
		Error:   &dErr,
	}

	res, err := proto.Marshal(downstream)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
