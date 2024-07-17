package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"literary-lions/frontend/src/models"
)


// HomeHandler handles the home page request.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

// Register handles user registration.
func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        renderTemplate(w, "register.html", nil)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    err = models.CreateUser(user.Email, user.Username, user.Password)
    if err != nil {
        log.Println(err)
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    // Redirect to the login page after successful registration
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login handles user login.
func Login(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)

	// Implement session or JWT token handling here

	w.WriteHeader(http.StatusOK)
}

// Logout handles user logout.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session or JWT token
	// This example assumes you're using sessions

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// PostComment handles post creation.
func PostComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "post-comment.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := models.PostComment(post.CategoryID, post.Title, post.Content, post.UserID); err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreateChannel handles comment creation.
func CreateChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "create-channel.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := models.CreateChannel(comment.PostID, comment.Content, comment.UserID); err != nil {
		log.Printf("Error creating comment: %v", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// renderTemplate renders the specified HTML template with optional data.
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + tmplName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
