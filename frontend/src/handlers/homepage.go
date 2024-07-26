package handlers

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"literary-lions/frontend/src/models"
)

var (
	conversations     = make(map[string][]models.Message) // map of roomID to messages
	conversationsLock = sync.RWMutex{}
)

// HomeHandler handles the home page request.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", nil)
}

// RenderTemplate renders the specified HTML template with optional data.
func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
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