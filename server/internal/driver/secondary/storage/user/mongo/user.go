package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/driver/secondary/storage/user"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

const (
	nameDb              = "db_user"
	nameUsersCollection = "users"
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

func (s *Storage) User(ctx context.Context, uid []byte) (model.User, error) {
	const op = "mongo.User"
	log := s.log.With(slog.String("op", op))

	var usr model.User
	res := s.usersCollection().FindOne(ctx, bson.D{{"_id", uid}})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return usr, fmt.Errorf("%s: %w", op, user.ErrUserNotFound)
		}

		log.Error("fetch user error", logger.Err(res.Err()))
		return usr, fmt.Errorf("%s: %w", op, user.ErrInternal)
	}

	if err := res.Decode(&usr); err != nil {
		log.Error("failed to decode user", logger.Err(err))
		return usr, fmt.Errorf("%s: %w", op, user.ErrInternal)
	}

	return usr, nil
}

func (s *Storage) Users(ctx context.Context) ([]model.User, error) {
	const op = "mongo.Users"
	log := s.log.With(slog.String("op", op))

	cursor, err := s.usersCollection().Find(ctx, bson.M{}, nil)
	if err != nil {
		log.Error("failed to fetch users", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, user.ErrInternal)
	}

	users := make([]model.User, 0)
	if err = cursor.All(ctx, &users); err != nil {
		log.Error("failed to fetch users", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, user.ErrInternal)
	}

	return users, nil
}

func (s *Storage) Save(ctx context.Context, email string, bio string) (model.User, error) {
	const op = "mongo.Save"
	log := s.log.With(slog.String("op", op))

	uid := uuid.New()
	usr := model.User{Id: uid[:], Email: email, Bio: bio}

	if _, err := s.usersCollection().InsertOne(ctx, usr); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.User{}, fmt.Errorf("%s: %w", op, user.ErrUserExists)
		}

		log.Error("failed to save user", logger.Err(err))
		return model.User{}, fmt.Errorf("%s: %w", op, user.ErrInternal)
	}

	return usr, nil
}

func (s *Storage) usersCollection() *mongo.Collection {
	return s.Client.Database(nameDb).Collection(nameUsersCollection)
}
