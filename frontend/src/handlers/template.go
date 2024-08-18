package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// RenderTemplate renders the specified HTML template with optional data.
func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + tmplName)
	// tmpl, err := template.ParseFiles("/frontend/src/templates/" + tmplName)  //this is working
	// tmpl, err := template.ParseFiles("./templates/" + tmplName)
	if err != nil {
		message := "Internal Server Error"
		StatusInternalServerError(w, message)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Error handling for render template
	err = tmpl.Execute(w, data)
	if err != nil {
		message := "Internal Server Error"
		StatusInternalServerError(w, message)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
