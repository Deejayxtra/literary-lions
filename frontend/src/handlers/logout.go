package handlers

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "html/template"
	// "log"
	"net/http"
	// "sync"
	// "time"
	// "io/ioutil"

	// "literary-lions/frontend/src/models"
)

// Logout handles user logout.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session or JWT token
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	})

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
