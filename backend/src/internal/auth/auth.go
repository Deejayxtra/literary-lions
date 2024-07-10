package auth

import (
    "database/sql"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "literary-lions/backend/src/internal/models"
)

func RegisterUser(db *sql.DB, email, username, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, string(hashedPassword))
    return err
}

func AuthenticateUser(db *sql.DB, email, password string) (models.User, error) {
    var user models.User
    row := db.QueryRow("SELECT id, email, username, password FROM users WHERE email = ?", email)
    err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
    if err != nil {
        return models.User{}, errors.New("user not found")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return models.User{}, errors.New("incorrect password")
    }

    return user, nil
}
