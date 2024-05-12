package handler

import (
	"context"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route/middleware"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/id"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"log/slog"
)

type AuthHandler struct {
	log      *slog.Logger
	registry primary.ClientRegistry
	auth     primary.Auth
}

func RegisterAuthHandler(
	log *slog.Logger,
	r *route.Router,
	registry primary.ClientRegistry,
	auth primary.Auth,
	authMiddleware *middleware.AuthMiddleware,
) {
	handler := AuthHandler{
		log:      log,
		registry: registry,
		auth:     auth,
	}

	r.RegisterHandler(
		frontendv1.UpstreamType_U_LOGIN,
		frontendv1.DownstreamType_D_LOGIN,
		&frontendv1.UpstreamLogin{},
		handler.Login,
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_LOGOUT,
		frontendv1.DownstreamType_D_LOGOUT,
		&frontendv1.UpstreamLogout{},
		authMiddleware.WithAuth(handler.Logout),
	)

	handler.log.Debug("auth handler registered")
}

func (a *AuthHandler) Login(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.Login"
	log := a.log.With(slog.String("op", op), slog.String("c", id.String(request.ClientId)))

	upstream := request.Payload.(*frontendv1.UpstreamLogin)
	token, err := a.auth.Login(ctx, upstream.Email, upstream.Password)
	if err != nil {
		log.Error("failed to login", logger.Err(err))
		return &route.UpstreamResponse{ErrCode: frontendv1.ErrorCode_BAD_LOGIN, ErrDesc: "failed to login"}
	}

	if err = a.registry.SetAuth(request.ClientId, token); err != nil {
		log.Error("failed to set auth", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{
		Payload: &frontendv1.DownstreamLogin{Token: token},
	}
}

func (a *AuthHandler) Logout(_ context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.Logout"
	log := a.log.With(slog.String("op", op), slog.String("c", id.String(request.ClientId)))

	if err := a.registry.UnsetAuth(request.ClientId); err != nil {
		log.Error("failed to logout", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{
		Payload: &frontendv1.DownstreamLogout{},
	}
}
