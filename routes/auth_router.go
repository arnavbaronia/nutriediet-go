package routes

import (
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
	
	// Password reset routes
	incomingRoutes.POST("/auth/forgot-password", userController.ForgotPassword)
	incomingRoutes.POST("/auth/reset-password", userController.ResetPassword)
}
