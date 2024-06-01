package main

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/controller"
	clientController "github.com/cd-Ishita/nutriediet-go/controller/client"
	"github.com/gin-gonic/gin"
	"os"

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

	//routes.AuthRoutes(router)
	//routes.ClientRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.POST("/create_user", controller.CreateUser)
	router.GET("/get_users", controller.GetUsers)

	// CLIENT - PROFILE
	router.POST(":client_id/my_profile", clientController.SaveProfileByClient)
	router.GET(":client_id/my_profile", clientController.GetProfileForClient)

	// CLIENT - WEIGHT_UPDATE
	router.POST(":client_id/weight_update", clientController.UpdateWeightForClient)
	router.GET(":client_id/weight_update", clientController.WeightUpdationStatus)

	// CLIENT - DIET
	router.GET(":client_id/diet", clientController.GetRegularDietForClient)
	router.GET(":client_id/detox_diet", clientController.GetDetoxDietForClient)

	// CLIENT _ EXERCISE
	router.GET(":client_id/exercise", controller.GetExercisesForClient)

	// ADMIN - EXERCISE
	router.GET("exercise", controller.GetExercisesForAdmin)
	router.GET("exercise/:exercise_id", controller.GetExercise)
	router.POST("exercise/:exercise_id/delete", controller.RemoveExerciseFromList)
	router.POST("exercise/:exercise_id/update", controller.UpdateExerciseFromList)
	router.POST("exercise/submit", controller.AddExerciseFromList)

	// ADMIN - DIET
	router.POST(":client_id/diet", controller.SaveDietForClient)

	router.Run(":" + port) // listen and serve on 0.0.0.0:8080
}
