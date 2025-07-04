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
	incomingRoutes.GET("/clients/:client_id/weight-history", clientController.GetWeightHistoryForClient)

	// CLIENT - DIET
	incomingRoutes.GET("/clients/:client_id/diet", clientController.GetDietsForClient)

	// CLIENT - EXERCISE
	incomingRoutes.GET("/clients/:client_id/exercise", clientController.GetExercisesForClient)
	incomingRoutes.POST("/clients/:client_id/exercise/favorite", clientController.ToggleFavoriteExercise)

	// CLIENT - PROFILE
	incomingRoutes.POST("/clients/:client_id/my_profile", clientController.UpdateProfileByClient)
	incomingRoutes.GET("/clients/:client_id/my_profile", clientController.GetProfileForClient)
	incomingRoutes.GET("/clients/:client_id/profile_created", clientController.HasClientCreatedProfile)

	// CLIENT - RECIPE
	incomingRoutes.GET("clients/:client_id/recipe", clientController.GetRecipeImageForClients)

	// CLIENT - MOTIVATION
	incomingRoutes.GET("clients/:client_id/motivation", clientController.GetActiveMotivationsForClients)

	// EMAIL-BASED PROFILE CREATION (Separate from client routes to avoid conflicts)
	incomingRoutes.POST("/users/:email/create_profile", clientController.CreateProfileByClient)

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
	incomingRoutes.POST("/admin/:client_id/edit_diet", adminController.EditDietForClient)
	incomingRoutes.POST("/admin/:client_id/weight_update", adminController.UpdateWeightForClientByAdmin)
	incomingRoutes.POST("/admin/:client_id/delete_diet", adminController.DeleteDietForClientByAdmin)
	incomingRoutes.POST("/admin/common_diet", adminController.SaveCommonDietForClients)
	incomingRoutes.GET("/admin/common_diet/:group_id", adminController.GetCommonDietsHistory)
	incomingRoutes.POST("/admin/common_diet/:group_id/update", adminController.EditCommonDiet)
	incomingRoutes.POST("/admin/common_diet/:group_id/delete_diet", adminController.DeleteCommonDiet)

	// <<<<<<<<===============================================================================>>>>>>

	// ADMIN - DIET TEMPLATES
	incomingRoutes.GET("/admin/diet_templates", adminController.GetDietTemplatesList)
	incomingRoutes.GET("/admin/diet_templates/:diet_template_id", adminController.GetDietTemplateByID)
	incomingRoutes.POST("/admin/diet_templates/new", adminController.CreateDietTemplate)
	incomingRoutes.POST("/admin/diet_templates/:diet_template_id", adminController.UpdateDietTemplate)
	incomingRoutes.POST("/admin/diet_templates/:diet_template_id/delete", adminController.DeleteDietTemplateByID)

	// ADMIN - RECIPES
	//incomingRoutes.GET("/admin/recipe/:id", adminController.GetRecipeByID)
	//incomingRoutes.POST("/admin/recipe/:id", adminController.UpdateRecipeByID)
	//incomingRoutes.POST("/admin/recipe/new", adminController.CreateRecipe)
	//incomingRoutes.POST("/admin/recipe/:id/delete", adminController.DeleteRecipeByID)
	incomingRoutes.GET("/admin/recipes", adminController.GetListOfRecipes)
	incomingRoutes.POST("/admin/recipes/upload", adminController.UploadRecipeImage)
	incomingRoutes.GET("/admin/recipes/:recipe_id", adminController.GetRecipeImageForAdmin)
	incomingRoutes.POST("/admin/recipes/:recipe_id/update", adminController.UpdateRecipeImageByAdmin)
	incomingRoutes.POST("/admin/recipes/:recipe_id/delete", adminController.DeleteRecipeImageByAdmin)

	// ADMIN - EXERCISES
	incomingRoutes.GET("/admin/exercises", adminController.GetListOfExercises)
	incomingRoutes.GET("/admin/exercise/:exercise_id", adminController.GetExerciseByID)
	incomingRoutes.POST("/admin/exercise/new", adminController.CreateExercise)
	incomingRoutes.POST("/admin/exercise/:exercise_id", adminController.UpdateExerciseByID)
	incomingRoutes.POST("/admin/exercise/:exercise_id/delete", adminController.DeleteExerciseByID)

	// ADMIN - MOTIVATION
	incomingRoutes.POST("/admin/motivations/new", adminController.CreateNewMotivation)
	incomingRoutes.POST("/admin/motivation/:motivation_id/unpost", adminController.UnpostMotivation)
	incomingRoutes.POST("/admin/motivation/:motivation_id/post", adminController.PostMotivation)
	incomingRoutes.GET("/admin/motivation", adminController.GetAllMotivations)
}
