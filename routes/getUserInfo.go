package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flowee-ru/flowee-api/models"
	"github.com/flowee-ru/flowee-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

func GetUserInfo(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	username := r.FormValue("username")
	token := r.FormValue("token")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	var me *models.Account = nil
	if token != "" {
		acc, err := utils.GetAccountFromToken(ctx, db, token)
		if err != mongo.ErrNoDocuments {
			me = acc
		}
	}

	var account models.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}, primitive.E{Key: "isActive", Value: true}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	followers, err := db.Collection("follows").CountDocuments(ctx, bson.D{primitive.E{Key: "user2", Value: account.ID}})
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	isFollowing := false
	if me != nil {
		var follow models.Follow
		err := db.Collection("follows").FindOne(ctx, bson.D{primitive.E{Key: "user1", Value: me.ID}, primitive.E{Key: "user2", Value: account.ID}}).Decode(&follow)
		if err != mongo.ErrNoDocuments {
			isFollowing = true
		}
	}

	if me != nil {
		fmt.Fprintf(w, `{"success": true, "username": "%s", "accountID": "%s", "avatar": "%s", "followers": %d, "isLive": %t, "streamName": "%s", "isFollowing": %t}`, account.Username, account.ID.Hex(), account.Avatar, followers, account.IsLive, account.StreamName, isFollowing)
	} else {
		fmt.Fprintf(w, `{"success": true, "username": "%s", "accountID": "%s", "avatar": "%s", "followers": %d, "isLive": %t, "streamName": "%s"}`, account.Username, account.ID.Hex(), account.Avatar, followers, account.IsLive, account.StreamName)
	}
}