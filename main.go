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

var ctx = context.Background()

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

	router.Use(setupCors)

	// users
	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.Login(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.Register(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/verifyAccount", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.VerifyAccount(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/verifyToken", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.VerifyToken(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/resendEmail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.ResendEmail(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/follow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.Follow(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/unfollow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.Unfollow(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.GetInfoByUsername(ctx, w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users/{accountID}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.GetInfo(ctx, w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.GetUsers(ctx, w, r, db)
	}).Methods("GET")

	// chat
	router.HandleFunc("/users/{accountID}/chat/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_chat.SendMessage(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/chat/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Ws(ctx, w, r, db)
	})

	http.Handle("/", router)

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}

func setupCors(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        h.ServeHTTP(w, r)
    })
}