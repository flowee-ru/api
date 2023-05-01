package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/flowee-ru/flowee-api/ws"
	"github.com/flowee-ru/flowee-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ChatSendMessage(w http.ResponseWriter, r *http.Request, db *mongo.Database, ctx context.Context) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	stream := r.FormValue("stream")
	token := r.FormValue("token")
	content := r.FormValue("content")

	if content == "" || token == "" || stream == "" || !primitive.IsValidObjectID(stream) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	streamID, _ := primitive.ObjectIDFromHex(stream)

	acc, err := utils.GetAccountFromToken(ctx, db, token)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	ws.OnChatMessage(int(time.Now().Unix()), acc.Username, acc.Avatar, content, streamID)

	fmt.Fprintf(w, `{"success": true}`)
}