package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID              primitive.ObjectID   `bson:"_id"`
	Username        string               `bson:"username"`
	Password        string               `bson:"passsword"`
	Email           string               `bson:"email"`
	Timestamp       int32                `bson:"timestamp"`
	AvatarURL       string               `bson:"avatarURL"`
	LastStream      int32                `bson:"lastStream"`
	LastEmailSend   int32                `bson:"lastEmailSend"`
	Token           string               `bson:"token"`
	VerifyToken     string               `bson:"verifyToken"`
	StreamToken     string               `bson:"streamToken"`
	StreamName      string               `bson:"streamName"`
	StreamURL       string               `bson:"streamURL"`
	IsLive          bool                 `bson:"isLive"`
	IsActive        bool                 `bson:"isActive"`
}