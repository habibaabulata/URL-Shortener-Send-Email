package controllers

import (
	"net/http"
	"url-shortener/config"
	"url-shortener/database"
	"url-shortener/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // Importing bcrypt for password hashing
)

// Register handles user registration
func Register(c *gin.Context) {
	var user models.User

	// Bind the JSON payload to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save the user to the database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Send a welcome email
	subject := newFunction()
	body := "<h1>Welcome to URL Shortener Service</h1><p>Thank you for registering!</p>"
	if err := config.SendEmail(user.Email, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func newFunction() string {
	subject := "Welcome to URL Shortener Service"
	return subject
}

// Login handles user login
func Login(c *gin.Context) {
	var user models.User

	// Bind the JSON payload to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User

	// Retrieve the user from the database using the provided email
	if err := database.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user", dbUser.Email)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
