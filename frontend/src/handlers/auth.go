package handlers

import (
	"net/http"
)

// Helper function to check authentication status
func isAuthenticated(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", false
	}

	// Check if user session exists
	token := cookie.Value
	user, exists := sessionStore.Get(token)
	if !exists {
		return "", false
	}

	return user.Username, true
}