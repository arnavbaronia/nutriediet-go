package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	database "github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// init runs before main() and before any package init() functions
func init() {
	// Load environment variables FIRST, before any other package initialization
	err := godotenv.Load()
	if err != nil {
		// Don't fail if .env file doesn't exist (useful for production)
		log.Println("Warning: .env file not found, using system environment variables")
	} else {
		log.Println("✅ Environment variables loaded from .env file")
	}
}

func main() {

	database.ConnectToDB()

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("no port found")
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Security headers - protects against XSS, clickjacking, etc.
	router.Use(middleware.SecurityHeaders())

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"https://nutriediet.netlify.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Client-Email", "Request-Client-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	
	// Serve static images with no-cache headers (for development)
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/images/") {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	})
	router.Static("/images", "./images")

	// Register routes
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Start server
	log.Printf("🚀 Server starting on port %s", port)
	router.Run(":" + port)
}
