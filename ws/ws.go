package ws

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var clients = make(map[*websocket.Conn] primitive.ObjectID)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ws(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	accountIDHex := mux.Vars(r)["accountID"]

	if !primitive.IsValidObjectID(accountIDHex) || accountIDHex == "" {
		return
	}
	
	accountID, _ := primitive.ObjectIDFromHex(accountIDHex)

	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: accountID}}).Decode(nil)
	if err == mongo.ErrNoDocuments {
		return
	}

	clients[conn] = accountID
	defer delete(clients, conn)

	for {
		msgType, _, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			break
		}
	}
}