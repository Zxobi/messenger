package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/storage/chat"
	"log/slog"
)

type ChatService struct {
	log *slog.Logger
	cp  ChatProvider
	cs  ChatSaver
	cn  ChatNotifier

	ucp UserChatProvider

	mp MessageProvider
	ms MessageSaver
}

type ChatProvider interface {
	Chat(cid []byte) (model.Chat, error)
}

type UserChatProvider interface {
	UserChats(ctx context.Context, uid []byte) ([]model.Chat, error)
}

type ChatSaver interface {
	Save(ctx context.Context, from []byte, to []byte) (model.Chat, error)
}

type MessageProvider interface {
	Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error)
}

type MessageSaver interface {
	Save(ctx context.Context, cid []byte, from []byte, text string) (model.ChatMessage, error)
}

type ChatNotifier interface {
	NewMessage(ctx context.Context, message *model.ChatMessage) error
	NewChat(ctx context.Context, chat *model.Chat) error
}

var (
	ErrChatExists        = errors.New("chat already exists")
	ErrUserChatsNotFound = errors.New("user chats not found")
	ErrChatNotFound      = errors.New("chat not found")

	ErrMessagesNotFound = errors.New("chat messages not found")
)

func NewService(
	log *slog.Logger,
	cp ChatProvider,
	cs ChatSaver,
	cn ChatNotifier,
	ucp UserChatProvider,
	mp MessageProvider,
	ms MessageSaver,
) *ChatService {
	return &ChatService{
		log: log,
		cp:  cp,
		cs:  cs,
		cn:  cn,
		ucp: ucp,
		mp:  mp,
		ms:  ms,
	}
}

func (s *ChatService) Create(ctx context.Context, from []byte, to []byte) (*model.Chat, error) {
	const op = "chat.Create"
	log := s.log.With(slog.String("op", op))

	log.Debug("creating chat")

	c, err := s.cs.Save(ctx, from, to)
	if err != nil {
		if errors.Is(err, chat.ErrChatExists) {
			return nil, fmt.Errorf("%s: %w", op, ErrChatExists)
		}

		log.Error("failed to save chat", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("chat created")

	err = s.cn.NewChat(ctx, &c)
	if err != nil {
		log.Error("failed to notify new chat", logger.Err(err))
	}

	return &c, nil
}

func (s *ChatService) Chat(ctx context.Context, cid []byte) (*model.Chat, error) {
	const op = "chat.Create"
	log := s.log.With(slog.String("op", op))

	log.Debug("getting chat")

	c, err := s.cp.Chat(cid)
	if err != nil {
		if errors.Is(err, chat.ErrChatNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrChatNotFound)
		}

		log.Error("failed to get chat", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("chat fetched")

	return &c, nil
}

func (s *ChatService) UserChats(ctx context.Context, uid []byte) ([]model.Chat, error) {
	const op = "chat.UserChats"
	log := s.log.With(slog.String("op", op))

	log.Debug("getting user chats")

	chats, err := s.ucp.UserChats(ctx, uid)
	if err != nil {
		if errors.Is(err, chat.ErrUserChatsNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserChatsNotFound)
		}

		log.Error("failed to get user chats", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("user chats fetched")
	return chats, nil
}

func (s *ChatService) Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error) {
	const op = "chat.Messages"
	log := s.log.With(slog.String("op", op))

	log.Debug("getting messages")

	messages, err := s.mp.Messages(ctx, cid)
	if err != nil {
		if errors.Is(err, chat.ErrUserChatsNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrMessagesNotFound)
		}

		log.Error("failed to get messages", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("messages fetched")
	return messages, nil
}

func (s *ChatService) SendMessage(ctx context.Context, cid []byte, from []byte, text string) (*model.ChatMessage, error) {
	const op = "chat.SendMessage"
	log := s.log.With(slog.String("op", op))

	log.Debug("saving message")

	m, err := s.ms.Save(ctx, cid, from, text)
	if err != nil {
		log.Error("failed to save chat", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("message saved")

	err = s.cn.NewMessage(ctx, &m)
	if err != nil {
		log.Error("failed to notify new message", logger.Err(err))
	}

	return &m, nil
}
