
package models

import (
	"database/sql"
	"errors"
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

func SetDatabase(databaseInstance *sql.DB) {
	db = databaseInstance
}

func RegisterUser(email, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", email, username, string(hashedPassword))
	return err
}

func AuthenticateUser(email, password string) (*User, error) {
	user := &User{}
	row := db.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func GetUser(userID int) (*User, error) {
	user := &User{}

	row := db.QueryRow("SELECT id, email, username, password FROM users WHERE id = ?", userID)

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}