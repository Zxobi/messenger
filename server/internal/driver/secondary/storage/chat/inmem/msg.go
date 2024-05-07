package inmem

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/driver/secondary/storage/chat"
	"github.com/dvid-messanger/internal/pkg/cutils"
	"github.com/google/uuid"
	"sync"
	"time"
)

type MessageStorage struct {
	hash map[[16]byte][]model.ChatMessage
	rw   *sync.RWMutex
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		hash: make(map[[16]byte][]model.ChatMessage),
		rw:   &sync.RWMutex{},
	}
}

func (s *MessageStorage) Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error) {
	op := "storage.inmem.Messages"

	s.rw.RLock()
	defer s.rw.RUnlock()
	c, ok := s.hash[[16]byte(cid)]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, chat.ErrMessagesNotFound)
	}

	return cutils.Copy(c), nil
}

func (s *MessageStorage) Save(ctx context.Context, cid []byte, from []byte, text string) (model.ChatMessage, error) {
	keyCid := [16]byte(cid)
	mid := uuid.New()
	cm := model.ChatMessage{
		Id:        mid[:],
		Cid:       cid,
		Uid:       from,
		Text:      text,
		Timestamp: time.Now().UnixMilli(),
	}

	s.rw.Lock()
	defer s.rw.Unlock()
	messages, ok := s.hash[keyCid]
	if !ok {
		s.hash[keyCid] = []model.ChatMessage{cm}
	}
	s.hash[keyCid] = append(messages, cm)

	return cm, nil
}
