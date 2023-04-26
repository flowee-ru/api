package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatMessage struct {
	ID          primitive.ObjectID   `bson:"_id"`
	AccountID   primitive.ObjectID   `bson:"accountID"`
	Content     string               `bson:"content"`
	Timestamp   int                  `bson:"timestamp"`
}