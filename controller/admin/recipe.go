package admin

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//func GetRecipeByID(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	db := database.DB
//
//	recipe := model.Recipe{}
//	if err := db.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			fmt.Errorf("error: GetRecipeByID | recipe does not exist with id: %d", c.Param("id"))
//			c.JSON(http.StatusNotFound, gin.H{"error": err})
//			return
//		}
//		fmt.Errorf("error: GetRecipeByID could not fetch recipe with id %d | err: %v", c.Param("id"), err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	ingredientsList := strings.Split(recipe.Ingredients, ";")
//	prepList := strings.Split(recipe.Preparation, ";")
//
//	res := model.GetRecipeResponse{
//		ID:          recipe.ID,
//		Name:        recipe.Name,
//		Ingredients: ingredientsList,
//		Preparation: prepList,
//	}
//
//	c.JSON(http.StatusOK, gin.H{"recipe": res})
//	return
//}
//
//func UpdateRecipeByID(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	var recipeReq model.UpdateRecipeRequest
//	if err := c.BindJSON(&recipeReq); err != nil {
//		fmt.Errorf("error: UpdateRecipeByID | could not extract request from context | err : %v", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	ingredients := strings.Join(recipeReq.Ingredients, ";")
//	steps := strings.Join(recipeReq.Preparation, ";")
//
//	recipe := model.Recipe{
//		ID:          recipeReq.ID,
//		Name:        recipeReq.Name,
//		Ingredients: ingredients,
//		Preparation: steps,
//	}
//	db := database.DB
//	if err := db.Model(&model.Recipe{}).Where("id = ?", c.Param("id")).Select("name", "ingredients", "preparation").Updates(&recipe).Error; err != nil {
//		fmt.Errorf("error: UpdateRecipeByID | could not save recipe %v | err: %v", recipe, err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true})
//	return
//}
//
//func CreateRecipe(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	var recipeReq model.CreateRecipeRequest
//	if err := c.BindJSON(&recipeReq); err != nil {
//		fmt.Errorf("error: CreateRecipe | could not extract request from context | err : %v", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	ingredients := ""
//	for _, ingredient := range recipeReq.Ingredients {
//		if ingredients == "" {
//			ingredients = ingredient
//		} else {
//			ingredients = ingredients + ";" + ingredient
//		}
//	}
//
//	steps := ""
//	for _, prep := range recipeReq.Preparation {
//		if steps == "" {
//			steps = prep
//		} else {
//			steps = steps + ";" + prep
//		}
//	}
//
//	recipe := model.Recipe{
//		Name:        recipeReq.Name,
//		Ingredients: ingredients,
//		Preparation: steps,
//	}
//	db := database.DB
//	if err := db.Save(&recipe).Error; err != nil {
//		fmt.Errorf("error: CreateRecipe | could not save recipe %v | err: %v", recipe, err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true})
//	return
//}
//
//func DeleteRecipeByID(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	recipeID, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meal ID"})
//		return
//	}
//
//	db := database.DB
//	if err := db.Model(&model.Recipe{}).Where("id = ?", recipeID).Update("deleted_at", time.Now()).Error; err != nil {
//		fmt.Errorf("error: DeleteRecipeByMealID | could not delete recipe with id: %v | err: %v", c.Param("id"), err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true})
//	return
//}

func GetListOfRecipes(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	db := database.DB
	var recipes []model.Recipe
	if err := db.Where("deleted_at IS NULL").Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch recipes"})
		return
	}

	res := make([]model.GetListOfRecipesResponse, len(recipes))
	for i, recipe := range recipes {
		res[i] = model.GetListOfRecipesResponse{
			Name:     recipe.Name,
			RecipeID: recipe.ID,
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "list": res})
}

func UploadRecipeImage(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file received"})
		return
	}

	imageName := strings.TrimSpace(c.PostForm("name"))
	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or empty image name"})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Read the file content
	imageData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // default if not specified
	}

	// Create recipe record
	recipe := model.Recipe{
		Name:      imageName,
		ImageData: imageData,
		ImageType: contentType,
	}

	db := database.DB
	if err := db.Create(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recipe to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "upload successful",
		"recipe":  recipe,
	})
}

func GetRecipeImageForAdmin(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	recipeID := c.Param("recipe_id")
	if recipeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing recipe id"})
		return
	}

	db := database.DB
	var recipe model.Recipe

	err := db.Where("id = ? AND deleted_at IS NULL", recipeID).First(&recipe).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe"})
		return
	}

	// Check Accept header to determine response type
	acceptHeader := c.GetHeader("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		// Return JSON response with recipe details
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"recipe": gin.H{
				"ID":   recipe.ID,
				"Name": recipe.Name,
			},
		})
	} else {
		// Return image data
		c.Data(http.StatusOK, recipe.ImageType, recipe.ImageData)
	}
}

func UpdateRecipeImageByAdmin(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	recipeID := c.Param("recipe_id")
	if recipeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing recipe id"})
		return
	}

	db := database.DB

	// First find the existing recipe
	var existingRecipe model.Recipe
	if err := db.Where("id = ? AND deleted_at IS NULL", recipeID).First(&existingRecipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe"})
		return
	}

	// Process the file if provided
	file, _ := c.FormFile("file")
	imageName := strings.TrimSpace(c.PostForm("name"))

	if file != nil {
		// Open the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer src.Close()

		// Read the file content
		imageData, err := io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		// Update image data and type
		existingRecipe.ImageData = imageData
		existingRecipe.ImageType = file.Header.Get("Content-Type")
		if existingRecipe.ImageType == "" {
			existingRecipe.ImageType = "application/octet-stream"
		}
	}

	if imageName != "" {
		existingRecipe.Name = imageName
	}

	// Update the recipe
	if err := db.Save(&existingRecipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update successful",
		"recipe":  existingRecipe,
	})
}

func DeleteRecipeImageByAdmin(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	recipeID := c.Param("recipe_id")
	if recipeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing recipe id"})
		return
	}

	db := database.DB

	// First find the existing recipe
	var recipe model.Recipe
	if err := db.Where("id = ? AND deleted_at IS NULL", recipeID).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe"})
		return
	}

	// Soft delete by setting DeletedAt
	now := time.Now()
	recipe.DeletedAt = &now

	if err := db.Save(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "recipe deleted successfully",
	})
}
