package routes

import (
	"context"
	"fmt"
	"net/http"
	
	"github.com/flowee-ru/flowee-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyToken(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	token := r.FormValue("token")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	var account models.Account
	err := db.Collection("accounts").FindOne(context.TODO(), bson.D{primitive.E{Key: "token", Value: token}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	fmt.Fprintf(w, `{"success": true, "username": "%s"}`, account.Username)
}