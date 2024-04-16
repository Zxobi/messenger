package chat

import "errors"

var (
	ErrChatExists           = errors.New("chat already exists")
	ErrChatNotFound         = errors.New("chat not found")
	ErrChatMessagesNotFound = errors.New("chat messages not found")
	ErrUserChatsNotFound    = errors.New("user chats not found")
)
