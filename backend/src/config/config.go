package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret   string
	DatabaseDSN string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		JWTSecret:   os.Getenv("JWT_SECRET"),
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}, nil
}
