package inmem

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/adapter/secondary/storage/chat"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/google/uuid"
	"slices"
	"sync"
)

type ChatStorage struct {
	hash map[[16]byte]model.Chat
	rw   *sync.RWMutex
}

func NewChatStorage() *ChatStorage {
	return &ChatStorage{
		hash: make(map[[16]byte]model.Chat),
		rw:   &sync.RWMutex{},
	}
}

func (s *ChatStorage) Chat(ctx context.Context, cid []byte) (model.Chat, error) {
	op := "storage.inmem.Chat"

	s.rw.RLock()
	defer s.rw.RUnlock()
	c, ok := s.hash[[16]byte(cid)]
	if !ok {
		return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrChatNotFound)
	}

	return c, nil
}

func (s *ChatStorage) UserChats(ctx context.Context, uid []byte) ([]model.Chat, error) {
	op := "storage.inmem.UserChats"

	s.rw.RLock()
	defer s.rw.RUnlock()

	uChats := make([]model.Chat, 0)
	for _, c := range s.hash {
		if containsUser(&c, uid) {
			uChats = append(uChats, c)
		}
	}

	if len(uChats) == 0 {
		return nil, fmt.Errorf("%s: %w", op, chat.ErrUserChatsNotFound)
	}
	return uChats, nil
}

func (s *ChatStorage) Save(ctx context.Context, from []byte, to []byte) (model.Chat, error) {
	op := "storage.inmem.Save"

	s.rw.Lock()
	defer s.rw.Unlock()
	for _, c := range s.hash {
		if containsUser(&c, from) && containsUser(&c, to) {
			return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrChatExists)
		}
	}

	cid := uuid.New()
	c := model.Chat{
		Id:      cid[:],
		Type:    model.CTPersonal,
		Members: []model.ChatMember{{Uid: from}, {Uid: to}},
	}
	s.hash[cid] = c

	return c, nil
}

func containsUser(c *model.Chat, uid []byte) bool {
	return slices.ContainsFunc(c.Members, func(m model.ChatMember) bool {
		return bytes.Equal(m.Uid, uid)
	})
}
