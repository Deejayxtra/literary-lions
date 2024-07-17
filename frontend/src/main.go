package main

import (
    "log"
    "net/http"
    "literary-lions/frontend/src/handlers"
)

func main() {
	
	// Define your handlers for routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/post-comment", handlers.PostComment)
	http.HandleFunc("/create-channel", handlers.CreateChannel)

	// Start the server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}