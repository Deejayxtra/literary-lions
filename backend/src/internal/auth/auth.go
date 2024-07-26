// package auth

// import (
// 	"database/sql"
// 	"errors"
// 	"literary-lions/backend/src/internal/models"

// 	"golang.org/x/crypto/bcrypt"
// )

// func RegisterUser(db *sql.DB, email, username, password string) error {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", email, username, string(hashedPassword))
// 	return err
// }

// func AuthenticateUser(db *sql.DB, email, password string) (models.User, error) {
// 	var user models.User
// 	row := db.QueryRow("SELECT id, email, username, password, role FROM users WHERE email = ?", email)
// 	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role)
// 	if err != nil {
// 		return models.User{}, errors.New("user not found")
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		return models.User{}, errors.New("incorrect password")
// 	}

// 	return user, nil
// }

package auth

import (
	"database/sql"
	"errors"
	"literary-lions/backend/src/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, email, username, password string) error {
	// Check if the email already exists
	var existingEmail string
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
	if err == nil {
		return errors.New("email already in use")
	} else if err != sql.ErrNoRows {
		return err
	}

	// Check if the username already exists
	var existingUsername string
	err = db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
	if err == nil {
		return errors.New("username already in use")
	} else if err != sql.ErrNoRows {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", email, username, string(hashedPassword))
	return err
}

// AuthenticateUser checks the email and password for login.
func AuthenticateUser(db *sql.DB, email, password string) (models.User, error) {
	var user models.User
	row := db.QueryRow("SELECT id, email, username, password, role FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return models.User{}, errors.New("email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("incorrect password")
	}

	return user, nil
}

// GetUserIDByUsername retrieves user ID by username.
func GetUserIDByUsername(db *sql.DB, username string) (int, error) {
	var userID int
	query := "SELECT id FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
