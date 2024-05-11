package auth

import (
	"context"
	"fmt"
	authv1 "github.com/dvid-messanger/protos/gen/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api authv1.AuthServiceClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "client.auth.New"

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
		api: authv1.NewAuthServiceClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) Create(ctx context.Context, uid []byte, email string, pass string) ([]byte, error) {
	const op = "client.auth.Register"

	resp, err := c.api.Create(ctx, &authv1.CreateRequest{Uid: uid, Email: email, Password: pass})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetUid(), nil
}

func (c *Client) Login(ctx context.Context, email string, pass string) (string, error) {
	const op = "client.auth.Login"

	resp, err := c.api.Login(ctx, &authv1.LoginRequest{Email: email, Password: pass})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetToken(), nil
}
