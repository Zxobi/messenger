package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/storage/user"
	"log/slog"
)

type Service struct {
	log *slog.Logger
	us  UserSaver
	up  UserProvider
}

type UserSaver interface {
	Save(ctx context.Context, email string, bio string) (model.User, error)
}

type UserProvider interface {
	User(ctx context.Context, uid []byte) (model.User, error)
	Users(ctx context.Context) ([]model.User, error)
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

func NewUser(log *slog.Logger, us UserSaver, up UserProvider) *Service {
	return &Service{
		log: log,
		us:  us,
		up:  up,
	}
}

func (u *Service) Create(ctx context.Context, email string, bio string) (*model.User, error) {
	const op = "user.Create"
	log := u.log.With(slog.String("op", op))

	log.Debug("registering user")

	usr, err := u.us.Save(ctx, email, bio)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			u.log.Debug("user already exists", logger.Err(err))

			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("user registered")

	return &usr, nil
}

func (u *Service) User(ctx context.Context, uid []byte) (*model.User, error) {
	const op = "user.User"
	log := u.log.With(slog.String("op", op))

	log.Debug("getting user")

	usr, err := u.up.User(ctx, uid)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			u.log.Debug("user not found", logger.Err(err))

			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		log.Error("failed to get user", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("user fetched")

	return &usr, nil
}

func (u *Service) Users(ctx context.Context) ([]model.User, error) {
	const op = "user.Users"
	log := u.log.With(slog.String("op", op))

	log.Debug("getting users")

	users, err := u.up.Users(ctx)
	if err != nil {
		log.Error("failed to get users", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("users fetched")

	return users, nil
}
