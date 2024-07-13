package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./literary_lions.db")
	if err != nil {
		return nil, err
	}

	// Adjust path calculation for schema.sql file
	// schemaPath := filepath.Join("..", "..", "internal", "db", "schema.sql")
	schemaPath := filepath.Join("..", "internal", "db", "schema.sql")
	absPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return nil, err
	}

	// Print the path for debugging
	log.Printf("Using schema file at: %s", absPath)

	schema, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}
