package main

import (
	"context"
	"flowee-api/routes"
	"flowee-api/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var ctx = context.TODO()

func main() {
	godotenv.Load()

	port := "8000"
	if os.Getenv("WS_PORT") != "" {
		port = os.Getenv("WS_PORT")
	}

	basePath := "/api"
	if os.Getenv("WS_BASE_PATH") == "/" {
		basePath = ""
	} else if os.Getenv("WS_BASE_PATH") != "" {
		basePath = os.Getenv("WS_BASE_PATH")
	}

	db, err := utils.ConnectMongo(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter()

	// auth
	router.HandleFunc(basePath + "/auth/login", func(w http.ResponseWriter, r *http.Request) {
		routes.Login(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/register", func(w http.ResponseWriter, r *http.Request) {
		routes.Register(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/verifyAccount", func(w http.ResponseWriter, r *http.Request) {
		routes.VerifyAccount(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/verifyToken", func(w http.ResponseWriter, r *http.Request) {
		routes.VerifyToken(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/auth/resendEmail", func(w http.ResponseWriter, r *http.Request) {
		routes.ResendEmail(w, r, db, ctx)
	})

	// actions
	router.HandleFunc(basePath + "/actions/follow", func(w http.ResponseWriter, r *http.Request) {
		routes.Follow(w, r, db, ctx)
	})
	router.HandleFunc(basePath + "/actions/unfollow", func(w http.ResponseWriter, r *http.Request) {
		routes.Unfollow(w, r, db, ctx)
	})

	http.Handle("/", router)

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}