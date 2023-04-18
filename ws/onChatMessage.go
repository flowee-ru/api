package ws

import (
	"encoding/base64"
	"strconv"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OnChatMessage(timestamp int, authorName string, content string, stream primitive.ObjectID) {
	for _, s := range clients {
		if s == stream {
			authorName = base64.StdEncoding.EncodeToString([]byte(authorName))
			content = base64.StdEncoding.EncodeToString([]byte(content))

			conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(timestamp)+"|"+authorName+"|"+content))
		}
	}
}