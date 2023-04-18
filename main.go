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

	basePath := "/api"
	if os.Getenv("API_BASE_PATH") == "/" {
		basePath = ""
	} else if os.Getenv("API_BASE_PATH") != "" {
		basePath = os.Getenv("API_BASE_PATH")
	}

	wsPath := "/ws"
	if os.Getenv("WS_PATH") == "/" {
		wsPath = ""
	} else if os.Getenv("WS_PATH") != "" {
		wsPath = os.Getenv("WS_PATH")
	}

	db, err := utils.ConnectMongo(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter()

	// auth
	router.HandleFunc(basePath + "/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.Login(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.Register(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/verifyAccount", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.VerifyAccount(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/verifyToken", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.VerifyToken(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/resendEmail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.ResendEmail(w, r, db, ctx)
	})

	// actions
	router.HandleFunc(basePath + "/actions/follow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.Follow(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/actions/unfollow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.Unfollow(w, r, db, ctx)
	})

	// chat
	router.HandleFunc(basePath + "/chat/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes.ChatSendMessage(w, r, db, ctx)
	})

	// websocket events
	router.HandleFunc(wsPath, func(w http.ResponseWriter, r *http.Request) {
		ws.Ws(wsUpgrader, w, r)
	})

	http.Handle("/", router)

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}