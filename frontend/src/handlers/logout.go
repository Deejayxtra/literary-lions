package handlers

import (
	"net/http"
	"sync"
	"literary-lions/frontend/src/config"
	"literary-lions/frontend/src/models"
)

// Logout handles user logout.
func Logout(w http.ResponseWriter, r *http.Request) {
	respChan := make(chan models.AuthResponse, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	// Extract the session cookie from the header
	cookieToken, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Failed to get session cookie", http.StatusUnauthorized)
		return
	}

	// Calls the function that sends request to the server
	go SendLogoutRequest(cookieToken, &wg, respChan)

	// Wait for the goroutine to finish
	wg.Wait()
	close(respChan)
	
	sessionStore.Delete(cookieToken.Value)

	// Defines the cookie
	cookie := http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // This deletes the cookie
	}
	http.SetCookie(w, &cookie)

	// Render the logout.html template
	RenderTemplate(w, "logout.html", nil)
}

// SendLogoutRequest handles user logout from the client side
func SendLogoutRequest(cookie *http.Cookie, wg *sync.WaitGroup, respChan chan models.AuthResponse) {
	defer wg.Done()

	// Create a POST request to the logout endpoint
	req, err := http.NewRequest(http.MethodPost, config.BaseApi+"/logout", nil)
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: "Error creating request",
		}
		return
	}

	// Set the session cookie in the request
	req.AddCookie(cookie)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: "Error sending request",
		}
		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		respChan <- models.AuthResponse{
			Success: false,
			Message: "Logout failed",
		}
		return
	}

	// Respond with a success message
	respChan <- models.AuthResponse{
		Success: true,
		Message: "Successfully logged out",
	}
}
