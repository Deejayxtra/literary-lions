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
	"time"

	"literary-lions/frontend/src/models"
)

var authToken string // Need to work on how and where to keep the token

// LoginHandler handles user login and redirects to the conversation room.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login-failure.html"))
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract credentials from form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Sample credentials
	credentials := models.Credentials{
		Email:    email,
		Password: password,
	}

	respChan := make(chan models.AuthResponse, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go SendLoginRequest(credentials, &wg, respChan)

	// Wait for the goroutine to finish
	wg.Wait()
	close(respChan)

	// Get the response
	select {
	case response := <-respChan:
		authToken = response.Token
		if response.Success {
			// Extract room_id query parameter
			roomID := r.URL.Query().Get("room_id")
			if roomID == "" {
				roomID = "category1" // Set a default room ID if not provided
			}

			// Redirect to the specific conversation room after successful login
			redirectURL := fmt.Sprintf("/conversation-room?room_id=%s", roomID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			// // Redirect to the conversation room after successful login
			// http.Redirect(w, r, "/conversation-room", http.StatusSeeOther)
		} else {
			// Log the failure reason
			log.Printf("Login failed: %s", response.Message)

			// Render the login failure template
			data := struct {
				Success bool
				Message string
			}{
				Success: response.Success,
				Message: response.Message,
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				log.Fatalf("Error rendering template: %v", err)
			}
		}
	case <-time.After(10 * time.Second):
		fmt.Println("Timeout while processing request")
	}

}

// SendLoginRequest
func SendLoginRequest(credentials models.Credentials, wg *sync.WaitGroup, respChan chan models.AuthResponse) {

	defer wg.Done()

	// Convert credentials to JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error marshaling credentials: %v", err),
		}
		return
	}

	// Create a POST request
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error creating request: %v", err),
		}
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error sending request: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error reading response: %v", err),
		}
		return
	}

	// Check response status code
	if resp.StatusCode != http.StatusOK {
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

		fmt.Printf("error response: %v", errorMessage)

		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintln(errorMessage),
		}
		return
	}

	// Optionally, you can further process the response body if needed
	var responseMessage map[string]interface{}
	if err := json.Unmarshal(body, &responseMessage); err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error unmarshaling response: %v", err),
		}
		return
	}

	respChan <- models.AuthResponse{
		Success: true,
		Token:   responseMessage["token"].(string), // We need to come back to this and figure out how to keep it for subsequent requests
	}
}
