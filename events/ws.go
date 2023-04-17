package events

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Ws(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	stream := r.URL.Query().Get("stream")

	if !primitive.IsValidObjectID(stream) {
		return
	}

	streamID, _ = primitive.ObjectIDFromHex(stream)

	for {
		msgType, _, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			break
		}
	}
}