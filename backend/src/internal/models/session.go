package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ValidateSession retrieves the user ID associated with a session UUID and checks if the session is still valid.
// Parameters:
//   - sessionUUID: The session UUID to validate.
//
// Returns:
//   - int: The user ID associated with the session if the session is valid.
//   - error: An error if the session is invalid or if any other issue occurs; otherwise, nil.
func ValidateSession(sessionUUID string) (int, error) {
	var userID int
	var expiresAt time.Time

	// Query to get the user ID and expiration time for the provided session UUID
	err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE uuid = ?", sessionUUID).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err // Return 0 and the error if the UUID is not found or there is a query error
	}

	// Check if the current time is past the expiration time of the session
	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired") // Return 0 and an error if the session has expired
	}

	return userID, nil // Return the user ID and no error if the session is valid
}

// InvalidateSession removes a session from the database using the provided UUID.
// Parameters:
//   - sessionUUID: The session UUID to invalidate.
//
// Returns:
//   - error: An error if the deletion fails; otherwise, nil.
func InvalidateSession(sessionUUID string) error {
	// Execute a DELETE statement to remove the session with the specified UUID
	_, err := db.Exec("DELETE FROM sessions WHERE uuid = ?", sessionUUID)
	return err // Return any error encountered during the deletion
}

// CreateSession generates a new session UUID and inserts it into the database for the specified user.
// Parameters:
//   - userID: The ID of the user for whom the session is created.
//
// Returns:
//   - string: The generated session UUID if successful.
//   - error: An error if UUID generation or insertion fails; otherwise, nil.
func CreateSession(userID int) (string, error) {
	// Generate a new session UUID
	sessionUUID, err := generateSessionUUID()
	if err != nil {
		return "", err // Return an empty string and the error if UUID generation fails
	}

	// Define the expiration time for the session (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert the new session into the database
	_, err = db.Exec("INSERT INTO sessions (user_id, uuid, expires_at) VALUES (?, ?, ?)", userID, sessionUUID, expiresAt)
	if err != nil {
		return "", err // Return an empty string and the error if insertion fails
	}

	return sessionUUID, nil // Return the generated UUID and no error if insertion is successful
}

// generateSessionUUID creates a new UUID for the session.
// Returns:
//   - string: The generated session UUID if successful.
//   - error: An error if the UUID generation fails; otherwise, nil.
func generateSessionUUID() (string, error) {
	// Generate a new UUID
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err // Return an empty string and the error if UUID generation fails
	}

	// Return the string representation of the UUID
	return newUUID.String(), nil
}
