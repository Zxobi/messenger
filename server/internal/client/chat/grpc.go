package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/converter"
	"github.com/dvid-messanger/internal/domain/model"
	chatv1 "github.com/dvid-messanger/protos/gen/chat"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

var ErrUserChatsNotFound = errors.New("user chats not found")

type Client struct {
	api chatv1.ChatServiceClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "client.chat.New"

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
		api: chatv1.NewChatServiceClient(cc),
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) Create(ctx context.Context, from []byte, to []byte) (*model.Chat, error) {
	const op = "client.chat.Create"

	resp, err := c.api.Create(ctx, &chatv1.CreateChatRequest{FromUid: from, ToUid: to})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ChatFromDTO(resp.GetChat()), nil
}

func (c *Client) Chat(ctx context.Context, cid []byte) (*model.Chat, error) {
	const op = "client.chat.Chat"

	resp, err := c.api.Chat(ctx, &chatv1.ChatRequest{Cid: cid})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ChatFromDTO(resp.GetChat()), nil
}

func (c *Client) UserChats(ctx context.Context, uid []byte) ([]model.Chat, error) {
	const op = "client.chat.UserChats"

	resp, err := c.api.UserChats(ctx, &chatv1.UserChatsRequest{Uid: uid})
	if err != nil {
		errStatus, ok := status.FromError(err)
		if ok && errStatus.Code() == codes.NotFound {
			return make([]model.Chat, 0), nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ChatsFromDTO(resp.GetChats()), nil
}

func (c *Client) SendMessage(ctx context.Context, cid []byte, uid []byte, text string) (*model.ChatMessage, error) {
	const op = "client.chat.SendMessage"

	resp, err := c.api.SendMessage(ctx, &chatv1.SendMessageRequest{Cid: cid, Uid: uid, Text: text})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ChatMessageFromDTO(resp.GetMessage()), nil
}

func (c *Client) Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error) {
	const op = "client.chat.Messages"

	resp, err := c.api.Messages(ctx, &chatv1.MessagesRequest{Cid: cid})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ChatMessagesFromDTO(resp.GetMessages()), nil
}
