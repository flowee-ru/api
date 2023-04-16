package utils

import (
	"context"
	"math/rand"
	"flowee-api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateToken(length int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetAccountFromToken(ctx context.Context, db *mongo.Database, token string) (*types.Account, error) {
	var account types.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "token", Value: token}}).Decode(&account)
	if err != nil {
		return nil, err
	}
	
	return &account, nil
}