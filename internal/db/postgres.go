package db

import (
	"fmt"
	"log"

	"github.com/aditya/finback/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgres establishes a connection to PostgreSQL and returns a GORM DB instance
func ConnectPostgres() *gorm.DB {
	host := config.GetEnv("DB_HOST", "localhost")
	user := config.GetEnv("DB_USER", "postgres")
	password := config.GetEnv("DB_PASSWORD", "postgrespassword")
	dbname := config.GetEnv("DB_NAME", "finance")
	port := config.GetEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL!")
	return db
}
