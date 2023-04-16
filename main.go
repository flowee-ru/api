package main

import (
	"context"
	"flowee-api/routes"
	"flowee-api/utils"
	"log"
	"net/http"
	"os"

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

	// auth
	http.HandleFunc(basePath + "/auth/login", func(w http.ResponseWriter, r *http.Request) {
		routes.Login(w, r, db, ctx)
	})
	http.HandleFunc(basePath + "/auth/register", func(w http.ResponseWriter, r *http.Request) {
		routes.Register(w, r, db, ctx)
	})

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}