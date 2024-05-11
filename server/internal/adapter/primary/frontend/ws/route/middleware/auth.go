package middleware

import (
	"context"
	"encoding/base64"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route"
	"github.com/dvid-messanger/internal/pkg/logger"
	"log/slog"
)

type AuthMiddleware struct {
	log      *slog.Logger
	registry primary.ClientRegistry
	tv       primary.TokenVerifier
}

func NewAuthMiddleware(log *slog.Logger, registry primary.ClientRegistry, tv primary.TokenVerifier) *AuthMiddleware {
	return &AuthMiddleware{log: log, registry: registry, tv: tv}
}

func (a *AuthMiddleware) WithAuth(fun route.HandlerFunc) route.HandlerFunc {
	return func(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
		const op = "middleware.WithAuth"
		log := a.log.With(slog.String("op", op))

		token, err := a.registry.Auth(request.ClientId)
		if err != nil {
			log.Error("failed to get client auth", logger.Err(err))
			return route.ErrResponseUnauthorized
		}

		claims, err := a.tv.Verify(token)
		if err != nil {
			log.Error("failed to verify auth token", logger.Err(err))
			return route.ErrResponseUnauthorized
		}

		uidClaim, ok := claims["uid"]
		if !ok {
			log.Error("uid claim missing")
			return route.ErrResponseUnauthorized
		}
		uidStr, ok := uidClaim.(string)
		if !ok {
			log.Error("uid claim is not a string")
			return route.ErrResponseUnauthorized
		}
		uid, err := base64.StdEncoding.DecodeString(uidStr)
		if err != nil {
			log.Error("failed to decode uid string", logger.Err(err))
			return route.ErrResponseUnauthorized
		}

		request.AuthUid = uid
		return fun(ctx, request)
	}
}
