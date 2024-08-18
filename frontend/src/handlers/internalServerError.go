// package handlers

// import (
// 	"html/template"
// 	"literary-lions/frontend/src/models"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"sync"
// )

// var errorTemplate = template.Must(template.ParseFiles("templates/notification.html"))

// func StatusInternalServerError(w http.ResponseWriter, message string) {
// 	// Log the error message for debugging purposes
// 	log.Print("Message: ", message)

// 	// Data to be passed to the template
// 	data := struct {
// 		Error string
// 	}{
// 		Error: message,
// 	}

// 	// Render the error template with the provided message
// 	if err := errorTemplate.Execute(w, data); err != nil {
// 		// If template rendering fails, log the error and send a generic error response
// 		log.Printf("Error rendering template: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func UnauthorizedErrorNotification(w http.ResponseWriter, r *http.Request, postID string, message string) {
//     respChan := make(chan models.PostDetails, 1)
//     var wg sync.WaitGroup

//     wg.Add(1)
//     go func() {
//         SendGetPostByIdRequest(postID, w, r, &wg, respChan)
//     }()

//     go func() {
//         wg.Wait()
//         close(respChan)
//     }()

// 	currentUser, authenticated := isAuthenticated(r)
//     response := <-respChan
// 	formattedDate := response.Post.CreatedAt.Format("January 2, 2006 at 3:04pm")
// 	response.Post.FormattedContent = strings.Split(response.Post.Content, "\n")

//     if response.Status == http.StatusOK {

//         data := struct {
// 		Post          models.Post
// 		FormattedDate string
// 		Authenticated bool
// 		Comments      []models.Comment
// 		Error         template.HTML
// 		Username      string
// 		Likes		  int
// 		Dislikes      int
// 	}{
// 		Post:          response.Post,
// 		FormattedDate: formattedDate,
// 		Authenticated: authenticated,
// 		Comments:      response.Comments,
// 		Error:         template.HTML(message),
// 		Username:      currentUser,
// 		Likes:		   response.Likes,
// 		Dislikes:      response.Dislikes,
//         }

//         // Render the template with posts and authentication status
//         RenderTemplate(w, "post.html", data)

//     }
// }

package handlers

import (
	"html/template"
	"literary-lions/frontend/src/models"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Define the HTML template directly in the Go file
const errorTemplateHTML = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Notification</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>

<body>
    <header>
        <h1>Literary Lions Forum</h1>
    </header>
    <main>
        {{if .Error}}
        <div class="notification notification-error">
            <p>{{.Error}}</p>
        </div>
        {{end}}
    </main>
    <footer>
        <p>&copy; 2024 Forum</p>
    </footer>
</body>

</html>
`

// Initialize the template
var errorTemplate *template.Template

func init() {
	// Parse the embedded HTML template
	var err error
	errorTemplate, err = template.New("error").Parse(errorTemplateHTML)
	if err != nil {
		log.Fatalf("Error parsing error template: %v", err)
	}
}

// StatusInternalServerError renders an internal server error page
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

// UnauthorizedErrorNotification handles unauthorized error notifications
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
