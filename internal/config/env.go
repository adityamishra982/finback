package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from a .env file if it exists
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on existing environment variables")
	}
}

// GetEnv retrieves the value of an environment variable or a fallback default
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
