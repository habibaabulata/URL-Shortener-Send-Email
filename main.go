package main

import (
	// Importing configuration, controllers & database package
	"net/http"
	"url-shortener/config"
	"url-shortener/controllers"
	"url-shortener/database"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig() // Load environment variables from .env file
	database.InitDB()   // Initialize the database connection

	// Create a new Gin router
	r := gin.Default()

	// Setup session store
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Register routes for authentication and URL shortening
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Group routes that require authentication
	auth := r.Group("/")
	auth.Use(AuthRequired)
	{
		// Protected routes
		auth.POST("/shorten", controllers.ShortenURL)        // Shorten a URL
		auth.GET("/:short_code", controllers.GetOriginalURL) // Retrieve the original URL
		auth.POST("/logout", controllers.Logout)             // Logout the user
	}

	r.Run(":8081")
}

// A middleware function to check if the user is authenticated
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
