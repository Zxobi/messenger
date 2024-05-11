package auth

import (
	"context"
	"errors"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/core/service/auth"
	gvalidate "github.com/dvid-messanger/internal/pkg/grpc"
	authv1 "github.com/dvid-messanger/protos/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	authv1.UnimplementedAuthServiceServer
	auth primary.Auth
}

func Register(gRpc *grpc.Server, auth primary.Auth) {
	authv1.RegisterAuthServiceServer(gRpc, &serverApi{auth: auth})
}

func (s *serverApi) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverApi) Register(ctx context.Context, req *authv1.CreateRequest) (*authv1.CreateResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	uid, err := s.auth.Create(ctx, req.GetUid(), req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.CreateResponse{Uid: uid}, nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateRegister(req *authv1.CreateRequest) error {
	if err := gvalidate.ValidateId(req.GetUid(), "uid"); err != nil {
		return err
	}
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
