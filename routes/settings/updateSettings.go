package settings

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateSettings(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	token := r.FormValue("token")

	streamName := r.FormValue("streamName")

	if streamName != "" {
		test, _ := regexp.MatchString("^\\s*$", streamName)
		if test {
			fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
			return
		}
	
		_, err := db.Collection("accounts").UpdateOne(ctx, bson.D{primitive.E{Key: "token", Value: token}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "streamName", Value: streamName}}}})
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(w, `{"success": false, "errorCode": 2}`)
			return
		}
	
		fmt.Fprintf(w, `{"success": true}`)
	} else {
		fmt.Fprintf(w, `{"success": false}`)
	}
}