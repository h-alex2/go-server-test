package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	dbName string
}

func GetMongoDBClient(uri, dbName string) (*MongoDB, error, context.Context, context.CancelFunc) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)

	return &MongoDB{client, dbName}, err, ctx, cancel
}

func (db *MongoDB) CloseMongoDB() error {
	return db.client.Disconnect(context.Background())
}

func (db *MongoDB) GetCollection(collection string) *mongo.Collection {
	return db.client.Database(db.dbName).Collection(collection)
}
