package routes

import (
	"github.com/cd-Ishita/nutriediet-go/controller"
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
	incomingRoutes.GET("/admin/clients", adminController.GetAllClients)
}
