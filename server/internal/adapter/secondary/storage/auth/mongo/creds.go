package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dvid-messanger/internal/adapter/secondary/storage/auth"
	"github.com/dvid-messanger/internal/core/domain/model"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

const (
	nameDb              = "db_auth"
	nameCredsCollection = "creds"
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

func (s *Storage) Save(ctx context.Context, uid []byte, email string, passHash []byte) (model.UserCredentials, error) {
	const op = "mongo.Save"
	log := s.log.With(slog.String("op", op))

	creds := model.UserCredentials{Id: uid[:], Email: email, PassHash: passHash}
	if _, err := s.credsCollection().InsertOne(ctx, creds); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.UserCredentials{}, auth.ErrUserExists
		}

		log.Error("failed to save creds", logger.Err(err))
		return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrInternal)
	}

	return creds, nil
}

func (s *Storage) User(ctx context.Context, email string) (model.UserCredentials, error) {
	const op = "mongo.User"
	log := s.log.With(slog.String("op", op))

	res := s.credsCollection().FindOne(ctx, bson.D{{"email", email}})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}

		log.Error("failed to fetch creds", logger.Err(res.Err()))
		return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrInternal)
	}

	var creds model.UserCredentials
	if err := res.Decode(&creds); err != nil {
		log.Error("failed to decode creds", logger.Err(err))
		return model.UserCredentials{}, fmt.Errorf("%s: %w", op, auth.ErrInternal)
	}

	return creds, nil
}

func (s *Storage) credsCollection() *mongo.Collection {
	return s.Client.Database(nameDb).Collection(nameCredsCollection)
}
