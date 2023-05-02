package ws

import (
	"encoding/base64"
	"strconv"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OnChatMessage(timestamp int, authorName string, avatar string, content string, stream primitive.ObjectID) {
	for conn, s := range clients {
		if s == stream {
			authorNameEncoded := base64.StdEncoding.EncodeToString([]byte(authorName))
			contentEncoded := base64.StdEncoding.EncodeToString([]byte(content))

			conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(timestamp)+"|"+authorNameEncoded+"|"+avatar+"|"+contentEncoded))
		}
	}
}