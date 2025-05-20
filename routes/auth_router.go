package routes

import (
	"github.com/cd-Ishita/nutriediet-go/api"
	userController "github.com/cd-Ishita/nutriediet-go/controller"
	clientController "github.com/cd-Ishita/nutriediet-go/controller/client"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// these are public routes, hence in AuthRoutes
	// open to all
	incomingRoutes.POST("signup", userController.SignUp)
	incomingRoutes.POST("login", userController.Login)
	incomingRoutes.POST("/create_profile/:email", clientController.CreateProfileByClient)
	incomingRoutes.POST("/password-reset/initiate", api.InitiatePasswordReset)
	incomingRoutes.POST("/password-reset/complete", api.CompletePasswordReset)
}
