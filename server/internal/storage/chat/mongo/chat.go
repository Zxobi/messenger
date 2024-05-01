package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/storage/chat"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"slices"
)

const (
	nameDb                = "db_chat"
	nameChatsCollection   = "chats"
	nameMembersCollection = "members"
)

type Storage struct {
	log *slog.Logger
	mongodb.MongoDatabase
}

func New(log *slog.Logger, opts ...mongodb.Option) *Storage {
	storage := &Storage{log: log}
	storage.ApplyOptions(opts...)
	return storage
}

func (s *Storage) Connect(ctx context.Context) error {
	const op = "mongo.Connect"
	log := s.log.With(slog.String("op", op))

	if err := s.MongoDatabase.Connect(ctx); err != nil {
		log.Error("db connection failed", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("connected to db")
	return nil
}

func (s *Storage) Chat(ctx context.Context, cid []byte) (model.Chat, error) {
	const op = "mongo.Chat"
	log := s.log.With(slog.String("op", op))

	res := s.chatsCollection().FindOne(ctx, bson.D{{"_id", cid}})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrChatNotFound)
		}

		log.Error("failed to fetch chat", logger.Err(res.Err()))
		return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	var fetched model.Chat
	if err := res.Decode(&fetched); err != nil {
		log.Error("failed to decode chat", logger.Err(res.Err()))
		return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return fetched, nil
}

func (s *Storage) Chats(ctx context.Context, cids [][]byte) ([]model.Chat, error) {
	const op = "mongo.Chats"
	log := s.log.With(slog.String("op", op))

	cursor, err := s.chatsCollection().Find(ctx, bson.M{"_id": bson.M{"$in": cids}})
	if err != nil {
		log.Error("failed to fetch chats", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	chats := make([]model.Chat, 0)
	if err = cursor.All(ctx, &chats); err != nil {
		log.Error("failed to fetch chats", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return chats, nil
}

func (s *Storage) UserChats(ctx context.Context, uid []byte) ([]model.Chat, error) {
	const op = "mongo.UserChats"
	log := s.log.With(slog.String("op", op))

	res := s.membersCollection().FindOne(ctx, bson.D{{"_id", uid}})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, chat.ErrUserChatsNotFound)
		}

		log.Error("failed to fetch user chats", logger.Err(res.Err()))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	var fetched model.UserChats
	if err := res.Decode(&fetched); err != nil {
		log.Error("failed to decode user chats", logger.Err(res.Err()))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return s.Chats(ctx, fetched.Chats)
}

func (s *Storage) Save(ctx context.Context, from []byte, to []byte) (model.Chat, error) {
	const op = "mongo.Save"
	log := s.log.With(slog.String("op", op))

	session, err := s.Client.StartSession()
	if err != nil {
		log.Error("failed to start session", logger.Err(err))
		return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}
	defer session.EndSession(ctx)

	transactionFunc := func(sessionContext mongo.SessionContext) (interface{}, error) {
		chats, err := s.UserChats(ctx, from)
		if err != nil && !errors.Is(err, chat.ErrUserChatsNotFound) {
			_ = session.AbortTransaction(sessionContext)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if slices.ContainsFunc(chats, func(chat model.Chat) bool {
			return slices.ContainsFunc(chat.Members, func(member model.ChatMember) bool {
				return bytes.Equal(member.Uid, to)
			})
		}) {
			_ = session.AbortTransaction(sessionContext)
			return nil, fmt.Errorf("%s: %w", op, chat.ErrChatExists)
		}

		cid := uuid.New()
		newChat := model.Chat{
			Id:      cid[:],
			Type:    model.CTPersonal,
			Members: []model.ChatMember{{Uid: from}, {Uid: to}},
		}
		if _, err = s.chatsCollection().InsertOne(ctx, newChat); err != nil {
			log.Error("failed to save chat", logger.Err(err))
			_ = session.AbortTransaction(ctx)
			return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
		}

		// Подтверждаем транзакцию
		if err = session.CommitTransaction(sessionContext); err != nil {
			log.Error("failed to commit transaction", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
		}

		return newChat, nil
	}

	ret, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		return model.Chat{}, err
	}

	newChat, ok := ret.(model.Chat)
	if !ok {
		log.Error("unexpected return type from transaction")
		return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return newChat, nil
}

func (s *Storage) chatsCollection() *mongo.Collection {
	return s.Client.Database(nameDb).Collection(nameChatsCollection)
}

func (s *Storage) membersCollection() *mongo.Collection {
	return s.Client.Database(nameDb).Collection(nameMembersCollection)
}
