package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"literary-lions/frontend/src/models"
)

// Logout handles user logout.
func Logout(w http.ResponseWriter, r *http.Request) {
	respChan := make(chan models.AuthResponse, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	cookieStr := "session_token=" + authToken

	go SendLogoutRequest(cookieStr, &wg, respChan)

	// Wait for the goroutine to finish
	wg.Wait()
	close(respChan)

	// Redirect to the login page after logout
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ConvertCookieString converts a cookie string to *http.Cookie
func ConvertCookieString(cookieStr string) (*http.Cookie, error) {
	parts := strings.SplitN(cookieStr, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cookie format")
	}
	return &http.Cookie{
		Name:  parts[0],
		Value: parts[1],
	}, nil
}

// SendLogoutRequest handles user logout from the client side
func SendLogoutRequest(cookieStr string, wg *sync.WaitGroup, respChan chan models.AuthResponse) {
    defer wg.Done()
	cookie, err := ConvertCookieString(cookieStr)
	if err != nil {
		respChan <- models.AuthResponse{
			Success: false,
			Message: fmt.Sprintf("error converting cookie string: %v", err),
		}
		return
	}
        // Create a POST request to the logout endpoint
    req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/logout", nil)
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