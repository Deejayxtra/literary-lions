package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
	"log"
	"net/http"
	"time"
)

func ShowUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)

	// Check if the user is authenticated
	if !authenticated {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Retrieve session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user details from session store using the token
	token := cookie.Value
	userData, exists := sessionStore.Get(token)
	if !exists {
		// If the user does not exist in the session store, redirect to login
		message := `You are not authorized! Please <a href="/login">login</a> before accessing your profile.`
		tmpl := template.Must(template.ParseFiles("templates/profile.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Error": template.HTML(message),
		})
		return
	}

	// Handle GET requests to render the profile page
	if r.Method == http.MethodGet {
		data := struct {
			Error	 bool
			Username string
			Email    string
		}{
			Error:    false,
			Username: currentUser,
			Email:    userData.Email,
		}	

		// Render the profile template with the user's data
		RenderTemplate(w, "profile.html", data)
		return
	}

	// If not a GET request, handle it accordingly
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
	currentUser, authenticated := isAuthenticated(r)
	if !authenticated {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Retrieve session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Print("Error retrieving session token: ", err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user details from session store using the token
	token := cookie.Value
	userData, exists := sessionStore.Get(token)
	if !exists {
		// If the user does not exist in the session store, redirect to login
		message := `You are not authorized! Please <a href="/login">login</a> before accessing your profile.`
		tmpl := template.Must(template.ParseFiles("templates/profile.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Error": template.HTML(message),
		})
		return
	}

	// Renders page to the user to update profile
	if r.Method == http.MethodGet {
		data := struct {
			Username string
			Email    string
		}{
			Username: currentUser,
			Email:    userData.Email,
		} 
		
		// Render the profile template with the user's data
		RenderTemplate(w, "profile-update.html", data)
		return
	} 

	// Handle POST request to process the update form submission
	if r.Method == http.MethodPost {
		// Extract credentials from form values
		email := r.FormValue("email")
		username := r.FormValue("username")

		respChan := make(chan models.ResponseDetails, 1)

		// Sample credentials
		credentials := models.User{
			Email:    email,
			Username: username,
		}

		// Calls the function that sends request to the server
		go func() {
			SendUpdateUserProfile(cookie, credentials, respChan)
		}()

		select {
		case response := <-respChan:
			if response.Success {
				// Update the session store with new data
				sessionStore.Set(token, response.Username, response.Email)
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			} else {
				// Pass error message to template
				tmpl := template.Must(template.ParseFiles("templates/profile.html"))
				data := struct {
					Error   template.HTML
					Email   string
					Username string
				}{
					Error:    template.HTML(response.Message),
					Username: currentUser,
					Email:    userData.Email,
				}
				if err := tmpl.Execute(w, data); err != nil {
					http.Error(w, "Error rendering template", http.StatusInternalServerError)
					log.Fatalf("Error rendering template: %v", err)
				}
			}
		case <-time.After(10 * time.Second):
			// Handle the case where the operation times out
			http.Error(w, "Profile update timed out", http.StatusGatewayTimeout)
		}
		return
	}

	// If the method is neither GET nor POST, return an error
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}


func SendUpdateUserProfile(cookie *http.Cookie, user models.User, respChan chan models.ResponseDetails) {
	defer close(respChan) // Ensure the channel is closed once this function completes
	// Mashalls data to JSON format
	jsonData, err := json.Marshal(user)
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error marshaling credentials: %v", err),
		}
		return
	}

	// Define the PUT request
	url := config.BaseApi + "/userprofile-update" // Update endpoint in the server side
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error creating request: %v", err),
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the session cookie in the request
	req.AddCookie(cookie)

	// Send the PUT request to the API
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

	username, usernameOK := responseMessage["username"].(string)
	email, emailOK := responseMessage["email"].(string)
	if !usernameOK || !emailOK {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: "Invalid response from authentication server",
		}
		return
	}

	// Handle successful response
	respChan <- models.ResponseDetails{
		Success: true,
		Message: "Profile updated successfully!",
		Username: username,
		Email:	  email,
	}
}