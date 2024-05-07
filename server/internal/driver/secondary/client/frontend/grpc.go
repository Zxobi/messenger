package frontend

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/core/domain/converter"
	"github.com/dvid-messanger/internal/core/domain/model"
	frontendv1 "github.com/dvid-messanger/protos/gen/frontend"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api frontendv1.NotifierClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "client.frontend.New"

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
		api: frontendv1.NewNotifierClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) NewMessage(ctx context.Context, message *model.ChatMessage) error {
	const op = "client.frontend.NewMessage"

	_, err := c.api.NewMessage(ctx, &frontendv1.NewMessageRequest{Message: converter.ChatMessageToDTO(message)})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) NewChat(ctx context.Context, chat *model.Chat) error {
	const op = "client.frontend.NewMessage"

	_, err := c.api.NewChat(ctx, &frontendv1.NewChatRequest{Chat: converter.ChatToDTO(chat)})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
