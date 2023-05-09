package chat

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/flowee-ru/flowee-api/utils"
	"github.com/flowee-ru/flowee-api/ws"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendMessage(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	accountIDHex := mux.Vars(r)["accountID"]
	token := r.FormValue("token")
	content := r.FormValue("content")

	if content == "" || token == "" || accountIDHex == "" || !primitive.IsValidObjectID(accountIDHex) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	accountID, _ := primitive.ObjectIDFromHex(accountIDHex)

	acc, err := utils.GetAccountFromToken(context.TODO(), db, token)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	ws.OnChatMessage(int(time.Now().Unix()), acc.Username, acc.Avatar, content, accountID)

	fmt.Fprintf(w, `{"success": true}`)
}