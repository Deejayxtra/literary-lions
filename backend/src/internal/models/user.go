package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	ID       int
	Email    string
	Username string
	Password string
	Role     string
}

// SetDatabase initializes the global database instance with the provided database connection.
// Parameters:
//   - databaseInstance: The database connection instance to be set as the global db variable.
func SetDatabase(databaseInstance *sql.DB) {
	db = databaseInstance
}

// RegisterUser creates a new user record in the database with the provided email, username, and password.
// The password is hashed before being stored.
//
// Parameters:
//   - email: The email address of the user.
//   - username: The username of the user.
//   - password: The plaintext password of the user.
//
// Returns:
//   - error: An error if the user registration fails; otherwise, nil.

// func RegisterUser(email, username, password string) error {
// 	// Hash the password using bcrypt
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err // Return the error if password hashing fails
// 	}

// 	// Insert the new user record into the database
// 	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", email, username, string(hashedPassword))
// 	return err // Return any error encountered during insertion
// }

func RegisterUser(email, username, password string) error {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert the new user record into the database
	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", email, username, string(hashedPassword))

	// Check for uniqueness constraint errors
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return errors.New("username already exists, please use another username")
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return errors.New("email already exists, please use another email")
		}
		return err // Return the error if it's not a uniqueness constraint error
	}

	return nil
}

// AuthenticateUser checks if the provided email and password match a user record in the database.
// It returns the user details if authentication is successful.
//
// Parameters:
//   - email: The email address of the user attempting to log in.
//   - password: The plaintext password provided by the user.
//
// Returns:
//   - *User: A pointer to the User object if authentication is successful; otherwise, nil.
//   - error: An error if the user is not found or if the password is incorrect; otherwise, nil.
func AuthenticateUser(email, password string) (*User, error) {
	user := &User{}

	// Query to retrieve user details based on the provided email
	row := db.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found") // Return an error if the user does not exist
		}
		return nil, err // Return any other error encountered during the query
	}

	// Compare the provided password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password") // Return an error if the password does not match
	}

	return user, nil // Return the user details if authentication is successful
}

// GetUser retrieves a user record from the database based on the provided user ID.
//
// Parameters:
//   - userID: The ID of the user to retrieve.
//
// Returns:
//   - *User: A pointer to the User object if the user is found; otherwise, nil.
//   - error: An error if the user is not found or if any other issue occurs; otherwise, nil.
func GetUser(userID int) (*User, error) {
	user := &User{}

	// Query to retrieve user details based on the provided user ID
	row := db.QueryRow("SELECT id, email, username, password FROM users WHERE id = ?", userID)

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found") // Return an error if the user does not exist
		}
		return nil, err // Return any other error encountered during the query
	}

	return user, nil // Return the user details if found
}

// ProfileUpdate updates the email and username for a user with the specified user ID.
//
// Parameters:
//   - userID: The ID of the user whose profile is being updated.
//   - email: The new email address for the user.
//   - username: The new username for the user.
//
// Returns:
//   - error: An error if the profile update fails; otherwise, nil.

func ProfileUpdate(userID int, email, username string) error {
	// Prepare the SQL statement for updating the user profile
	stmt, err := db.Prepare("UPDATE users SET email = ?, username = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %w", err)
	}
	defer stmt.Close()

	// Execute the SQL statement to update the user profile
	_, err = stmt.Exec(email, username, userID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return errors.New("username already exists, please use another username")
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return errors.New("email already exists, please use another email")
		}
		return fmt.Errorf("failed to execute SQL statement: %w", err)
	}

	return nil // Return nil if the profile update is successful
}

func FindUserByEmail(tx *sql.Tx, email string) (*User, error) {
	var user User
	err := tx.QueryRow("SELECT id, email, username, password FROM users WHERE LOWER(email) = LOWER(?)", email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		// Log the error for debugging
		fmt.Printf("Error querying user by email: %v\n", err)
		return nil, err
	}

	// Log user details for debugging
	// fmt.Printf("User found: ID=%d, Email=%s\n", user.ID, user.Email)
	fmt.Printf("User found: ID=%d, Email=%s\n, Username=%s\n", user.ID, user.Email, user.Username)
	// return &user, nil
	return &user, nil
}

