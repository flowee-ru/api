package routes

import (
	"context"
	"flowee-api/types"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.TODO()

func Login(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var account types.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	log.Printf("%s logged in", r.FormValue("username"))
	fmt.Fprintf(w, `{"success": true, "token": "%s"}`, account.Token)
}