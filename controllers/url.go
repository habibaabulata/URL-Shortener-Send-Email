package controllers

import (
	"math/rand"
	"net/http"
	"time"
	"url-shortener/database"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
)

// Constants for generating short codes
const shortCodeLength = 8
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generates a random short code for the URL
func generateShortCode() string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	shortCode := make([]byte, shortCodeLength)
	for i := range shortCode {
		shortCode[i] = charset[random.Intn(len(charset))]
	}
	return string(shortCode)
}

// handles shortening a URL
func ShortenURL(c *gin.Context) {
	var request struct {
		OriginalURL string `json:"original_url"` // JSON payload containing the original URL
	}

	// Bind the JSON payload to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a short code for the URL
	shortCode := generateShortCode()
	url := models.URL{
		ShortCode:   shortCode,
		OriginalURL: request.OriginalURL,
		UserID:      1, // Replace with actual user ID
	}

	// Save the URL to the database
	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_code": shortCode})
}

// handles retrieving the original URL from the short code
func GetOriginalURL(c *gin.Context) {
	shortCode := c.Param("short_code") // Retrieve the short code from the URL parameter
	var url models.URL

	// Find the URL in the database using the short code
	if err := database.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Redirect to the original URL
	c.Redirect(http.StatusFound, url.OriginalURL)
}
