package handlers

import (
	"html/template"
	"log"
	"net/http"
)


var errorTemplate = template.Must(template.ParseFiles("templates/notification.html"))

func StatusInternalServerError(w http.ResponseWriter, message string) {
	// Log the error message for debugging purposes
	log.Print("Message: ", message)

	// Data to be passed to the template
	data := struct {
		Error string
	}{
		Error: message,
	}

	// Render the error template with the provided message
	if err := errorTemplate.Execute(w, data); err != nil {
		// If template rendering fails, log the error and send a generic error response
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}