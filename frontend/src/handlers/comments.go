package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("postID")

	if r.Method == http.MethodPost {
		r.ParseForm()
		content := r.FormValue("content")

		respChan := make(chan models.ResponseDetails, 1)
		var wg sync.WaitGroup

		payload := models.Comment{
			PostID:  postID,
			Content: content,
		}

		// Extract the session cookie from the header
		cookieToken, err := r.Cookie("session_token")
		if err != nil {
			// http.Error(w, "Failed to get session cookie", http.StatusUnauthorized)
			// User must be logged-in to continue
			message := `You are not authorized! Please <a href="/login">login</a> before adding comment.`
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(message),
			})
			return
		}

		wg.Add(1)
		go func() {
			SendAddCommentRequest(cookieToken, payload, &wg, respChan)
		}()

		go func() {
			wg.Wait()
			close(respChan)
		}()

		responseDetails := <-respChan

		if responseDetails.Status == http.StatusCreated {
			log.Print("comment added")
			http.Redirect(w, r, "post?id="+postID, http.StatusSeeOther)
		} else if responseDetails.Status == http.StatusUnauthorized {
			// responseDetails.Status = http.StatusUnauthorized
			responseDetails.Message = `You are not authorized! Please <a href="/login">login</a> before adding comment.`
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": template.HTML(responseDetails.Message),
			})
		} else {
			log.Print("Something went wrong")
			// responseDetails.Status = resp.StatusCode
			responseDetails.Message = "Oops! Something went wrong. Failed to add comment."
			tmpl := template.Must(template.ParseFiles("templates/post.html"))
			tmpl.Execute(w, map[string]interface{}{
				"Error": responseDetails.Message,
			})
		}
	}
}

func SendAddCommentRequest(cookie *http.Cookie, payload models.Comment, waitGroup *sync.WaitGroup, respChan chan models.ResponseDetails) {
	defer waitGroup.Done()

	// Convert payload to JSON
	commentData, err := json.Marshal(payload)
	if err != nil {
		respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to marshal payload"}
		return
	}

	// Create a POST request
	req, err := http.NewRequest("POST", config.BaseApi+"/post/"+payload.PostID+"/comment", bytes.NewBuffer(commentData))
	if err != nil {
		respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Failed to create request"}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the session cookie in the request
	req.AddCookie(cookie)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respChan <- models.ResponseDetails{Status: http.StatusInternalServerError, Message: "Request failed"}
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error reading response: %v", err),
		}
		return
	}

	// Check response status code
	if resp.StatusCode != http.StatusCreated {
		// Attempt to parse the error message from the response
		var errorResponse map[string]interface{}
		var errorMessage string
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			errorMessage = string(body) // Use raw body as fallback
		} else {
			if errMsg, exists := errorResponse["error"]; exists {
				errorMessage = fmt.Sprintf("%v", errMsg)
			} else {
				errorMessage = "unknown error"
			}
		}

		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintln(errorMessage),
			Status:  resp.StatusCode,
		}
		return
	}

	// Optionally, you can further process the response body if needed
	var responseMessage map[string]interface{}
	if err := json.Unmarshal(body, &responseMessage); err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error unmarshaling response: %v", err),
			Status:  resp.StatusCode,
		}
		return
	}

	// Extracting the message from the response map
	message, ok := responseMessage["message"].(string)
	if !ok {
		message = "Unexpected response format"
	}

	respChan <- models.ResponseDetails{
		Success: true,
		Message: fmt.Sprintln(message), // displays server response to the user
		Status:  resp.StatusCode,
	}
}
