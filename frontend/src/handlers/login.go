package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
)

// LoginHandler handles user login and redirects appropriately.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the login form HTML
		RenderTemplate(w, "login.html", nil)
		return
	} else {

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

		// Calls the function that sends request to the server
		go SendLoginRequest(credentials, &wg, respChan)

		// Wait for the goroutine to finish
		wg.Wait()
		close(respChan)

		// Get the response
		select {
		case response := <-respChan:
			if response.Success {

				// Set the session token as a cookie
				expiration := time.Now().Add(24 * time.Hour)
				cookie := http.Cookie{
					Name:     "session_token",
					Value:    response.Token,
					Expires:  expiration,
					HttpOnly: true,
				}
				http.SetCookie(w, &cookie)

				// Keep the token and username in store for later usage
				sessionStore.Set(response.Token, response.Username, response.Email)

				// Redirect to the index page after successful login
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				// Display error notification in case of login failure
				tmpl := template.Must(template.ParseFiles("templates/login.html"))
				tmpl.Execute(w, map[string]interface{}{
					"Error": template.HTML(response.Message),
				})
			}
		case <-time.After(10 * time.Second):
			fmt.Println("Timeout while processing request")
		}

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
	req, err := http.NewRequest(http.MethodPost, config.BaseApi+"/login", bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error creating request: %v", err),
		}
		return
	}

	// Request Header
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

	token, tokenOK := responseMessage["token"].(string)
	username, usernameOK := responseMessage["username"].(string)
	email, emailOK := responseMessage["email"].(string)
	if !tokenOK || !usernameOK || !emailOK {
		respChan <- models.AuthResponse{
			Success: false,
			Message: "Invalid response from authentication server",
		}
		return
	}

	respChan <- models.AuthResponse{
		Success:  true,
		Token:    token,
		Username: username,
		Email:	  email,
	}
}
