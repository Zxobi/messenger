package scylla

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/adapter/secondary/storage/chat"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/database/scylladb"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"log/slog"
)

const (
	statementSelectByCid = "SELECT mid, uid, text FROM messages WHERE cid=?"
	statementInsert      = "INSERT INTO messages(cid, mid, uid, text) VALUES(?,?,?,?)"
)

type Storage struct {
	log     *slog.Logger
	session *gocql.Session
}

func New(log *slog.Logger) *Storage {
	return &Storage{
		log: log,
	}
}

func (s *Storage) Connect(cluster *gocql.ClusterConfig) error {
	const op = "scylla.Connect"
	log := s.log.With(slog.String("op", op))

	session, err := scylladb.NewSession(cluster)
	if err != nil {
		log.Info("scylla connection failed", logger.Err(err))
		return fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	s.session = session
	return nil
}

func (s *Storage) Messages(ctx context.Context, cid []byte) ([]model.ChatMessage, error) {
	const op = "scylla.Messages"
	log := s.log.With(slog.String("op", op))

	iter := s.session.Query(statementSelectByCid, cid).Iter()
	scanner := iter.Scanner()
	defer iter.Close()

	var mid, uid = make([]byte, 0), make([]byte, 0)
	var text string

	res := make([]model.ChatMessage, 0)
	for scanner.Next() {
		if err := scanner.Scan(mid, uid, &text); err != nil {
			log.Error("failed to scan", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
		}
		res = append(res, model.ChatMessage{Id: mid, Cid: cid, Uid: uid, Text: text})
	}

	if scanner.Err() != nil {
		log.Error("failed to fetch messages", logger.Err(scanner.Err()))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return res, nil
}

func (s *Storage) Save(ctx context.Context, cid []byte, from []byte, text string) (model.ChatMessage, error) {
	const op = "scylla.Save"
	log := s.log.With(slog.String("op", op))

	mid := uuid.New()
	chatMessage := model.ChatMessage{Id: mid[:], Cid: cid, Uid: from, Text: text}

	if err := s.session.Query(
		statementInsert,
		chatMessage.Cid,
		chatMessage.Id,
		chatMessage.Uid,
		chatMessage.Text,
	).Exec(); err != nil {
		log.Error("failed to save message", logger.Err(err))
		return model.ChatMessage{}, chat.ErrInternal
	}

	return chatMessage, nil
}
