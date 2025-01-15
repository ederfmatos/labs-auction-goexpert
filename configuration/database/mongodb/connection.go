package mongodb

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const (
	_mongoDBURL      = "MONGODB_URL"
	_mongoDBDatabase = "MONGODB_DB"
)

func NewMongoDBConnection() (*mongo.Database, error) {
	mongoURL := os.Getenv(_mongoDBURL)
	mongoDatabase := os.Getenv(_mongoDBDatabase)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.Error("Error trying to connect to mongodb database", err)
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		logger.Error("Error trying to ping mongodb database", err)
		return nil, err
	}

	return client.Database(mongoDatabase), nil
}
