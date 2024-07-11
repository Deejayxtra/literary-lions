package handlers

import (
	"database/sql"
	"encoding/json"
	"literary-lions/backend/src/internal/auth"
	"literary-lions/backend/src/internal/models"
	"net/http"
)

var db *sql.DB

func InitHandlers(database *sql.DB) {
	db = database
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = auth.RegisterUser(db, user.Email, user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := auth.AuthenticateUser(db, credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: user.Username,
		Path:  "/",
	})

	json.NewEncoder(w).Encode(user)
}

// Implement more handlers for creating posts, comments, liking/disliking, filtering, etc.
