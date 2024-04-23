package chat

import (
	"bytes"
	"context"
	"errors"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/storage/chat"
	"github.com/dvid-messanger/mocks/mock_chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"time"
)

var log = slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

func TestCreateChat(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockChatSaver := &mock_chat.MockChatSaver{}
		mockChatNotifier := &mock_chat.MockChatNotifier{}

		expected := model.Chat{
			Id:      []byte("mockChatId"),
			Members: []model.ChatMember{{Id: []byte("from")}, {Id: []byte("to")}},
		}
		mockChatSaver.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(expected, nil)
		mockChatNotifier.On("NewChat", mock.Anything, mock.Anything).Return(nil)

		service := NewService(log, nil, mockChatSaver, mockChatNotifier, nil, nil, nil)

		createdChat, err := service.Create(context.Background(), []byte("from"), []byte("to"))
		require.NoError(t, err)
		assert.Equal(t, expected, *createdChat)
	})
	t.Run("ChatExistsError", func(t *testing.T) {
		t.Parallel()

		mockChatSaver := &mock_chat.MockChatSaver{}
		mockChatNotifier := &mock_chat.MockChatNotifier{}

		mockChatSaver.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(
			model.Chat{}, chat.ErrChatExists)

		service := NewService(log, nil, mockChatSaver, mockChatNotifier, nil, nil, nil)

		_, err := service.Create(context.Background(), []byte("from"), []byte("to"))
		assert.ErrorIs(t, err, ErrChatExists)
	})
}

func TestGetChat(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockChatProvider := &mock_chat.MockChatProvider{}
		expected := model.Chat{
			Id:      []byte("mockChatId"),
			Members: []model.ChatMember{{Id: []byte("from")}, {Id: []byte("to")}},
		}
		mockChatProvider.On("Chat", mock.Anything, mock.Anything).Return(expected, nil)

		service := NewService(log, mockChatProvider, nil, nil, nil, nil, nil)

		res, err := service.Chat(context.Background(), []byte("mockChatId"))
		require.NoError(t, err)
		assert.Equal(t, expected, *res)
	})
	t.Run("NotFoundError", func(t *testing.T) {
		t.Parallel()

		mockChatProvider := &mock_chat.MockChatProvider{}

		mockChatProvider.On("Chat", mock.Anything, mock.Anything).Return(
			model.Chat{}, chat.ErrChatNotFound)

		service := NewService(log, mockChatProvider, nil, nil, nil, nil, nil)

		_, err := service.Chat(context.Background(), []byte("nonexistentChatId"))
		assert.ErrorIs(t, err, ErrChatNotFound)
	})
}

func TestUserChats(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockUserChatProvider := &mock_chat.MockUserChatProvider{}

		expectedChats := []model.Chat{
			{Id: []byte("chat1"), Members: []model.ChatMember{{Id: []byte("user1")}, {Id: []byte("user2")}}},
			{Id: []byte("chat2"), Members: []model.ChatMember{{Id: []byte("user1")}, {Id: []byte("user3")}}},
		}

		mockUserChatProvider.On("UserChats", mock.Anything, mock.Anything).Return(expectedChats, nil)

		service := NewService(log, nil, nil, nil, mockUserChatProvider, nil, nil)

		chats, err := service.UserChats(context.Background(), []byte("user1"))
		require.NoError(t, err)
		assert.Equal(t, expectedChats, chats)
	})
	t.Run("NotFoundError", func(t *testing.T) {
		t.Parallel()

		mockUserChatProvider := &mock_chat.MockUserChatProvider{}

		mockUserChatProvider.On("UserChats", mock.Anything, mock.Anything).Return(nil, chat.ErrUserChatsNotFound)

		service := NewService(log, nil, nil, nil, mockUserChatProvider, nil, nil)

		_, err := service.UserChats(context.Background(), []byte("nonexistentUser"))
		assert.ErrorIs(t, err, ErrUserChatsNotFound)
	})
}

func TestMessages(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockMessageProvider := &mock_chat.MockMessageProvider{}

		expectedMessages := []model.ChatMessage{
			{Id: []byte("msg1"), Text: "Hello", Uid: []byte("user1"), Cid: []byte("chat1"), Timestamp: time.Now().UnixMilli()},
			{Id: []byte("msg2"), Text: "Hi", Uid: []byte("user2"), Cid: []byte("chat1"), Timestamp: time.Now().UnixMilli()},
		}

		mockMessageProvider.On("Messages", mock.Anything, mock.Anything).Return(expectedMessages, nil)

		service := NewService(log, nil, nil, nil, nil, mockMessageProvider, nil)

		messages, err := service.Messages(context.Background(), []byte("chatID"))
		require.NoError(t, err)
		assert.Equal(t, expectedMessages, messages)
	})
	t.Run("NotFoundError", func(t *testing.T) {
		t.Parallel()

		mockMessageProvider := &mock_chat.MockMessageProvider{}
		mockMessageProvider.On("Messages", mock.Anything, mock.Anything).
			Return(nil, chat.ErrMessagesNotFound)

		service := NewService(log, nil, nil, nil, nil, mockMessageProvider, nil)

		_, err := service.Messages(context.Background(), []byte("nonexistentChatID"))
		assert.ErrorIs(t, err, ErrMessagesNotFound)
	})
}

func TestSendMessage(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockMessageSaver := &mock_chat.MockMessageSaver{}
		mockChatNotifier := &mock_chat.MockChatNotifier{}

		expectedMessage := model.ChatMessage{
			Id:        []byte("msgId"),
			Text:      "Hello",
			Uid:       []byte("user"),
			Cid:       []byte("chat"),
			Timestamp: time.Now().UnixMilli(),
		}

		mockMessageSaver.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(expectedMessage, nil)
		mockChatNotifier.On("NewMessage", mock.Anything, mock.Anything).Return(nil)

		service := NewService(log, nil, nil, mockChatNotifier, nil, nil, mockMessageSaver)

		message, err := service.SendMessage(context.Background(), []byte("chat"), []byte("user"), "Hello")
		require.NoError(t, err)
		assert.Equal(t, expectedMessage, *message)
	})
	t.Run("MessageSaveFailedError", func(t *testing.T) {
		t.Parallel()

		mockMessageSaver := &mock_chat.MockMessageSaver{}
		mockChatNotifier := &mock_chat.MockChatNotifier{}

		mockMessageSaver.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(model.ChatMessage{}, errors.New("failed to save message"))

		service := NewService(log, nil, nil, mockChatNotifier, nil, nil, mockMessageSaver)

		_, err := service.SendMessage(context.Background(), []byte("chat"), []byte("user"), "Hello")
		assert.Error(t, err)
	})
	t.Run("NotifyFailedNoError", func(t *testing.T) {
		t.Parallel()

		mockMessageSaver := &mock_chat.MockMessageSaver{}
		mockChatNotifier := &mock_chat.MockChatNotifier{}

		expectedMessage := model.ChatMessage{
			Id:        []byte("msgId"),
			Text:      "Hello",
			Uid:       []byte("user"),
			Cid:       []byte("chat"),
			Timestamp: time.Now().UnixMilli(),
		}

		mockMessageSaver.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(expectedMessage, nil)
		mockChatNotifier.On("NewMessage", mock.Anything, mock.Anything).
			Return(errors.New("failed to notify new message"))

		service := NewService(log, nil, nil, mockChatNotifier, nil, nil, mockMessageSaver)

		message, err := service.SendMessage(context.Background(), []byte("chat"), []byte("user"), "Hello")
		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, *message)
	})
}
