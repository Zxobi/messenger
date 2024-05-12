package handler

import (
	"context"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route/middleware"
	"github.com/dvid-messanger/internal/core/domain/converter"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/id"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

type InfoHandler struct {
	log      *slog.Logger
	registry primary.ClientRegistry
	chat     primary.Chat
	user     primary.User
}

func RegisterInfoHandler(
	log *slog.Logger,
	r *route.Router,
	registry primary.ClientRegistry,
	chat primary.Chat,
	user primary.User,
	authMiddleware *middleware.AuthMiddleware,
) {
	handler := InfoHandler{
		log:      log,
		registry: registry,
		chat:     chat,
		user:     user,
	}

	r.RegisterHandler(
		frontendv1.UpstreamType_U_INFO_INIT,
		frontendv1.DownstreamType_D_INFO_INIT,
		&frontendv1.UpstreamInfoInit{},
		authMiddleware.WithAuth(handler.InitInfo),
	)

	handler.log.Debug("info handler registered")
}

func (a *InfoHandler) InitInfo(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.InitInfo"
	log := a.log.With(slog.String("op", op), slog.String("c", id.String(request.ClientId)))

	var user *model.User
	var chats []model.Chat

	eg, gCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		usr, err := a.user.User(gCtx, request.AuthUid)
		if err != nil {
			return err
		}

		user = usr
		return nil
	})
	eg.Go(func() error {
		uChats, err := a.chat.UserChats(gCtx, request.AuthUid)
		if err != nil {
			return err
		}

		chats = uChats
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Error("failed to get info", logger.Err(err))
		return route.ErrResponseInternal
	}
	if err := a.registry.SetInfo(request.ClientId, user, chats); err != nil {
		log.Error("failed to set info", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{
		Payload: &frontendv1.DownstreamInfoInit{User: converter.UserToDTO(user), Chats: converter.ChatsToDTO(chats)},
	}
}
