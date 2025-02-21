package routes

import (
	"github.com/cd-Ishita/nutriediet-go/controller"
	clientController "github.com/cd-Ishita/nutriediet-go/controller/client"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/gin-gonic/gin"

	adminController "github.com/cd-Ishita/nutriediet-go/controller/admin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// Authentication middleware applies to all routes
	incomingRoutes.Use(middleware.Authenticate)

	// USER ROUTES
	incomingRoutes.GET("/users", controller.GetUsers)
	incomingRoutes.GET("/user/:user_id", controller.GetUser)

	// <<<<<<<<===============================================================================>>>>>>
	// CLIENT ROUTES (Prefix with `/clients` for all client-related routes)

	// CLIENT - WEIGHT UPDATE
	incomingRoutes.POST("/clients/:client_id/weight_update", clientController.UpdateWeightForClient)
	incomingRoutes.GET("/clients/:client_id/weight_update", clientController.WeightUpdationStatus)

	// CLIENT - DIET
	incomingRoutes.GET("/clients/:client_id/diet", clientController.GetRegularDietForClient)
	incomingRoutes.GET("/clients/:client_id/detox_diet", clientController.GetDetoxDietForClient)

	// CLIENT - EXERCISE
	incomingRoutes.GET("/clients/:client_id/exercise", clientController.GetExercisesForClient)

	// CLIENT - PROFILE
	incomingRoutes.POST("/clients/:client_id/my_profile", clientController.UpdateProfileByClient)
	incomingRoutes.GET("/clients/:client_id/my_profile", clientController.GetProfileForClient)
	incomingRoutes.GET("/clients/:client_id/profile_created", clientController.HasClientCreatedProfile)
	// EMAIL-BASED PROFILE CREATION (Separate from client routes to avoid conflicts)
	incomingRoutes.POST("/users/:email/create_profile", clientController.CreateProfileByClient)

	// CLIENT - RECIPE
	incomingRoutes.GET(":meal_id/recipe", clientController.GetRecipeByMealIDForClient)

	// <<<<<<<<===============================================================================>>>>>>
	// ADMIN ROUTES (Prefix with `/admin` for all admin-related routes)

	incomingRoutes.GET("/admin/clients", adminController.GetAllClients)
	incomingRoutes.GET("/admin/client/:client_id", adminController.GetClientInfo)
	incomingRoutes.POST("/admin/client/:client_id", adminController.UpdateClientInfo)
	incomingRoutes.POST("/admin/client/:client_id/activation", adminController.ActivateOrDeactivateClientAccount)
	incomingRoutes.GET("/admin/client/:client_id/weight_history", adminController.GetWeightHistoryForClient)
	incomingRoutes.GET("/admin/client/:client_id/diet_history", adminController.GetDietHistoryForClient)

	// ADMIN - DIET
	incomingRoutes.GET("/admin/meal_list", adminController.GetMealList)
	incomingRoutes.GET("/admin/quantity_list", adminController.GetQuantityList)
	incomingRoutes.POST("/admin/:client_id/diet", adminController.SaveDietForClient)

	// <<<<<<<<===============================================================================>>>>>>

	// ADMIN - DIET TEMPLATES
	incomingRoutes.GET("/admin/diet_templates", adminController.GetDietTemplatesList)
	incomingRoutes.GET("/admin/diet_templates/:diet_template_id", adminController.GetDietTemplateByID)
	incomingRoutes.POST("/admin/diet_templates/new", adminController.CreateDietTemplate)
	incomingRoutes.POST("/admin/diet_templates/:diet_template_id", adminController.UpdateDietTemplate)
	incomingRoutes.POST("/admin/diet_templates/:diet_template_id/delete", adminController.DeleteDietTemplateByID)

	// ADMIN - RECIPES
	incomingRoutes.GET("/admin/recipe/:id", adminController.GetRecipeByID)
	incomingRoutes.POST("/admin/recipe/:id", adminController.UpdateRecipeByID)
	incomingRoutes.POST("/admin/recipe/new", adminController.CreateRecipe)
	incomingRoutes.POST("/admin/recipe/:id/delete", adminController.DeleteRecipeByID)
	incomingRoutes.GET("/admin/recipes", adminController.GetListOfRecipes)

	// ADMIN - EXERCISES
	incomingRoutes.GET("/admin/exercises", adminController.GetListOfExercises)
	incomingRoutes.GET("/admin/exercise/:exercise_id", adminController.GetExerciseByID)
	incomingRoutes.POST("/admin/exercise/new", adminController.CreateExercise)
	incomingRoutes.POST("/admin/exercise/:exercise_id", adminController.UpdateExerciseByID)
	incomingRoutes.POST("/admin/exercise/:exercise_id/delete", adminController.DeleteExerciseByID)
}
