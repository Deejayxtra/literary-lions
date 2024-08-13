package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

// Config holds application configuration values.
// JWTSecret: Secret key used for signing JWT tokens.
// DatabaseDSN: Data Source Name for connecting to the database.
type Config struct {
	JWTSecret   string // Secret key for JWT authentication
	DatabaseDSN string // Data Source Name for database connection
}

// LoadConfig loads configuration values from environment variables and returns a Config struct.
// It also handles loading the .env file if present.
//
// Returns:
//   - *Config: A pointer to a Config struct containing the loaded configuration values.
//   - error: An error if there was an issue loading the .env file.
func LoadConfig() (*Config, error) {
	// Load environment variables from a .env file if it exists
	err := godotenv.Load()
	if err != nil {
		// Log an error message if the .env file could not be loaded
		// This message will be visible in the application logs
		log.Printf("Error loading .env file: %v", err)
	}

	// Create and return a Config struct populated with values from environment variables
	// os.Getenv retrieves the value of the environment variable specified
	return &Config{
		JWTSecret:   os.Getenv("JWT_SECRET"),   // JWT secret key for token generation and verification
		DatabaseDSN: os.Getenv("DATABASE_DSN"), // Data Source Name for database connection
	}, nil
}
