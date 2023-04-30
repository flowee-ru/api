package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clients = make(map[*websocket.Conn] primitive.ObjectID)
var conn *websocket.Conn

func Ws(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request) {
	conn, _ = upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	stream := r.URL.Query().Get("stream")

	if !primitive.IsValidObjectID(stream) || stream == "" {
		conn.Close()
	}

	streamID, _ := primitive.ObjectIDFromHex(stream)

	clients[conn] = streamID
	defer delete(clients, conn)

	for {
		msgType, _, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			break
		}
	}
}