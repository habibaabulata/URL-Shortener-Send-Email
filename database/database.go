package database

import (
	// Importing MySQL driver for Gorm, Gorm package, configuration package & models package
	"log"
	"url-shortener/config"
	"url-shortener/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB //database instance

// initializes the database connection
func InitDB() {
	dsn := config.GetDSN() // Get the DSN from environment variables

	// Open a connection to the database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// Auto-migrate the User and URL models
	DB.AutoMigrate(&models.User{}, &models.URL{})
}
