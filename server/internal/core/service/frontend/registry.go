package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/adapter/primary"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/pkg/id"
	"log/slog"
	"slices"
	"strconv"
	"sync"
)

var (
	ErrClientNotRegistered     = errors.New("client not registered")
	ErrClientAlreadyRegistered = errors.New("client already registered")
	ErrClientNotAuthorized     = errors.New("client not authorized")
	ErrClientAlreadyAuthorized = errors.New("client already authorized")

	ErrChatAlreadyRegistered = errors.New("chat already registered")
)

const (
	claimAuth = "auth"
)

type Client interface {
	Send(msg []byte) error
	GetId() []byte
}

type clientWithClaims struct {
	client Client
	claims map[string]interface{}
}

type ClientRegistry struct {
	log         *slog.Logger
	clients     map[[16]byte]*clientWithClaims
	clientUser  map[[16]byte]*model.User
	chatClients map[[16]byte][]Client

	mu *sync.RWMutex
}

func NewClientRegistry(log *slog.Logger) *ClientRegistry {
	return &ClientRegistry{
		log:         log,
		clients:     make(map[[16]byte]*clientWithClaims),
		clientUser:  make(map[[16]byte]*model.User),
		chatClients: make(map[[16]byte][]Client),
		mu:          &sync.RWMutex{},
	}
}

func (cr *ClientRegistry) Register(client primary.Client) error {
	const op = "registry.Register"
	log := cr.log.With(slog.String("op", op), slog.String("c", id.String(client.GetId())))

	log.Debug("registering")

	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, exists := cr.clients[id.Id(client.GetId())]
	if exists {
		log.Error("already registered")
		return fmt.Errorf("%s: %w", op, ErrClientAlreadyRegistered)
	}

	cr.clients[id.Id(client.GetId())] = &clientWithClaims{client: client, claims: make(map[string]interface{})}

	log.Debug("registered")
	return nil
}

func (cr *ClientRegistry) SetAuth(clientId []byte, auth string, user *model.User, chats []model.Chat) error {
	const op = "registry.SetAuth"
	log := cr.log.With(slog.String("op", op), slog.String("c", id.String(clientId)))

	log.Debug("setting auth for user " + id.String(user.Id))

	cr.mu.Lock()
	defer cr.mu.Unlock()

	clientClaims, ok := cr.clients[id.Id(clientId)]
	if !ok {
		log.Error("not registered")
		return fmt.Errorf("%s: %w", op, ErrClientNotRegistered)
	}
	if clientClaims.claims[claimAuth] != nil {
		log.Error("already authorized")
		return fmt.Errorf("%s: %w", op, ErrClientAlreadyAuthorized)
	}

	clientClaims.claims[claimAuth] = auth
	cr.clientUser[id.Id(clientId)] = user
	for _, chat := range chats {
		cr.chatClients[id.Id(chat.Id)] = append(cr.chatClients[id.Id(chat.Id)], clientClaims.client)
	}

	log.Debug("authorized as " + id.String(user.Id))
	return nil
}

func (cr *ClientRegistry) Auth(clientId []byte) (string, error) {
	const op = "registry.Auth"
	log := cr.log.With(slog.String("op", op), slog.String("c", id.String(clientId)))

	log.Debug("getting auth")

	cr.mu.Lock()
	defer cr.mu.Unlock()
	clientClaims, ok := cr.clients[id.Id(clientId)]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrClientNotRegistered)
	}

	auth, ok := clientClaims.claims[claimAuth]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrClientNotAuthorized)
	}

	return auth.(string), nil
}

func (cr *ClientRegistry) Unregister(clientId []byte) error {
	const op = "registry.Unregister"
	log := cr.log.With(slog.String("op", op), slog.String("c", id.String(clientId)))

	log.Debug("unregistering")

	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, exists := cr.clients[id.Id(clientId)]
	if !exists {
		log.Error("not registered")
		return fmt.Errorf("%s: %w", op, ErrClientNotRegistered)
	}

	delete(cr.clients, id.Id(clientId))
	delete(cr.clientUser, id.Id(clientId))
	for cid, clients := range cr.chatClients {
		idx := slices.IndexFunc(clients, func(client Client) bool {
			return bytes.Equal(client.GetId(), clientId)
		})
		if idx == -1 {
			continue
		}

		clients[idx] = clients[len(clients)-1]
		cr.chatClients[cid] = clients[:len(clients)-1]
	}

	log.Debug("unregistered")

	return nil
}

func (cr *ClientRegistry) Clients(cid []byte) ([]Client, error) {
	const op = "registry.Clients"
	log := cr.log.With(slog.String("op", op))

	log.Debug("getting clients " + id.String(cid))

	cr.mu.RLock()
	defer cr.mu.RUnlock()
	clients := cr.chatClients[id.Id(cid)]

	log.Debug("fetched clients " + id.String(cid))

	return clients, nil
}

func (cr *ClientRegistry) RegisterChat(chat *model.Chat) error {
	const op = "registry.RegisterChat"
	log := cr.log.With(slog.String("op", op))

	log.Debug("registering new chat " + id.String(chat.Id))

	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, exists := cr.chatClients[id.Id(chat.Id)]
	if exists {
		return fmt.Errorf("%s: %w", op, ErrChatAlreadyRegistered)
	}

	clients := make([]Client, 0)
	for client, user := range cr.clientUser {
		if slices.ContainsFunc(chat.Members, func(member model.ChatMember) bool {
			return bytes.Equal(member.Uid, user.Id)
		}) {
			clients = append(clients, cr.clients[client].client)
		}
	}

	cr.chatClients[id.Id(chat.Id)] = clients

	log.Debug("registered new chat for " + strconv.Itoa(len(clients)) + " clients")

	return nil
}
