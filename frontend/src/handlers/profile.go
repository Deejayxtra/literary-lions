package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
	"net/http"
	"strconv"
	"sync"
	"time"
)


func ShowUserProfile(w http.ResponseWriter, r *http.Request) {

	// Get the authentication status and the currentUser if any
	currentUser, authenticated := isAuthenticated(r)
	if r.Method == http.MethodGet {
		data := struct {
			Username    string
		}{
			Username:      currentUser,
		}
		// Render the template with posts and authentication status
		RenderTemplate(w, "profile.html", data)
	
		return
	}

	if !authenticated {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// // Create a new GET request to fetch the user profile
	// req, err := http.NewRequest("GET", config.BaseApi+"/user/"+userID, nil)
	// if err != nil {
	// 	http.Error(w, "Failed to create request", http.StatusInternalServerError)
	// 	return
	// }

	// // Use an http.Client to make the request
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
	// 	return
	// }
	// defer resp.Body.Close()

	// // Read the response body
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	http.Error(w, "Failed to read response", http.StatusInternalServerError)
	// 	return
	// }

	// // Log the response body for debugging
	// log.Printf("API Response: %s", body)

	// // Parse the JSON response into a User model
	// var user models.User
	// err = json.Unmarshal(body, &user)
	// if err != nil {
	// 	http.Error(w, "Failed to parse response: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// data := struct {
	// 	User          models.User
	// 	Authenticated bool
	// 	Username      string
	// }{
	// 	User:          user,
	// 	Authenticated: authenticated,
	// 	Username:      user.Username,
	// }

	// // Render the template with user profile and authentication status
	// RenderTemplate(w, "profile.html", data)
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