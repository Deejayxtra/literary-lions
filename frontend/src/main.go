package main

import (
    "log"
    "net/http"
    "literary-lions/frontend/src/handlers"
)

func main() {
    
    // Define a test route to process a login request
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        // Render the login form HTML
        handlers.RenderTemplate(w, "login.html", nil)
    })

    // Define your handlers for routes
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/register", handlers.Register)
   	// http.HandleFunc("/login", handlers.Login)
    http.HandleFunc("/login-handler", handlers.LoginHandler) // New handler for processing login
    http.HandleFunc("/logout", handlers.Logout)
   // http.HandleFunc("/post-comment", handlers.PostComment)
   // http.HandleFunc("/create-channel", handlers.CreateChannel)
    http.HandleFunc("/conversation-room", handlers.ConversationRoom)
    //http.HandleFunc("/conversation-room", handlers.ShowPosts)
    http.HandleFunc("/create-post", handlers.CreatePost)

    // Start the server
    log.Println("Server started on :8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
