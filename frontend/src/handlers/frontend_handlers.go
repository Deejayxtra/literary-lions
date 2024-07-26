package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
	"io/ioutil"

	"literary-lions/frontend/src/models"
)

var (
	conversations     = make(map[string][]models.Message) // map of roomID to messages
	conversationsLock = sync.RWMutex{}
)

var authToken string // Need to work on how and where to keep the token

// HomeHandler handles the home page request.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", nil)
}

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
		// errChan := make(chan error, 1)
		var wg sync.WaitGroup

		// Sample credentials
		credentials := models.Credentials {
			Email:    email,
			Username: username,
			Password: password,
		}

		wg.Add(1)
		go func() {
			fetchForUserRegisterAsync(credentials, &wg, respChan)
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

// fetchForUserRegisterAsync
func fetchForUserRegisterAsync(credentials models.Credentials, wg *sync.WaitGroup, respChan chan models.ResponseDetails) {

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
	url := "http://localhost:8080/register"
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

	// Print credentials for debugging
	fmt.Printf("Credentials: email=%s, password=%s\n", email, password)

	// Sample credentials
	credentials := models.Credentials {
		Email:    email,
		Password: password,
	}

	respChan := make(chan models.LoginResponse, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go fetchForUserLoginAsync(credentials, &wg, respChan)

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
				roomID = "channel1" // Set a default room ID if not provided
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

// fetchForUserLoginAsync
func fetchForUserLoginAsync(credentials models.Credentials, wg *sync.WaitGroup, respChan chan models.LoginResponse) {

	defer wg.Done()

	// Convert credentials to JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		respChan <- models.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("error marshaling credentials: %v", err),
		}
		return
	}

	// Create a POST request
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.LoginResponse{
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
		respChan <- models.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("error sending request: %v", err),
		}
		return
	}
	defer resp.Body.Close()


	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respChan <- models.LoginResponse{
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

		respChan <- models.LoginResponse{
			Success: false,
			Message: fmt.Sprintln(errorMessage),
		}
		return
	}

	// Optionally, you can further process the response body if needed
	var responseMessage map[string]interface{}
	if err := json.Unmarshal(body, &responseMessage); err != nil {
		respChan <- models.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("error unmarshaling response: %v", err),
		}
		return
	}

	respChan <- models.LoginResponse{
		Success: true,
		Token: responseMessage["token"].(string), // We need to come back to this and figure out how to keep it for subsequent requests
	}
}


// Logout handles user logout.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear session or JWT token
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	})

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ConversationRoom handles the conversation room.
func ConversationRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			http.Error(w, "Missing room_id", http.StatusBadRequest)
			return
		}

		conversationsLock.Lock()
		messages := conversations[roomID]
		conversationsLock.Unlock()

		data := struct {
			RoomID   string
			RoomName string
			Messages []models.Message
		}{
			RoomID:   roomID,
			RoomName: getRoomName(roomID), // Function to get the room name based on roomID
			Messages: messages,
		}
		RenderTemplate(w, "conversation-room.html", data)
		return
	} else if r.Method == http.MethodPost {
		// Extract Message from form values
		content := r.FormValue("content")

		// Sample Message
		comment := models.Message{
			Content: content,
		}

		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			http.Error(w, "Missing room_id", http.StatusBadRequest)
			return
		}

		conversationsLock.Lock()
		conversations[roomID] = append(conversations[roomID], comment)
		conversationsLock.Unlock()

		// Redirect to the same conversation room to display the updated conversation
		redirectURL := fmt.Sprintf("/conversation-room?room_id=%s", roomID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Helper function to get room name based on roomID
func getRoomName(roomID string) string {
	roomNames := map[string]string{
		"channel1": "General",
		"channel2": "News",
		"channel3": "Entertainment",
		"channel4": "Music",
		"channel5": "Sports",
		"channel6": "Random",
	}
	if name, ok := roomNames[roomID]; ok {
		return name
	}
	return "Unknown Room"
}


// RenderTemplate renders the specified HTML template with optional data.
func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + tmplName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}