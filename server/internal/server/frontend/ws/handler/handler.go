package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/server/frontend/ws"
	"github.com/dvid-messanger/internal/service/frontend"
	feproto "github.com/dvid-messanger/protos/gen/go/frontend"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/protobuf/proto"
	"log/slog"
)

var ErrUnknownUpstreamType = errors.New("unknown upstream type")

type Handler struct {
	log      *slog.Logger
	user     User
	auth     Auth
	chat     Chat
	tv       TokenVerifier
	registry ClientRegistry
}

type User interface {
	Create(ctx context.Context, email string, bio string) (*model.User, error)
	User(ctx context.Context, uid []byte) (*model.User, error)
	Users(ctx context.Context) ([]model.User, error)
}

type Auth interface {
	Register(ctx context.Context, uid []byte, email string, pass string) ([]byte, error)
	Login(ctx context.Context, email string, pass string) (string, error)
}

type Chat interface {
	Create(ctx context.Context, fromUid []byte, toUid []byte) (*model.Chat, error)

	Chat(ctx context.Context, cid []byte) (*model.Chat, error)
	UserChats(ctx context.Context, uid []byte) ([]model.Chat, error)

	SendMessage(ctx context.Context, cid []byte, uid []byte, text string) (*model.ChatMessage, error)
	Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error)
}

type ClientRegistry interface {
	Register(client frontend.Client, user *model.User, chats []model.Chat) error
	Unregister(client frontend.Client) error
}

type TokenVerifier interface {
	Verify(token string) (jwt.MapClaims, error)
}

func NewHandler(log *slog.Logger, user User, auth Auth, chat Chat, tv TokenVerifier, registry ClientRegistry) *Handler {
	return &Handler{
		log:      log,
		user:     user,
		auth:     auth,
		chat:     chat,
		tv:       tv,
		registry: registry,
	}
}

func (h *Handler) HandleMsg(c *ws.Client, msg []byte) error {
	const op = "handler.HandleMsg"
	log := h.log.With(slog.String("op", op))

	log.Debug("handling client msg")

	upstream := &feproto.Upstream{}
	if err := proto.Unmarshal(msg, upstream); err != nil {
		log.Error("failed to unmarshal upstream")
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := h.handleUpstream(c, upstream); err != nil {
		log.Error("failed to handle msg")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("msg handled", slog.Attr{
		Key:   "type",
		Value: slog.StringValue(upstream.Type.String()),
	})
	return nil
}

func (h *Handler) Register(c *ws.Client) {
	const op = "handler.Register"
	log := h.log.With(slog.String("op", op))

	log.Debug("register clint", slog.String("ip", c.GetAddr().String()))
}

func (h *Handler) Unregister(c *ws.Client) {
	const op = "handler.Unregister"
	log := h.log.With(slog.String("op", op))

	log.Info("unregister client", slog.String("ip", c.GetAddr().String()))

	err := h.registry.Unregister(c)
	if err != nil {
		log.Error("failed to unregister client", logger.Err(err))
	}
}

func (h *Handler) handleUpstream(c *ws.Client, upstream *feproto.Upstream) error {
	const op = "handler.handleUpstream"
	log := h.log.With(slog.String("op", op))

	var err error
	switch upstream.Type {
	case feproto.UpstreamType_U_ECHO:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamEcho{}, h.handleEcho)
	case feproto.UpstreamType_U_GET_USERS:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamGetUsers{}, h.handleGetUsers)
	case feproto.UpstreamType_U_GET_USER:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamGetUser{}, h.handleGetUser)
	case feproto.UpstreamType_U_CUR_USER:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamCurUser{}, h.handleCurUser)
	case feproto.UpstreamType_U_REG_USER:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamRegUser{}, h.handleRegUser)
	case feproto.UpstreamType_U_LOGIN:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamLogin{}, h.handleLogin)
	case feproto.UpstreamType_U_GET_USER_CHATS:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamGetUserChats{}, h.handleGetUserChats)
	case feproto.UpstreamType_U_CREATE_CHAT:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamCreateChat{}, h.handleCreateChat)
	case feproto.UpstreamType_U_SEND_MESSAGE:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamSendMessage{}, h.handleSendMessage)
	case feproto.UpstreamType_U_CHAT_MESSAGES:
		err = handleUpstreamFunc(c, upstream, &feproto.UpstreamChatMessages{}, h.handleChatMessages)

	default:
		log.Error("unknown upstream type " + upstream.Type.String())
		return ErrUnknownUpstreamType
	}

	if err != nil {
		log.Error("failed to handle", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func handleUpstreamFunc[T proto.Message](
	c *ws.Client,
	upstream *feproto.Upstream,
	msg T,
	handleFunc func(c *ws.Client, msg T) error,
) error {
	const op = "handler.handleUpstreamFunc"

	err := proto.Unmarshal(upstream.Payload, msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = handleFunc(c, msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
