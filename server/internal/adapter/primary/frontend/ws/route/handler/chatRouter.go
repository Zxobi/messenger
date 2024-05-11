package handler

import (
	"context"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route/middleware"
	"github.com/dvid-messanger/internal/core/domain/converter"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/protos/gen/frontend"
	"log/slog"
)

type ChatHandler struct {
	log  *slog.Logger
	chat primary.Chat
}

func RegisterChatHandler(log *slog.Logger, r *route.Router, chat primary.Chat, auth *middleware.AuthMiddleware) {
	cr := ChatHandler{log: log, chat: chat}

	r.RegisterHandler(
		frontendv1.UpstreamType_U_GET_USER_CHATS,
		frontendv1.DownstreamType_D_GET_USER_CHATS,
		&frontendv1.UpstreamGetUserChats{},
		auth.WithAuth(cr.GetUserChats),
	)
}

func (r *ChatHandler) GetUserChats(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.GetUserChats"
	log := r.log.With(slog.String("op", op))

	rChat, err := r.chat.UserChats(ctx, request.AuthUid)
	if err != nil {
		log.Error("failed to get user chats", logger.Err(err))
		return route.ErrResponseInternal
	}

	downstream := &frontendv1.DownstreamGetUserChats{Chats: converter.ChatsToDTO(rChat)}
	return &route.UpstreamResponse{Payload: downstream}
}
