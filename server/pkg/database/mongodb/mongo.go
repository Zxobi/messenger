package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Client  *mongo.Client
	options *options.ClientOptions
}

func (db *MongoDatabase) ApplyOptions(opts ...Option) {
	db.options = &options.ClientOptions{}
	for _, opt := range opts {
		opt(db)
	}
}

func (db *MongoDatabase) Connect(ctx context.Context) error {
	if err := db.options.Validate(); err != nil {
		return err
	}

	client, err := mongo.Connect(ctx, db.options)
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	db.Client = client
	return nil
}

func (db *MongoDatabase) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}
