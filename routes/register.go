package routes

import (
	"context"
	"flowee-api/types"
	"flowee-api/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request, db *mongo.Database, ctx context.Context) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	captcha := r.FormValue("captcha")

	if username == "" || password == "" || email == "" || captcha == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	check, err := utils.VerifyCaptcha(captcha)
	if err != nil || !check {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	err = db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}).Decode(nil)
	if err != mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 2}`)
		return
	}

	err = db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "email", Value: email}}).Decode(nil)
	if err != mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 3}`)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	db.Collection("accounts").InsertOne(ctx, types.Account{
		ID: primitive.NewObjectID(),
		Username: username,
		Email: email,
		Password: string(hash),
		Timestamp: int32(time.Now().Unix()),
		LastStream: 0,
		Token: utils.GenerateToken(),
		StreamToken: utils.GenerateToken(),
		StreamName: username + "'s stream",
		IsAcive: false,
	})

	log.Printf("%s registered", username)
	fmt.Fprintf(w, `{"success": true}`)
}