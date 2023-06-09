package users

import (
	"context"
	"fmt"
	"net/http"
	
	"github.com/flowee-ru/flowee-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	var account models.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}, primitive.E{Key: "isActive", Value: true}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	fmt.Fprintf(w, `{"success": true, "token": "%s"}`, account.Token)
}