package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cd-Ishita/nutriediet-go/controller"
	database "github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// Don't fail if .env file doesn't exist (useful for production)
		log.Println("Warning: .env file not found, using system environment variables")
	} else {
		log.Println("✅ Environment variables loaded from .env file")
	}

	database.ConnectToDB()

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("no port found")
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

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
	router.Static("/images", "./images")

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.POST("/create_user", controller.CreateUser)
	router.GET("/get_users", controller.GetUsers)

	// ADMIN - EXERCISE
	router.GET("exercise", controller.GetExercisesForAdmin)
	router.GET("exercise/:exercise_id", controller.GetExercise)
	router.POST("exercise/:exercise_id/delete", controller.RemoveExerciseFromList)
	router.POST("exercise/:exercise_id/update", controller.UpdateExerciseFromList)
	router.POST("exercise/submit", controller.AddExerciseFromList)

	// ADMIN - DIET

	router.Run(":" + port) // listen and serve on 0.0.0.0:8081
}
