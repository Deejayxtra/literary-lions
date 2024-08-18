package handlers

import (
	"html/template"
	"literary-lions/frontend/src/models"
	"log"
	"net/http"
	"strings"
	"sync"
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

func UnauthorizedErrorNotification(w http.ResponseWriter, r *http.Request, postID string, message string) {
	respChan := make(chan models.PostDetails, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		SendGetPostByIdRequest(postID, w, r, &wg, respChan)
	}()

	go func() {
		wg.Wait()
		close(respChan)
	}()

	currentUser, authenticated := isAuthenticated(r)
	response := <-respChan
	formattedDate := response.Post.CreatedAt.Format("January 2, 2006 at 3:04pm")
	response.Post.FormattedContent = strings.Split(response.Post.Content, "\n")

	if response.Status == http.StatusOK {

		data := struct {
			Post          models.Post
			FormattedDate string
			Authenticated bool
			Comments      []models.Comment
			Error         template.HTML
			Username      string
			Likes         int
			Dislikes      int
		}{
			Post:          response.Post,
			FormattedDate: formattedDate,
			Authenticated: authenticated,
			Comments:      response.Comments,
			Error:         template.HTML(message),
			Username:      currentUser,
			Likes:         response.Likes,
			Dislikes:      response.Dislikes,
		}

		// Render the template with posts and authentication status
		RenderTemplate(w, "post.html", data)

	}
}
