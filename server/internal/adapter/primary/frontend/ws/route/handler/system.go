package handler

import (
	"context"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"log/slog"
)

type SystemHandler struct {
	log *slog.Logger
}

func RegisterSystemHandler(
	log *slog.Logger,
	r *route.Router,
) {
	handler := SystemHandler{log: log}

	r.RegisterHandler(
		frontendv1.UpstreamType_U_ECHO,
		frontendv1.DownstreamType_D_ECHO,
		&frontendv1.UpstreamEcho{},
		handler.Echo,
	)

	handler.log.Debug("system handler registered")
}

func (r *SystemHandler) Echo(_ context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	upstream := request.Payload.(*frontendv1.UpstreamEcho)
	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamEcho{Content: upstream.GetContent()}}
}
