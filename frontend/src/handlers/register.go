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

	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
)

// Register handles user registration.
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "register.html", nil)
		return
	} else if r.Method == http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/registration-status.html"))
		// Extract credentials from form values
		email := r.FormValue("email")
		password := r.FormValue("password")
		username := r.FormValue("username")

		// Print credentials for debugging
		fmt.Printf("Credentials: email=%s, password=%s\n, username=%s\n", email, password, username)

		respChan := make(chan models.ResponseDetails, 1)
		var wg sync.WaitGroup

		// Sample credentials
		credentials := models.Credentials{
			Email:    email,
			Username: username,
			Password: password,
		}

		wg.Add(1)
		go func() {
			SendRegistrationRequest(credentials, &wg, respChan)
		}()

		// Wait for the goroutine to finish
		go func() {
			wg.Wait()
			close(respChan)
		}()

		var responseDetails models.ResponseDetails
		select {
		case responseDetails = <-respChan:
			// Handle response details
			fmt.Println("Response:", responseDetails.Message)
		case <-time.After(10 * time.Second):
			responseDetails = models.ResponseDetails{
				Success: false,
				Message: "timeout while processing request",
			}
		}

		data := struct {
			Success bool
			Message string
		}{
			Success: responseDetails.Success,
			Message: responseDetails.Message,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			log.Fatalf("Error rendering template: %v", err)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// SendRegistrationRequest
func SendRegistrationRequest(credentials models.Credentials, wg *sync.WaitGroup, respChan chan models.ResponseDetails) {

	defer wg.Done() // Notify the wait group when this goroutine completes

	// Marshal the user object to JSON.
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error marshaling credentials: %v", err),
		}
		return
	}

	// Define the POST request
	url := config.BaseApi + "/register"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error creating request: %v", err),
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error sending request: %v", err),
		}
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

		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintln(errorMessage),
		}
		return
	}

	// Optionally, you can further process the response body if needed
	var responseMessage map[string]interface{}
	if err := json.Unmarshal(body, &responseMessage); err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error unmarshaling response: %v", err),
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
	}
}
