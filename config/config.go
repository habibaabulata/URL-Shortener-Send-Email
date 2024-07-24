package config

import (
	"log"
	"os"
	"gopkg.in/gomail.v2"
	"github.com/joho/godotenv" // Importing the godotenv package to load environment variables
)

// LoadConfig loads the environment variables from the .env file
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
}

// GetDSN constructs the Data Source Name (DSN) for connecting to the MySQL database
func GetDSN() string {
	return os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// EmailConfig holds the email configuration details
type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}

// SendEmail sends an email using the provided configuration
func SendEmail(to, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", os.Getenv("EMAIL_USERNAME"))
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 587, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))

    return d.DialAndSend(m)
}
