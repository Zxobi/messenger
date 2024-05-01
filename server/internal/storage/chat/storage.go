package chat

import "errors"

var (
	ErrChatExists        = errors.New("chat already exists")
	ErrChatNotFound      = errors.New("chat not found")
	ErrMessagesNotFound  = errors.New("chat messages not found")
	ErrUserChatsNotFound = errors.New("user chats not found")
	ErrInternal          = errors.New("storage error")
)
