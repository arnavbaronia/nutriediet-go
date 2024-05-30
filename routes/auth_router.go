package routes

import (
	userController "github.com/cd-Ishita/nutriediet-go/controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", userController.SignUp)
	incomingRoutes.POST("users/login", userController.Login)
}
