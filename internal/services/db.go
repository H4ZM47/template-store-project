package services

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	if os.Getenv("GIN_MODE") == "debug" {
		// Use SQLite for local development
		return gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	// Use PostgreSQL for production
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, name, sslmode)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
