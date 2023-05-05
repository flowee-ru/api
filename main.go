package main

import (
	"context"
	"log"
	"net/http"
	"os"

	routes_chat "github.com/flowee-ru/flowee-api/routes/chat"
	routes_users "github.com/flowee-ru/flowee-api/routes/users"

	"github.com/flowee-ru/flowee-api/utils"
	"github.com/flowee-ru/flowee-api/ws"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := "8000"
	if os.Getenv("API_PORT") != "" {
		port = os.Getenv("API_PORT")
	}

	db, err := utils.ConnectMongo(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter()

	// users
	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.Login(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.Register(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/verifyAccount", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.VerifyAccount(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/verifyToken", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.VerifyToken(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/resendEmail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.ResendEmail(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/follow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.Follow(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/unfollow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.Unfollow(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.GetInfoByUsername(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users/{accountID}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_users.GetInfo(w, r, db)
	}).Methods("GET")

	// chat
	router.HandleFunc("/users/{accountID}/chat/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		routes_chat.SendMessage(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/chat/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ws.Ws(w, r, db)
	})

	http.Handle("/", router)

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}