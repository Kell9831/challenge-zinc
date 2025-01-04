package main

import (
	"data-backend/handlers"
	"data-backend/middlewares"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)


func main() {

	// Cargar variables desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Cambiar a http.Handle
	http.Handle("/", middlewares.HandleCORS(http.DefaultServeMux))
	http.Handle("/search", middlewares.HandleCORS(http.HandlerFunc(handlers.SearchHandler)))

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


