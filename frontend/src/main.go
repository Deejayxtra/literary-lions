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
	http.HandleFunc("/post/like", handlers.LikePost)
	http.HandleFunc("/post/dislike", handlers.DislikePost)
	http.HandleFunc("/comment", handlers.AddComment)
	http.HandleFunc("/profile", handlers.ShowUserProfile)
	http.HandleFunc("/update_profile", handlers.UpdateUserProfile)
	http.HandleFunc("/delete_profile", handlers.DeleteUserProfile)

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout-handler", handlers.Logout)
	http.HandleFunc("/create-post", handlers.CreatePost)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
