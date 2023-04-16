package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID          primitive.ObjectID `bson:"_id"`
	Username    string             `bson:"username"`
	Password    string             `bson:"passsword"`
	Email       string             `bson:"email"`
	Timestamp   int                `bson:"timestamp"`
	Token       string             `bson:"token"`
	StreamToken string             `bson:"streamToken"`
	StreamName  string             `bson:"streamName"`
}