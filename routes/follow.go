package routes

import (
	"context"
	"flowee-api/types"
	"flowee-api/utils"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Follow(w http.ResponseWriter, r *http.Request, db *mongo.Database, ctx context.Context) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	token := r.FormValue("token")
	targetIDHex := r.FormValue("targetID")

	if token == "" || targetIDHex == "" || !primitive.IsValidObjectID(targetIDHex) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	accountID, err := primitive.ObjectIDFromHex(targetIDHex)

	w.Header().Set("Content-Type", "application/json")

	account, err := utils.GetAccountFromToken(ctx, db, token)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	err = db.Collection("follows").FindOne(ctx, bson.D{
		primitive.E{Key: "user1", Value: account.ID},
		primitive.E{Key: "user2", Value: accountID},
	}).Decode(nil)
	if err != mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 2}`)
		return
	}

	db.Collection("follows").InsertOne(ctx, types.Follow{
		ID: primitive.NewObjectID(),
		User1: account.ID,
		User2: accountID,
		Timestamp: int32(time.Now().Unix()),
	})

	fmt.Fprintf(w, `{"success": true}`)
}