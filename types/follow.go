package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Follow struct {
	ID          primitive.ObjectID   `bson:"_id"`
	User1       primitive.ObjectID   `bson:"user1"`
	User2       primitive.ObjectID   `bson:"user2"`
	Timestamp   int32                `bson:"timestamp"`
}