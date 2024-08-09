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
	"strconv"
	"sync"
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

	// Handle GET requests to render the profile page
	if r.Method == http.MethodGet {
		data := struct {
			Username string
			Email    string
		}{
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Convert the user ID from string to integer
	userID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Create a User struct from the form data
	user := models.User{
		ID:         userID,
		Username:   r.FormValue("username"),
		Email:      r.FormValue("email"),
	//	ProfilePic: r.FormValue("profile_pic"),
	}

	// Convert the User struct to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		return
	}

	// Create a new POST request to the API to update the user profile
	req, err := http.NewRequest("POST", config.BaseApi+"/user/update", bytes.NewBuffer(userData))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Use an http.Client to make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Check for errors in the API response
	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	// Redirect to the profile page
	http.Redirect(w, r, "/profile?id="+r.FormValue("id"), http.StatusSeeOther)
}

// DeleteUserProfile handles deleting a user's profile
func DeleteUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the user ID from the form data
	userID := r.FormValue("id")

	// Create a new DELETE request to the API to delete the user profile
	req, err := http.NewRequest("DELETE", config.BaseApi+"/user/"+userID, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Use an http.Client to make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to delete user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check for errors in the API response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	// Redirect to the homepage or a confirmation page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SendProfileRequest(user models.User, wg *sync.WaitGroup, respChan chan models.ResponseDetails) {

	defer wg.Done()

	jsonData, err := json.Marshal(user)
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

}