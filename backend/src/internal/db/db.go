package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./literary_lions.db")
	if err != nil {
		return nil, err
	}

	createTables(db)
	createDefaultAdmin(db)
	return db, nil
}

func createTables(db *sql.DB) {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            username TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            role TEXT NOT NULL CHECK (role IN ('user', 'admin'))
        )`,
		`CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        )`,
		`CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            category_id INTEGER,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (category_id) REFERENCES categories(id)
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
		`CREATE TABLE IF NOT EXISTS likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER,
            comment_id INTEGER,
            user_id INTEGER NOT NULL,
            is_like BOOLEAN NOT NULL,
            FOREIGN KEY (post_id) REFERENCES posts(id),
            FOREIGN KEY (comment_id) REFERENCES comments(id),
            UNIQUE (user_id, post_id, comment_id)
        )`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatalf("Could not create table: %v", err)
		}
	}
}

func createDefaultAdmin(db *sql.DB) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE username = 'admin'").Scan(&username)

	if err == sql.ErrNoRows {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err := db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'admin')", "admin@mail.com", "admin", hashedPassword)
		if err != nil {
			log.Fatalf("Could not create admin user: %v", err)
		}
		log.Println("Default admin user created")
	} else if err != nil {
		log.Fatalf("Could not check admin user: %v", err)
	}
}
