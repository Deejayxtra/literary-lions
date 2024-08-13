package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

// ValidateSession retrieves the user ID associated with a session token and checks if the session is still valid.
// Parameters:
//   - token: The session token to validate.
//
// Returns:
//   - int: The user ID associated with the session if the session is valid.
//   - error: An error if the session is invalid or if any other issue occurs; otherwise, nil.
func ValidateSession(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	// Query to get the user ID and expiration time for the provided token
	err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err // Return 0 and the error if the token is not found or there is a query error
	}

	// Check if the current time is past the expiration time of the session
	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired") // Return 0 and an error if the session has expired
	}

	return userID, nil // Return the user ID and no error if the session is valid
}

// InvalidateSession removes a session from the database using the provided token.
// Parameters:
//   - token: The session token to invalidate.
//
// Returns:
//   - error: An error if the deletion fails; otherwise, nil.
func InvalidateSession(token string) error {
	// Execute a DELETE statement to remove the session with the specified token
	_, err := db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err // Return any error encountered during the deletion
}

// CreateSession generates a new session token and inserts it into the database for the specified user.
// Parameters:
//   - userID: The ID of the user for whom the session is created.
//
// Returns:
//   - string: The generated session token if successful.
//   - error: An error if token generation or insertion fails; otherwise, nil.
func CreateSession(userID int) (string, error) {
	// Generate a new session token
	token, err := generateSessionToken()
	if err != nil {
		return "", err // Return an empty string and the error if token generation fails
	}

	// Define the expiration time for the session (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert the new session into the database
	_, err = db.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		return "", err // Return an empty string and the error if insertion fails
	}

	return token, nil // Return the generated token and no error if insertion is successful
}

// generateSessionToken creates a new random session token.
// Returns:
//   - string: The generated session token if successful.
//   - error: An error if the token generation fails; otherwise, nil.
func generateSessionToken() (string, error) {
	// Define the length of the token in bytes
	const tokenLength = 32

	// Create a byte slice to hold the random bytes
	tokenBytes := make([]byte, tokenLength)

	// Fill the byte slice with random bytes
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err // Return an empty string and the error if random byte generation fails
	}

	// Encode the byte slice to a base64 string
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Return the generated token
	return token, nil
}
