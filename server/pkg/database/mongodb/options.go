package mongodb

import "time"

type Option func(db *MongoDatabase)

func Timeout(timeout time.Duration) Option {
	return func(db *MongoDatabase) {
		db.options.Timeout = &timeout
	}
}

func URI(uri string) Option {
	return func(db *MongoDatabase) {
		db.options.ApplyURI(uri)
	}
}
