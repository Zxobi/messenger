package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/driver/secondary/storage/auth"
	"github.com/dvid-messanger/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Service struct {
	log      *slog.Logger
	us       UserSaver
	up       UserProvider
	tm       TokenMaker
	tokenTtl time.Duration
}

type UserSaver interface {
	Save(ctx context.Context, uid []byte, email string, passHash []byte) (model.UserCredentials, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (model.UserCredentials, error)
}

type TokenMaker interface {
	MakeToken(user model.UserCredentials, duration time.Duration) (string, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

func NewService(log *slog.Logger, us UserSaver, up UserProvider, tm TokenMaker, tokenTtl time.Duration) *Service {
	return &Service{
		log:      log,
		us:       us,
		up:       up,
		tm:       tm,
		tokenTtl: tokenTtl,
	}
}

func (a *Service) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))

	log.Debug("attempt to login user")

	user, err := a.up.User(ctx, email)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			a.log.Debug("user not found", logger.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Debug("invalid credentials", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Debug("user logged in successfully")

	token, err := a.tm.MakeToken(user, a.tokenTtl)
	if err != nil {
		a.log.Error("failed to make token", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Service) Create(ctx context.Context, uid []byte, email string, password string) ([]byte, error) {
	const op = "auth.Create"
	log := a.log.With(slog.String("op", op))

	log.Debug("creating")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", logger.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userCred, err := a.us.Save(ctx, uid, email, passHash)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			a.log.Debug("user already exists", logger.Err(err))

			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("created")

	return userCred.Id, nil
}
