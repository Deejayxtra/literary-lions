package db

import (
	"context"
	"database/sql"
	"log"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// InitDB initializes the database connection, sets up tables, and creates a default admin user.
// Returns a pointer to the database connection and an error if any.
func InitDB() (*sql.DB, error) {
    // Open a connection to the SQLite database file
	db, err := sql.Open("sqlite3", "./literary_lions.db")
	if err != nil {
		return nil, err // Return the error if the database connection fails
	}

	// Create a context with a timeout to avoid hanging if the database doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure the context is canceled to free up resources

	// Ping the database to verify the connection is alive
	if err := db.PingContext(ctx); err != nil {
		return nil, err // Return the error if the ping fails
	}

	// Create necessary tables if they don't already exist
	createTables(db)

	// Create a default admin user if one doesn't already exist
	createDefaultAdmin(db)

	// Return the database connection
	return db, nil
}

// createTables creates the necessary tables for the application if they don't already exist.
// It accepts a pointer to the database connection.
func createTables(db *sql.DB) {
    // List of SQL statements to create tables
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            username TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            role TEXT NOT NULL CHECK (role IN ('user', 'admin'))
        )`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			UUID VARCHAR(36) NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
    		category TEXT,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`,
		`CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			is_like BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
        )`,
		`CREATE TABLE IF NOT EXISTS post_dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER,
			is_dislike BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id)
        )`,
		`CREATE TABLE IF NOT EXISTS comment_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			comment_id INTEGER,
			is_like BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id)
        )`,
		`CREATE TABLE IF NOT EXISTS comment_dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			comment_id INTEGER,
			is_dislike BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id)
        )`,
	}

	// Iterate over the table creation statements and execute them
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatalf("Could not create table: %v", err) // Log and exit if any table creation fails
		}
	}
}

// createDefaultAdmin creates a default admin user if one does not already exist in the database.
// It accepts a pointer to the database connection.
func createDefaultAdmin(db *sql.DB) {
	var username string

	// Check if the admin user already exists by querying the database
	err := db.QueryRow("SELECT username FROM users WHERE username = 'admin'").Scan(&username)

	// If the admin user does not exist, create one
	if err == sql.ErrNoRows {
		// Generate a hashed password for the admin user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

		// Insert the admin user into the database
		_, err := db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'admin')",
			"admin@mail.com", "admin", hashedPassword)
		if err != nil {
			log.Fatalf("Could not create admin user: %v", err) // Log and exit if the admin creation fails
		}
		log.Println("Default admin user created") // Log success if the admin is created
	} else if err != nil {
		log.Fatalf("Could not check admin user: %v", err) // Log and exit if the user check fails for any other reason
	}
}
