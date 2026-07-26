package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetListOfRecipes(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	db := database.DB
	var recipeRows []model.Recipe
	if err := db.Where("deleted_at IS NULL").Find(&recipeRows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch recipes"})
		return
	}

	res := make([]model.GetListOfRecipesResponse, len(recipeRows))
	for i, recipe := range recipeRows {
		res[i] = model.GetListOfRecipesResponse{
			Name:     model.RecipeDisplayName(recipe),
			RecipeID: recipe.ID,
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "list": res})
}

func CreateRecipe(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	var req model.RecipeCard
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	row := model.RecipeFromPayload(req)
	db := database.DB
	if err := db.Create(&row).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"recipe":  model.ToRecipeCard(row),
	})
}

func UpdateRecipeByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	recipeID, err := strconv.ParseUint(c.Param("recipe_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	var req model.RecipeCard
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	db := database.DB
	var existing model.Recipe
	if err := db.Where("id = ? AND deleted_at IS NULL", recipeID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe"})
		return
	}

	updated := model.RecipeFromPayload(req)
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt
	if updated.Slug == "" {
		updated.Slug = existing.Slug
	}
	if updated.ImageURL == "" {
		updated.ImageURL = existing.ImageURL
	}

	if err := db.Save(&updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"recipe":  model.ToRecipeCard(updated),
	})
}

func GetRecipeByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	recipeID, err := strconv.ParseUint(c.Param("recipe_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	db := database.DB
	var row model.Recipe
	if err := db.Where("id = ? AND deleted_at IS NULL", recipeID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"recipe":  model.ToRecipeCard(row),
	})
}

func DeleteRecipeByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access by client"})
		return
	}

	recipeID, err := strconv.ParseUint(c.Param("recipe_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	db := database.DB
	now := time.Now()
	result := db.Model(&model.Recipe{}).
		Where("id = ? AND deleted_at IS NULL", recipeID).
		Update("deleted_at", now)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete recipe"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "recipe deleted successfully",
	})
}
