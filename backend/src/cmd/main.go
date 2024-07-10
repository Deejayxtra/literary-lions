package main

import (
	"literary-lions/backend/src/internal/db"
	"literary-lions/backend/src/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	handlers.InitHandlers(database)

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Add more routes for posts, comments, likes, etc.

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
