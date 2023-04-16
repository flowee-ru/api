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
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	basePath := "/api"
	if os.Getenv("BASE_PATH") == "/" {
		basePath = ""
	} else if os.Getenv("BASE_PATH") != "" {
		basePath = os.Getenv("BASE_PATH")
	}

	db, err := utils.ConnectMongo(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc(basePath + "/auth/login", func(w http.ResponseWriter, r *http.Request) {
		routes.Login(w, r, db)
	})

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}