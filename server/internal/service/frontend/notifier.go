package frontend

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/domain/converter"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/id"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/lib/proto"
	frontendv1 "github.com/dvid-messanger/protos/gen/go/frontend"
	"log/slog"
	"strconv"
)

type Notifier struct {
	log *slog.Logger

	cp ClientProvider
	cr ChatRegistry
}

type ClientProvider interface {
	Clients(cid []byte) ([]Client, error)
}

type ChatRegistry interface {
	RegisterChat(chat *model.Chat) error
}

func NewNotifier(log *slog.Logger, cp ClientProvider, cr ChatRegistry) *Notifier {
	return &Notifier{log: log, cp: cp, cr: cr}
}

func (n *Notifier) NewMessage(ctx context.Context, message *model.ChatMessage) error {
	const op = "frontend.NewMessage"
	log := n.log.With(slog.String("op", op))

	log.Debug("notifying new chat message " + id.String(message.Cid))

	clients, err := n.cp.Clients(message.Cid)
	if err != nil {
		log.Error("failed to get clients", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	downstream, err := makeDownstream(message)
	if err != nil {
		log.Error("failed to make downstream message", logger.Err(err))
	}

	for _, c := range clients {
		c.Send(downstream)
	}

	log.Debug("notified " + strconv.Itoa(len(clients)) + " clients in " + id.String(message.Cid))

	return nil
}

func (n *Notifier) NewChat(ctx context.Context, chat *model.Chat) error {
	const op = "frontend.RegisterChat"
	log := n.log.With(slog.String("op", op))

	log.Debug("notifying new chat " + id.String(chat.Id))

	err := n.cr.RegisterChat(chat)
	if err != nil {
		log.Error("failed create new chat", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("notified new chat " + id.String(chat.Id))

	return nil
}

func makeDownstream(message *model.ChatMessage) ([]byte, error) {
	downstream := &frontendv1.DownstreamNewMessage{Message: converter.ChatMessageToDTO(message)}

	return proto.MarshalDownstream[*frontendv1.DownstreamNewMessage](
		downstream,
		frontendv1.DownstreamType_D_NEW_MESSAGE,
		nil,
	)
}
