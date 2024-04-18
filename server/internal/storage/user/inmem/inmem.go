package inmem

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/id"
	"github.com/dvid-messanger/internal/storage/user"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"sync"
)

type Storage struct {
	hash map[[16]byte]model.User
	rw   *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		hash: make(map[[16]byte]model.User),
		rw:   &sync.RWMutex{},
	}
}

func (s *Storage) SaveUser(ctx context.Context, email string, bio string) (model.User, error) {
	const op = "storage.inmem.SaveUser"

	s.rw.Lock()
	defer s.rw.Unlock()

	for _, usr := range s.hash {
		if usr.Email == email {
			return model.User{}, fmt.Errorf("%s: %w", op, user.ErrUserExists)
		}
	}

	uid := uuid.New()
	usr := model.User{Id: uid[:], Email: email, Bio: bio}
	s.hash[uid] = usr
	return usr, nil
}

func (s *Storage) User(ctx context.Context, uid []byte) (model.User, error) {
	const op = "storage.inmem.User"

	s.rw.RLock()
	defer s.rw.RUnlock()

	u, ok := s.hash[id.Id(uid)]
	if !ok {
		return model.User{}, fmt.Errorf("%s: %w", op, user.ErrUserNotFound)
	}

	return u, nil
}

func (s *Storage) Users(ctx context.Context) ([]model.User, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	return maps.Values(s.hash), nil
}
