package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/flowee-ru/flowee-api/routes"
	"github.com/flowee-ru/flowee-api/utils"
	"github.com/flowee-ru/flowee-api/ws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var ctx = context.TODO()

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	godotenv.Load()

	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	db, err := utils.ConnectMongo(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter()

	// auth
	router.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.Login(w, r, db, ctx)
	})
	router.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.Register(w, r, db, ctx)
	})
	router.HandleFunc("/auth/verifyAccount", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.VerifyAccount(w, r, db, ctx)
	})
	router.HandleFunc("/auth/verifyToken", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.VerifyToken(w, r, db, ctx)
	})
	router.HandleFunc("/auth/resendEmail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.ResendEmail(w, r, db, ctx)
	})

	// actions
	router.HandleFunc("/actions/follow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.Follow(w, r, db, ctx)
	})
	router.HandleFunc("/actions/unfollow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.Unfollow(w, r, db, ctx)
	})

	// data fetching
	router.HandleFunc("/users/getInfo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.GetUserInfo(w, r, db)
	})

	// chat
	router.HandleFunc("/chat/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes.ChatSendMessage(w, r, db, ctx)
	})

	// websocket events
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ws.Ws(wsUpgrader, w, r)
	})

	http.Handle("/", router)

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}