package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/core/domain/modelutil"
	"github.com/dvid-messanger/internal/driver/secondary/storage/chat"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

const (
	nameDb                  = "db_chat"
	nameChatsCollection     = "chats"
	nameUserChatsCollection = "user_chats"
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

func (s *Storage) UserChats(ctx context.Context, uid []byte) (model.UserChats, error) {
	const op = "mongo.UserChats"
	log := s.log.With(slog.String("op", op))

	res := s.userChatsCollection().FindOne(ctx, bson.D{{"_id", uid}})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return model.UserChats{}, fmt.Errorf("%s: %w", op, chat.ErrUserChatsNotFound)
		}

		log.Error("failed to fetch user chats", logger.Err(res.Err()))
		return model.UserChats{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	var fetched model.UserChats
	if err := res.Decode(&fetched); err != nil {
		log.Error("failed to decode user chats", logger.Err(res.Err()))
		return model.UserChats{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return fetched, nil
}

func (s *Storage) UserChatsMany(ctx context.Context, uids [][]byte) ([]model.UserChats, error) {
	const op = "mongo.UserChats"
	log := s.log.With(slog.String("op", op))

	cursor, err := s.userChatsCollection().Find(ctx, bson.M{"_id": bson.M{"$in": uids}})
	if err != nil {
		log.Error("failed to fetch user chats", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	userChats := make([]model.UserChats, 0)
	if err = cursor.All(ctx, &userChats); err != nil {
		log.Error("failed to fetch user chats", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return userChats, nil
}

func (s *Storage) SaveUserChats(ctx context.Context, uid []byte) (model.UserChats, error) {
	const op = "mongo.SaveUserChats"
	log := s.log.With(slog.String("op", op))

	userChats := model.UserChats{Uid: uid}

	if _, err := s.userChatsCollection().InsertOne(ctx, userChats); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.UserChats{}, fmt.Errorf("%s: %w", op, chat.ErrUserChatsExists)
		}

		log.Error("failed to save user chats", logger.Err(err))
		return model.UserChats{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return userChats, nil
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

	transactionFunc := func(sessionContext mongo.SessionContext) (any, error) {
		userChats, err := s.UserChatsMany(sessionContext, [][]byte{from, to})
		if err != nil {
			return model.Chat{}, err
		}
		userChatsTo, err := s.getOrInsertUserChats(sessionContext, userChats, to)
		if err != nil {
			return model.Chat{}, err
		}
		userChatsFrom, err := s.getOrInsertUserChats(sessionContext, userChats, from)
		if err != nil {
			return model.Chat{}, err
		}

		if modelutil.HaveChatWith(&userChatsTo, from) || modelutil.HaveChatWith(&userChatsFrom, to) {
			return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrChatExists)
		}

		cid := uuid.New()
		newChat := model.NewPersonalChat(cid[:], from, to)

		if _, err = s.chatsCollection().InsertOne(sessionContext, newChat); err != nil {
			log.Error("failed to save chat", logger.Err(err))
			return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
		}
		if res, err := s.userChatsCollection().UpdateByID(
			sessionContext,
			userChatsFrom.Uid,
			bson.M{"$push": bson.M{"chats": modelutil.AddPersonalChat(&userChatsFrom, newChat.Id, to)}},
		); err != nil || res.ModifiedCount < 1 {
			if err == nil {
				log.Error("failed to update user chats, nothing modified")
			} else {
				log.Error("failed to update user chats", logger.Err(err))
			}
			return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
		}
		if res, err := s.userChatsCollection().UpdateByID(
			sessionContext,
			userChatsTo.Uid,
			bson.M{"$push": bson.M{"chats": modelutil.AddPersonalChat(&userChatsTo, newChat.Id, from)}},
		); err != nil || res.ModifiedCount < 1 {
			if err == nil {
				log.Error("failed to update user chats, nothing modified")
			} else {
				log.Error("failed to update user chats", logger.Err(err))
			}
			return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
		}

		return *newChat, nil
	}

	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		if errors.Is(err, chat.ErrChatExists) {
			return model.Chat{}, err
		}
		log.Error("failed to execute transaction", logger.Err(err))
		return model.Chat{}, fmt.Errorf("%s: %w", op, chat.ErrInternal)
	}

	return result.(model.Chat), nil
}

func (s *Storage) chatsCollection() *mongo.Collection {
	return s.Client.Database(nameDb).Collection(nameChatsCollection)
}

func (s *Storage) userChatsCollection() *mongo.Collection {
	return s.Client.Database(nameDb).Collection(nameUserChatsCollection)
}

func (s *Storage) getOrInsertUserChats(
	ctx context.Context,
	userChats []model.UserChats,
	uid []byte,
) (model.UserChats, error) {
	for _, v := range userChats {
		if bytes.Equal(v.Uid, uid) {
			return v, nil
		}
	}

	res, err := s.SaveUserChats(ctx, uid)
	if err != nil {
		return model.UserChats{}, err
	}

	return res, nil
}
