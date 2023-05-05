package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flowee-ru/flowee-api/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Unfollow(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	token := r.FormValue("token")
	accountIDHex := mux.Vars(r)["accountID"]

	if token == "" || accountIDHex == "" || !primitive.IsValidObjectID(accountIDHex) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	accountID, _ := primitive.ObjectIDFromHex(accountIDHex)

	account, err := utils.GetAccountFromToken(ctx, db, token)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	err = db.Collection("follows").FindOne(ctx, bson.D{
		primitive.E{Key: "user1", Value: account.ID},
		primitive.E{Key: "user2", Value: accountID},
	}).Decode(nil)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": true}`)
		return
	}

	db.Collection("follows").DeleteOne(ctx, bson.D{
		primitive.E{Key: "user1", Value: account.ID},
		primitive.E{Key: "user2", Value: accountID},
	})

	fmt.Fprintf(w, `{"success": true}`)
}