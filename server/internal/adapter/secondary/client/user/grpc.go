package user

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/core/domain/converter"
	"github.com/dvid-messanger/internal/core/domain/model"
	userv1 "github.com/dvid-messanger/protos/gen/user"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api userv1.UserServiceClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "client.user.New"

	retryOpts := []retry.CallOption{
		retry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		retry.WithMax(uint(retriesCount)),
		retry.WithPerRetryTimeout(timeout),
	}

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			retry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: userv1.NewUserServiceClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) Create(ctx context.Context, email string, bio string) (*model.User, error) {
	const op = "client.user.Create"

	resp, err := c.api.Create(ctx, &userv1.CreateRequest{Email: email, Bio: bio})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.UserFromDTO(resp.GetUser()), nil
}

func (c *Client) User(ctx context.Context, uid []byte) (*model.User, error) {
	const op = "client.user.User"

	resp, err := c.api.User(ctx, &userv1.UserRequest{Uid: uid})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.UserFromDTO(resp.GetUser()), nil
}

func (c *Client) Users(ctx context.Context) ([]model.User, error) {
	const op = "client.user.Users"

	resp, err := c.api.Users(ctx, &userv1.UsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.UsersFromDTO(resp.GetUsers()), nil
}
