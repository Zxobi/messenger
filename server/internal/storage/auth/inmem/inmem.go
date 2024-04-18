package inmem

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/storage/auth"
	"sync"
)

type Storage struct {
	hash map[string]model.UserCredentials
	rw   *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		hash: make(map[string]model.UserCredentials),
		rw:   &sync.RWMutex{},
	}
}

func (s *Storage) SaveUser(ctx context.Context, uid []byte, email string, passHash []byte) (model.UserCredentials, error) {
	const op = "storage.inmem.SaveUser"

	s.rw.Lock()
	defer s.rw.Unlock()

	if _, ok := s.hash[email]; ok {
		return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrUserExists)
	}
	for _, userCred := range s.hash {
		if bytes.Equal(userCred.Id, uid) {
			return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrUserExists)
		}
	}

	userCred := model.UserCredentials{Email: email, Id: uid, PassHash: passHash}
	s.hash[email] = userCred

	return userCred, nil
}

func (s *Storage) User(ctx context.Context, email string) (model.UserCredentials, error) {
	const op = "storage.inmem.User"

	s.rw.RLock()
	defer s.rw.RUnlock()

	user, ok := s.hash[email]
	if !ok {
		return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
	}

	return user, nil
}
