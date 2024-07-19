package main

import (
    "log"
    "net/http"
    "literary-lions/frontend/src/handlers"
	"fmt"
)

func main() {
    
    // Define a test route to process a login request
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            // Render the login form HTML
			handlers.RenderTemplate(w, "login.html", nil)
			return
            // tmpl, err := template.ParseFiles("templates/login_form.html")
            // if err != nil {
            //     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            //     log.Println("Error parsing template:", err)
            //     return
            // }
            // err = tmpl.Execute(w, nil)
            // if err != nil {
            //     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            //     log.Println("Error executing template:", err)
            //     return
            // }
        } else if r.Method == http.MethodPost {
            // Extract credentials from form values
            email := r.FormValue("email")
            password := r.FormValue("password")

            // Print credentials for debugging
            fmt.Printf("Credentials: email=%s, password=%s\n", email, password)

            // Call sendLoginRequest to process the login
            resp, err := handlers.SendLoginRequest(email, password)
            if err != nil {
                http.Error(w, "Failed to send login request", http.StatusInternalServerError)
                return
            }
            defer resp.Body.Close()

            // Write response status and body
            w.WriteHeader(resp.StatusCode)
            _, err = w.Write([]byte("Login request processed. Status: " + resp.Status))
            if err != nil {
                http.Error(w, "Failed to write response", http.StatusInternalServerError)
            }
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
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

    // Start the server
    log.Println("Server started on :8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
