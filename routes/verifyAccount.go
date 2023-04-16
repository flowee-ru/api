package routes

import (
	"context"
	"fmt"
	"net/http"
	"flowee-api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyAccount(w http.ResponseWriter, r *http.Request, db *mongo.Database, ctx context.Context) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	token := r.FormValue("verifyToken")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var account types.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "verifyToken", Value: token}, primitive.E{Key: "isActive", Value: false}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	db.Collection("accounts").UpdateOne(ctx, bson.D{primitive.E{Key: "verifyToken", Value: token}},
	bson.D{primitive.E{
		Key: "$set", Value: bson.D{primitive.E{
			Key: "isActive", Value: true,
		}},
	}})

	fmt.Fprintf(w, `{"success": true, "token": "%s"}`, account.Token)
}