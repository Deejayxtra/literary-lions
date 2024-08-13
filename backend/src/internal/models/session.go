package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

// ValidateSession retrieves the user ID from a session token
func ValidateSession(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}

// InvalidateSession removes or marks the session as invalid
func InvalidateSession(token string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

func CreateSession(userID int) (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = db.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateSessionToken() (string, error) {
	// Define the length of the token in bytes
	const tokenLength = 32

	// Create a byte slice to hold the random bytes
	tokenBytes := make([]byte, tokenLength)

	// Fill the byte slice with random bytes
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a base64 string
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Return the generated token
	return token, nil
}
