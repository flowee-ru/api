package settings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flowee-ru/flowee-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSettings(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	token := r.FormValue("token")

	var account models.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "token", Value: token}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}
	
	streamToken := account.ID.Hex() + "?t=" + account.StreamToken

	fmt.Fprintf(w, `{"success": true, "settings": {
		"streamToken": "%s",
		"streamName": "%s"
	}}`, streamToken, account.StreamName)
}