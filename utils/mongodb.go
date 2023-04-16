package utils

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context) (*mongo.Database, error) {
	mongoUri := "mongodb://localhost:27017"
	if os.Getenv("MONGO_URI") != "" {
		mongoUri = os.Getenv("MONGO_URI")
	}

	mongoDB := "flowee"
	if os.Getenv("MONGO_DB") != "" {
		mongoDB = os.Getenv("MONGO_DB")
	}

	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(mongoDB), nil
}