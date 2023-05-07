package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var clients = make(map[*websocket.Conn] primitive.ObjectID)
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func Ws(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close()

	accountIDHex := mux.Vars(r)["accountID"]

	if !primitive.IsValidObjectID(accountIDHex) || accountIDHex == "" {
		return
	}
	
	accountID, _ := primitive.ObjectIDFromHex(accountIDHex)

	err = db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: accountID}}).Decode(nil)
	if err == mongo.ErrNoDocuments {
		return
	}

	clients[conn] = accountID
	defer delete(clients, conn)

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if string(p) == "ping" {
			err := conn.WriteMessage(mt, []byte("pong"))
			if err != nil {
				return
			}
		}
	}
}