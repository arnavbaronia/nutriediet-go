package routes

import (
	userController "github.com/cd-Ishita/nutriediet-go/controller"
	clientController "github.com/cd-Ishita/nutriediet-go/controller/client"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	// Rate limiter for authentication endpoints (5 requests/minute)
	authRateLimit := middleware.RateLimitAuth()
	
	// Rate limiter for password reset (3 requests/minute - stricter)
	strictRateLimit := middleware.RateLimitStrict()
	
	// Public routes with rate limiting to prevent abuse
	incomingRoutes.POST("/signup", authRateLimit, userController.SignUp)
	incomingRoutes.POST("/login", authRateLimit, userController.Login)
	incomingRoutes.POST("/create_user", authRateLimit, userController.CreateUser) // Public signup endpoint
	incomingRoutes.POST("/create_profile/:email", authRateLimit, clientController.CreateProfileByClient)
	
	// Password reset routes with strict rate limiting
	incomingRoutes.POST("/auth/forgot-password", strictRateLimit, userController.ForgotPassword)
	incomingRoutes.POST("/auth/reset-password", strictRateLimit, userController.ResetPassword)
}
