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

	fmt.Printf("jsonData: %s\n", jsonData)

	// Define the POST request
	url := "http://localhost:8080/register"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error creating request: %v", err),
		}
		// fmt.Println("Error creating request: \n", err)
		// errChan <- fmt.Errorf("error creating request: %w", err)
		// close(errChan)
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
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("received non-OK response status: %s, body: %s", resp.Status, string(body)), // Edit response here
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

	respChan <- models.ResponseDetails{
		Success: true,
		Message: fmt.Sprintf("User registered successfully. Server response: %v", responseMessage),
	}

	fmt.Println("Response from server:", responseMessage)

	fmt.Printf("response: %v\n", resp.StatusCode)
}

func SendLoginRequest(email, password string) (*http.Response, error) {
	loginData := map[string]string{"email": email, "password": password}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/login-handler", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}


// LoginHandler handles user login and redirects to the conversation room.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	respChan := make(chan models.ResponseDetails, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go fetchForUserLoginAsync(credentials, &wg, respChan)

	// Wait for the goroutine to finish
	wg.Wait()
	close(respChan)

	// Get the response
	responseDetails := <-respChan
	fmt.Printf("Login response: %+v\n", responseDetails)

	// Redirect to the conversation room after successful login
	http.Redirect(w, r, "/conversation-room", http.StatusSeeOther)
}

// fetchForUserLoginAsync
func fetchForUserLoginAsync(credentials models.Credentials, wg *sync.WaitGroup, respChan chan models.ResponseDetails) {

	defer wg.Done()

	// Convert credentials to JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error marshaling credentials: %v", err),
		}
		return
	}

	// Create a POST request
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(jsonData))
	if err != nil {
		respChan <- models.ResponseDetails{
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
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error sending request: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	// Read the response
	var responseDetails models.ResponseDetails
	if err := json.NewDecoder(resp.Body).Decode(&responseDetails); err != nil {
		respChan <- models.ResponseDetails{
			Success: false,
			Message: fmt.Sprintf("error decoding response: %v", err),
		}
		return
	}

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		responseDetails.Success = false
		responseDetails.Message = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
	}

	respChan <- responseDetails
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "Missing room_id", http.StatusBadRequest)
		return
	}

	conversationsLock.Lock()
	conversations[roomID] = append(conversations[roomID], message)
	conversationsLock.Unlock()

	w.WriteHeader(http.StatusCreated)
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