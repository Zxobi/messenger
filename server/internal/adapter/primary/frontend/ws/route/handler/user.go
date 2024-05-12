package handler

import (
	"context"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route"
	"github.com/dvid-messanger/internal/adapter/primary/frontend/ws/route/middleware"
	"github.com/dvid-messanger/internal/core/domain/converter"
	"github.com/dvid-messanger/internal/pkg/logger"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"log/slog"
)

type UserHandler struct {
	log  *slog.Logger
	user primary.User
	auth primary.Auth
}

func RegisterUserHandler(
	log *slog.Logger,
	r *route.Router,
	user primary.User,
	auth primary.Auth,
	authMiddleware *middleware.AuthMiddleware,
) {
	handler := UserHandler{log: log, user: user, auth: auth}

	r.RegisterHandler(
		frontendv1.UpstreamType_U_CUR_USER,
		frontendv1.DownstreamType_D_CUR_USER,
		&frontendv1.UpstreamCurUser{},
		authMiddleware.WithAuth(handler.GetCurUser),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_GET_USER,
		frontendv1.DownstreamType_D_GET_USER,
		&frontendv1.UpstreamGetUser{},
		authMiddleware.WithAuth(handler.GetUser),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_GET_USERS,
		frontendv1.DownstreamType_D_GET_USERS,
		&frontendv1.UpstreamGetUsers{},
		authMiddleware.WithAuth(handler.GetUsers),
	)
	r.RegisterHandler(
		frontendv1.UpstreamType_U_REG_USER,
		frontendv1.DownstreamType_D_REG_USER,
		&frontendv1.UpstreamRegUser{},
		authMiddleware.WithAuth(handler.RegUser),
	)

	handler.log.Debug("user handler registered")
}

func (r *UserHandler) GetCurUser(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.GetCurUser"
	log := r.log.With(slog.String("op", op))

	user, err := r.user.User(ctx, request.AuthUid)
	if err != nil {
		log.Error("failed to get cur user", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamCurUser{User: converter.UserToDTO(user)}}
}

func (r *UserHandler) GetUser(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.GetUser"
	log := r.log.With(slog.String("op", op))

	upstream := request.Payload.(*frontendv1.UpstreamGetUser)

	user, err := r.user.User(ctx, upstream.GetUid())
	if err != nil {
		log.Error("failed to get user", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamGetUser{User: converter.UserToDTO(user)}}
}

func (r *UserHandler) GetUsers(ctx context.Context, _ *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.GetUsers"
	log := r.log.With(slog.String("op", op))

	users, err := r.user.Users(ctx)
	if err != nil {
		log.Error("failed to get users", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamGetUsers{Users: converter.UsersToDTO(users)}}
}

func (r *UserHandler) RegUser(ctx context.Context, request *route.UpstreamRequest) *route.UpstreamResponse {
	const op = "handler.RegUser"
	log := r.log.With(slog.String("op", op))

	upstream := request.Payload.(*frontendv1.UpstreamRegUser)

	usr, err := r.user.Create(ctx, upstream.GetEmail(), upstream.GetBio())
	if err != nil {
		log.Error("failed to create user", logger.Err(err))
		return route.ErrResponseInternal
	}
	_, err = r.auth.Create(ctx, usr.Id, upstream.GetEmail(), upstream.GetPassword())
	if err != nil {
		log.Error("failed to create auth", logger.Err(err))
		return route.ErrResponseInternal
	}

	return &route.UpstreamResponse{Payload: &frontendv1.DownstreamRegUser{User: converter.UserToDTO(usr)}}
}
