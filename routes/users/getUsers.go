package users

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/flowee-ru/flowee-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	cur, err := db.Collection("accounts").Find(ctx, bson.D{primitive.E{Key: "isActive", Value: true}, primitive.E{Key: "isLive", Value: true}})
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	var result []string
	for cur.Next(ctx) {
		var account models.Account
		cur.Decode(&account)
		result = append(result, `{"username": "` + account.Username + `", ` +
			`"avatar": "` + account.Avatar + `", ` +
			`"streamName": "` + account.StreamName + `"}`)
	}

	fmt.Fprintf(w, `{"success": true, "users": [%s]}`, strings.Join(result, ", "))
}