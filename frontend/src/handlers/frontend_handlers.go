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

	"github.com/dgrijalva/jwt-go"
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

		// credChan := make(chan models.Credentials)
		errChan := make(chan error)

		// Sample credentials
		credentials := models.Credentials {
			Email:    email,
			Username: username,
			Password: password,
		}

		go func() {
			fetchForUserRegisterAsync(credentials, errChan)
		}()

		// Handle errors if any
		var err error
		select {
		case err := <-errChan:
			fmt.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			// Pass data to the template
			data := struct {
				Success bool
				Message string
			}{
				Success: err == nil,
				Message: "",
			}
			if err != nil {
				data.Message = "Registration failed: " + err.Error()
			} else {
				data.Message = "User registered successfully"
			}
			tmpl.Execute(w, data)
		}

		// // Call sendLoginRequest to process the login
		// resp, err := handlers.SendLoginRequest(email, password)
		// if err != nil {
		//     http.Error(w, "Failed to send login request", http.StatusInternalServerError)
		//     return
		// }
		// defer resp.Body.Close()

		// // Write response status and body
		// w.WriteHeader(resp.StatusCode)
		// _, err = w.Write([]byte("Login request processed. Status: " + resp.Status))
		// if err != nil {
		//     http.Error(w, "Failed to write response", http.StatusInternalServerError)
		// }
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// // Sample credentials
	// credential := models.Credentials {
	// 	Email:    "test@example.com",
	// 	Username: "testuser",
	// 	Password: "password123",
	// }

	// // Marshal the user object to JSON.
	// jsonData, err := json.Marshal(credential)
	// if err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// 	http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
	// 	return
	// }

	// // Create a new HTTP request.
	// req, err := http.NewRequest("POST", "http://localhost:8000/register", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	fmt.Println("Error creating HTTP request:", err)
	// 	http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
	// 	return
	// }

	// // Set the Content-Type header to application/json.
	// req.Header.Set("Content-Type", "application/json")

	// // client := &http.Client{}
	// // resp, err := client.Do(req)
	// // if err != nil {
	// // 	fmt.Println("Error sending HTTP request:", err)
	// // 	http.Error(w, "Error sending HTTP request", http.StatusInternalServerError)
	// // 	return
	// // }
	// // defer resp.Body.Close()

	// // if resp.StatusCode != http.StatusOK {
	// // 	http.Error(w, "Error registering user", resp.StatusCode)
	// // 	return
	// // }

	// Redirect to the login page after successful registration
	//http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// fetchForUserRegisterAsync
func fetchForUserRegisterAsync(credentials models.Credentials, errChan chan error) {
	
	// Marshal the user object to JSON.
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		errChan <- fmt.Errorf("error marshaling credentials: %w", err)
		close(errChan)
		return
	}

	fmt.Printf("jsonData: %s\n", jsonData)

	// Define the POST request
	url := "http://localhost:8080/register"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		errChan <- fmt.Errorf("error creating request: %w", err)
		close(errChan)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("error sending request: %w", err)
		close(errChan)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("received non-OK response status: %s", resp.Status)
		close(errChan)
		return
	}

	fmt.Printf("response: %v\n", resp.StatusCode)
}

// // Login displays the login page.
// func Login(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		RenderTemplate(w, "login.html", nil)
// 		return
// 	}

// 	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// }

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

	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := models.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		fmt.Printf("Invalid credentials: %s\n", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Printf("Provided password: %s\n", credentials.Password)
	fmt.Printf("Stored password hash: %s\n", user.Password)

	if !user.CheckPassword(credentials.Password) {
		fmt.Printf("password error: %s\n", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	sessionToken, err := generateJWTToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to the conversation room after successful login
	http.Redirect(w, r, "/conversation-room", http.StatusSeeOther)
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

// generateJWTToken generates a JWT token for the given user.
func generateJWTToken(user *models.User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with a secret
	tokenString, err := token.SignedString([]byte("$2a$12$wJ89JKZa/nH/jf/Y0BZhKuGrOq1BF9N6ZOHYpDkqI9lRdfq9nWJ.e"))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
