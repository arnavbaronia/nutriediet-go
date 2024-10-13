package routes

import (
	"github.com/cd-Ishita/nutriediet-go/controller"
	clientController "github.com/cd-Ishita/nutriediet-go/controller/client"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/gin-gonic/gin"

	adminController "github.com/cd-Ishita/nutriediet-go/controller/admin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// to be used by admin and client routes both
	// only after authentication, these routes can be used
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.GET("/users", controller.GetUsers)
	incomingRoutes.GET("/user:user_id", controller.GetUser)

	// <<<<<<<<===============================================================================>>>>>>
	// CLIENT ROUTES

	// CLIENT - WEIGHT UPDATE
	incomingRoutes.POST(":client_id/weight_update", clientController.UpdateWeightForClient)
	incomingRoutes.GET(":client_id/weight_update", clientController.WeightUpdationStatus)

	// CLIENT - DIET
	incomingRoutes.GET(":client_id/diet", clientController.GetRegularDietForClient)
	incomingRoutes.GET(":client_id/detox_diet", clientController.GetDetoxDietForClient)

	// CLIENT - EXERCISE
	incomingRoutes.GET(":client_id/exercise", clientController.GetExercisesForClient)

	// CLIENT - PROFILE
	incomingRoutes.POST(":client_id/my_profile", clientController.UpdateProfileByClient)
	incomingRoutes.GET(":client_id/my_profile", clientController.GetProfileForClient)
	incomingRoutes.POST("/:email/create_profile", clientController.CreateProfileByClient)

	// <<<<<<<<===============================================================================>>>>>>
	// ADMIN ROUTES
	incomingRoutes.GET("/admin/clients", adminController.GetAllClients)
	incomingRoutes.GET("/admin/client/:client_id", adminController.GetClientInfo)
	incomingRoutes.POST("/admin/client/:client_id", adminController.UpdateClientInfo)
	incomingRoutes.POST("/admin/client/:client_id/activation", adminController.ActivateOrDeactivateClientAccount)

}
