package user

import (
	"context"
	"errors"
	"github.com/dvid-messanger/internal/core/domain/converter"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/core/service/user"
	grpcutil "github.com/dvid-messanger/internal/pkg/grpc"
	userv1 "github.com/dvid-messanger/protos/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User interface {
	Create(ctx context.Context, email string, bio string) (*model.User, error)
	User(ctx context.Context, uid []byte) (*model.User, error)
	Users(ctx context.Context) ([]model.User, error)
}

type serverApi struct {
	userv1.UnimplementedUserServiceServer
	user User
}

func Register(gRpc *grpc.Server, user User) {
	userv1.RegisterUserServiceServer(gRpc, &serverApi{user: user})
}

func (s *serverApi) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	usr, err := s.user.Create(ctx, req.GetEmail(), req.GetBio())
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userv1.CreateResponse{User: converter.UserToDTO(usr)}, nil
}

func (s *serverApi) User(ctx context.Context, req *userv1.UserRequest) (*userv1.UserResponse, error) {
	if err := validateUser(req); err != nil {
		return nil, err
	}

	usr, err := s.user.User(ctx, req.GetUid())
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userv1.UserResponse{User: converter.UserToDTO(usr)}, nil
}

func (s *serverApi) Users(ctx context.Context, req *userv1.UsersRequest) (*userv1.UsersResponse, error) {
	users, err := s.user.Users(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userv1.UsersResponse{Users: converter.UsersToDTO(users)}, nil
}

func validateUser(req *userv1.UserRequest) error {
	if err := grpcutil.ValidateId(req.GetUid(), "uid"); err != nil {
		return err
	}

	return nil
}

func validateRegister(req *userv1.CreateRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}
