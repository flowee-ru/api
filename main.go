package main

import (
	"flowee-api/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	http.HandleFunc("/", routes.Login)

	log.Println("Starting server on port " + port)
	http.ListenAndServe(":" + port, nil)
}