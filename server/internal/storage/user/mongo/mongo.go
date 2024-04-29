package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/domain/model"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/storage/user"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type Storage struct {
	log *slog.Logger
	mongodb.MongoDatabase
	users        *mongo.Collection
	dbName       string
	usersColName string
}

func New(log *slog.Logger, dbName string, usersColName string, opts ...mongodb.Option) *Storage {
	storage := &Storage{log: log, dbName: dbName, usersColName: usersColName}
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

	db := s.Client.Database(s.dbName)
	s.users = db.Collection(s.usersColName)
	return nil
}

func (s *Storage) User(ctx context.Context, uid []byte) (model.User, error) {
	const op = "mongo.User"
	log := s.log.With(slog.String("op", op))

	var usr model.User
	res := s.users.FindOne(ctx, bson.D{{"_id", uid}})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return usr, user.ErrUserNotFound
		}

		log.Error("fetch user error", logger.Err(res.Err()))
		return usr, user.ErrInternal
	}

	if err := res.Decode(&usr); err != nil {
		log.Error("failed to decode user", logger.Err(err))
		return usr, user.ErrInternal
	}

	return usr, nil
}

func (s *Storage) Users(ctx context.Context) ([]model.User, error) {
	const op = "mongo.Users"
	log := s.log.With(slog.String("op", op))

	cursor, err := s.users.Find(ctx, bson.D{{}}, nil)
	if err != nil {
		log.Error("failed to fetch users", logger.Err(err))
		return nil, user.ErrInternal
	}

	users := make([]model.User, 0)
	for cursor.Next(ctx) {
		var usr model.User
		if err = cursor.Decode(&usr); err != nil {
			log.Error("failed to decode users", logger.Err(err))
			return nil, user.ErrInternal
		}

		users = append(users, usr)
	}

	if cursor.Err() != nil {
		log.Error("cursor error", logger.Err(cursor.Err()))
		return nil, user.ErrInternal
	}
	_ = cursor.Close(ctx)

	return users, nil
}

func (s *Storage) Save(ctx context.Context, email string, bio string) (model.User, error) {
	const op = "mongo.Save"
	log := s.log.With(slog.String("op", op))

	uid := uuid.New()
	usr := model.User{Id: uid[:], Email: email, Bio: bio}

	if _, err := s.users.InsertOne(ctx, usr); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.User{}, user.ErrUserExists
		}

		log.Error("failed to save user", logger.Err(err))
		return model.User{}, user.ErrInternal
	}

	return usr, nil
}
