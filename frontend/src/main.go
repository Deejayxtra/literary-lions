package main

import (
	"literary-lions/frontend/src/handlers"
	"log"
	"net/http"
)

func main() {

	// Define your handlers for routes

	http.HandleFunc("/", handlers.ShowPosts)
	http.HandleFunc("/posts/category", handlers.ShowPostsByCategory)
	http.HandleFunc("/post", handlers.ShowPostByID)
	http.HandleFunc("/comment", handlers.AddComment)

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout-handler", handlers.Logout)
	http.HandleFunc("/create-post", handlers.CreatePost)
	http.HandleFunc("/conversation-room", handlers.ConversationRoom)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
