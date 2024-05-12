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
	handler := ChatHandler{log: log, chat: chat}

	r.RegisterHandler(
		frontendv1.UpstreamType_U_GET_USER_CHATS,
		frontendv1.DownstreamType_D_GET_USER_CHATS,
		&frontendv1.UpstreamGetUserChats{},
		auth.WithAuth(handler.GetUserChats),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_GET_CHAT,
		frontendv1.DownstreamType_D_GET_CHAT,
		&frontendv1.UpstreamGetChat{},
		auth.WithAuth(handler.GetChat),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_CREATE_CHAT,
		frontendv1.DownstreamType_D_CREATE_CHAT,
		&frontendv1.UpstreamCreateChat{},
		auth.WithAuth(handler.CreateChat),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_SEND_MESSAGE,
		frontendv1.DownstreamType_D_SEND_MESSAGE,
		&frontendv1.UpstreamSendMessage{},
		auth.WithAuth(handler.SendMessage),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_CHAT_MESSAGES,
		frontendv1.DownstreamType_D_CHAT_MESSAGES,
		&frontendv1.UpstreamChatMessages{},
		auth.WithAuth(handler.ChatMessages),
	)

	handler.log.Debug("chat handler registered")
}

func (r *ChatHandler) GetUserChats(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.GetUserChats"
	log := r.log.With(slog.String("op", op))

	chats, err := r.chat.UserChats(ctx, request.AuthUid)
	if err != nil {
		log.Error("failed to get user chats", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamGetUserChats{Chats: converter.ChatsToDTO(chats)}}
}

func (r *ChatHandler) GetChat(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.GetChat"
	log := r.log.With(slog.String("op", op))

	upstream := request.Payload.(*frontendv1.UpstreamGetChat)

	chat, err := r.chat.Chat(ctx, upstream.GetCid())
	if err != nil {
		log.Error("failed to get chat", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamGetChat{Chat: converter.ChatToDTO(chat)}}
}

func (r *ChatHandler) CreateChat(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.CreateChat"
	log := r.log.With(slog.String("op", op))

	upstream := request.Payload.(*frontendv1.UpstreamCreateChat)

	chat, err := r.chat.Create(ctx, request.AuthUid, upstream.GetUid())
	if err != nil {
		log.Error("failed to create chat", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamCreateChat{Chat: converter.ChatToDTO(chat)}}
}

func (r *ChatHandler) SendMessage(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.SendMessage"
	log := r.log.With(slog.String("op", op))

	upstream := request.Payload.(*frontendv1.UpstreamSendMessage)

	msg, err := r.chat.SendMessage(ctx, upstream.GetCid(), request.AuthUid, upstream.GetText())
	if err != nil {
		log.Error("failed to send message", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamNewMessage{Message: converter.ChatMessageToDTO(msg)}}
}

func (r *ChatHandler) ChatMessages(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.ChatMessages"
	log := r.log.With(slog.String("op", op))

	upstream := request.Payload.(*frontendv1.UpstreamChatMessages)

	msgs, err := r.chat.Messages(ctx, upstream.GetCid())
	if err != nil {
		log.Error("failed to get chat messages", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamChatMessages{Messages: converter.ChatMessagesToDTO(msgs)}}
}
