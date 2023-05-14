package main

import (
	"context"
	"log"
	"net/http"
	"os"

	routes_chat "github.com/flowee-ru/flowee-api/routes/chat"
	routes_settings "github.com/flowee-ru/flowee-api/routes/settings"
	routes_users "github.com/flowee-ru/flowee-api/routes/users"
	routes_actions "github.com/flowee-ru/flowee-api/routes/actions"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/flowee-ru/flowee-api/ws"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func connectMongo(ctx context.Context) (*mongo.Database, error) {
	mongoUri := "mongodb://localhost:27017"
	if os.Getenv("MONGO_URI") != "" {
		mongoUri = os.Getenv("MONGO_URI")
	}

	mongoDB := "flowee"
	if os.Getenv("MONGO_DB") != "" {
		mongoDB = os.Getenv("MONGO_DB")
	}

	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(mongoDB), nil
}

func main() {
	godotenv.Load()

	port := "8000"
	if os.Getenv("API_PORT") != "" {
		port = os.Getenv("API_PORT")
	}

	db, err := connectMongo(context.TODO())
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
		routes_users.Register(w, r, db)
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

	router.HandleFunc("/users/username/{username}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.GetInfoByUsername(ctx, w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users/{accountID}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.GetInfo(ctx, w, r, db)
	}).Methods("GET")

	router.HandleFunc("/users/{accountID}/avatar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.SetAvatar(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_users.GetUsers(ctx, w, r, db)
	}).Methods("GET")

	// actions
	router.HandleFunc("/users/{accountID}/follow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_actions.Follow(ctx, w, r, db)
	}).Methods("POST")

	router.HandleFunc("/users/{accountID}/unfollow", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_actions.Unfollow(ctx, w, r, db)
	}).Methods("POST")

	// settings
	router.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_settings.GetSettings(ctx, w, r, db)
	}).Methods("GET")

	router.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_settings.UpdateSettings(ctx, w, r, db)
	}).Methods("POST")

	// chat
	router.HandleFunc("/users/{accountID}/chat/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		routes_chat.SendMessage(ctx, w, r, db)
	}).Methods("POST")

	// websocket
	router.HandleFunc("/users/{accountID}/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.Ws(ctx, w, r, db)
	})

	log.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}

func setupCors(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        h.ServeHTTP(w, r)
    })
}