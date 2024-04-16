package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/id"
	"log/slog"
	"slices"
	"strconv"
	"sync"
)

var (
	ErrClientAlreadyRegistered = errors.New("client already registered")
	ErrClientNotRegistered     = errors.New("client not registered")
	ErrChatAlreadyRegistered   = errors.New("chat already registered")
)

type Client interface {
	Send(payload []byte)
}

type ClientRegistry struct {
	log         *slog.Logger
	clientUser  map[Client]*model.User
	chatClients map[[16]byte][]Client

	mu sync.RWMutex
}

func NewClientRegistry(log *slog.Logger) *ClientRegistry {
	return &ClientRegistry{
		log:         log,
		clientUser:  make(map[Client]*model.User),
		chatClients: make(map[[16]byte][]Client),
	}
}

func (cr *ClientRegistry) Register(client Client, user *model.User, chats []model.Chat) error {
	const op = "frontend.Register"
	log := cr.log.With(slog.String("op", op))

	log.Debug("registering client " + id.String(user.Id))

	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, exists := cr.clientUser[client]
	if exists {
		log.Error("client already registered")
		return fmt.Errorf("%s: %w", op, ErrClientAlreadyRegistered)
	}

	cr.clientUser[client] = user
	for _, chat := range chats {
		cKey := [16]byte(chat.Id)
		clients, ok := cr.chatClients[cKey]
		if ok {
			cr.chatClients[cKey] = append(clients, client)
		} else {
			cr.chatClients[cKey] = []Client{client}
		}
	}

	log.Debug("client registered " + id.String(user.Id))

	return nil
}

func (cr *ClientRegistry) Unregister(client Client) error {
	const op = "frontend.Unregister"
	log := cr.log.With(slog.String("op", op))

	log.Debug("unregistering client")

	cr.mu.Lock()
	defer cr.mu.Unlock()
	user, exists := cr.clientUser[client]
	if !exists {
		log.Error("client not registered")
		return fmt.Errorf("%s: %w", op, ErrClientNotRegistered)
	}

	for cid, clients := range cr.chatClients {
		idx := slices.Index(clients, client)
		if idx == -1 {
			continue
		}

		clients[idx] = clients[len(clients)-1]
		cr.chatClients[cid] = clients[:len(clients)-1]
	}

	log.Debug("client unregistered " + id.String(user.Id))

	return nil
}

func (cr *ClientRegistry) Clients(cid []byte) ([]Client, error) {
	const op = "frontend.Clients"
	log := cr.log.With(slog.String("op", op))

	log.Debug("getting clients " + id.String(cid))

	cr.mu.RLock()
	defer cr.mu.RUnlock()
	clients := cr.chatClients[[16]byte(cid)]

	log.Debug("fetched clients " + id.String(cid))

	return clients, nil
}

func (cr *ClientRegistry) RegisterChat(chat *model.Chat) error {
	const op = "frontend.RegisterChat"
	log := cr.log.With(slog.String("op", op))

	log.Debug("registering new chat")

	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, exists := cr.chatClients[id.Id(chat.Id)]
	if exists {
		return fmt.Errorf("%s: %w", op, ErrChatAlreadyRegistered)
	}

	clients := make([]Client, 0)
	for client, user := range cr.clientUser {
		if slices.ContainsFunc(chat.Members, func(member model.ChatMember) bool {
			return bytes.Equal(member.Id, user.Id)
		}) {
			clients = append(clients, client)
		}
	}

	cr.chatClients[id.Id(chat.Id)] = clients

	log.Debug("registered new chat for " + strconv.Itoa(len(clients)) + " clients")

	return nil
}
