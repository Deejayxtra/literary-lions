package db

import (
	"database/sql"

	"os"

	_ "github.com/mattn/go-sqlite3"
	// "log"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./literary_lions.db")
	if err != nil {
		return nil, err
	}

	schema, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}
