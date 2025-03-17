package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cd-Ishita/nutriediet-go/controller"
	"github.com/cd-Ishita/nutriediet-go/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	database "github.com/cd-Ishita/nutriediet-go/database"
)

func main() {
	database.ConnectToDB()

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("no port found")
		port = "8081"
	}

	router := gin.New()
	router.Use(gin.Logger())

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"https://nutriediet.netlify.app/"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Client-Email", "Request-Client-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

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
