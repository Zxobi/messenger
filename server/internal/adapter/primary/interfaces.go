package primary

import (
	"context"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

type User interface {
	Create(ctx context.Context, email string, bio string) (*model.User, error)
	User(ctx context.Context, uid []byte) (*model.User, error)
	Users(ctx context.Context) ([]model.User, error)
}

type Auth interface {
	Create(ctx context.Context, uid []byte, email string, pass string) ([]byte, error)
	Login(ctx context.Context, email string, pass string) (string, error)
}

type Chat interface {
	Create(ctx context.Context, fromUid []byte, toUid []byte) (*model.Chat, error)

	Chat(ctx context.Context, cid []byte) (*model.Chat, error)
	UserChats(ctx context.Context, uid []byte) ([]model.Chat, error)

	SendMessage(ctx context.Context, cid []byte, uid []byte, text string) (*model.ChatMessage, error)
	Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error)
}

type Notifier interface {
	NewMessage(ctx context.Context, message *model.ChatMessage) error
	NewChat(ctx context.Context, message *model.Chat) error
}

type TokenVerifier interface {
	Verify(token string) (jwt.MapClaims, error)
}

type Client interface {
	GetId() []byte
	Send(msg []byte) error
}

type ClientRegistry interface {
	Register(client Client) error
	SetAuth(id []byte, auth string, user *model.User, chats []model.Chat) error
	Auth(id []byte) (string, error)
	Unregister(id []byte) error
}
