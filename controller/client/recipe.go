package client

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRecipeImageForClients fetches all recipes for clients
func GetRecipeImageForClients(c *gin.Context) {
	clientEmail := c.GetString("email")
	clientID := c.Param("client_id")

	// Authentication check
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, clientID)
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized access",
			"details": gin.H{
				"clientEmail":     clientEmail,
				"requestClientID": clientID,
			},
		})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	db := database.DB
	var recipes []model.Recipe

	// Fetch only non-deleted recipes
	err := db.Where("deleted_at IS NULL").Find(&recipes).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"isActive": true,
				"recipes":  []interface{}{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch recipes",
			"details": err.Error(),
		})
		return
	}

	// Return response with recipe data
	response := make([]gin.H, len(recipes))
	for i, recipe := range recipes {
		response[i] = gin.H{
			"id":   recipe.ID,
			"name": recipe.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isActive": true,
		"recipes":  response,
	})
}

func GetSingleRecipeImageForClient(c *gin.Context) {
	clientEmail := c.GetString("email")
	clientID := c.Param("client_id")
	recipeID := c.Param("recipe_id")

	// Authentication check
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, clientID)
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized access",
			"details": gin.H{
				"clientEmail":     clientEmail,
				"requestClientID": clientID,
			},
		})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

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
		if recipe.ImageData != nil && len(recipe.ImageData) > 0 {
			c.Data(http.StatusOK, recipe.ImageType, recipe.ImageData)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "no image data available"})
		}
	}
}
