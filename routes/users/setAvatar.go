package users

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flowee-ru/flowee-api/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetAvatar(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	token := r.FormValue("token")
	accountIDHex := mux.Vars(r)["accountID"]

	if token == "" || accountIDHex == "" || !primitive.IsValidObjectID(accountIDHex) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	apiHost := os.Getenv("VITE_API_HOST")

	accountID, _ := primitive.ObjectIDFromHex(accountIDHex)

	var account models.Account
	err := db.Collection("accounts").FindOne(context.TODO(), bson.D{
		primitive.E{Key: "_id", Value: accountID},
		primitive.E{Key: "token", Value: token},
	}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	file, info, err := r.FormFile("avatar")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"success": false, "errorCode": 2}`)
		return
	}
	defer file.Close()

	ext := filepath.Ext(info.Filename)
	if ext != ".png" && ext != ".jpg" {
		log.Println(err)
		fmt.Fprintf(w, `{"success": false, "errorCode": 3}`)
		return
	}

	dst, err := os.Create("data/avatars/" + account.ID.Hex() + ".png")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"success": false, "errorCode": 4}`)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Println(err)
		fmt.Fprintf(w, `{"success": false, "errorCode": 5}`)
		return
	}

	avatarURL := apiHost + "/users/" + accountIDHex + "/avatar"

	db.Collection("accounts").UpdateOne(context.TODO(), bson.D{
		primitive.E{Key: "token", Value: token},
		primitive.E{Key: "_id", Value: accountID},
	}, bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "avatar", Value: avatarURL},
		}},
	})

	fmt.Fprintf(w, `{"success": true}`)
}