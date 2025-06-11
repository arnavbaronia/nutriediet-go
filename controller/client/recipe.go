package client

import (
	"errors"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	// Fetch only non-deleted recipes with their image data
	err := db.Select("id, name, image_data, image_type").Where("deleted_at IS NULL").Find(&recipes).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"recipes": []interface{}{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch recipes",
			"details": err.Error(),
		})
		return
	}

	// Prepare response with base64 encoded images
	response := make([]gin.H, len(recipes))
	for i, recipe := range recipes {
		imageBase64 := ""
		if len(recipe.ImageData) > 0 {
			imageBase64 = "data:" + recipe.ImageType + ";base64," +
				helpers.BytesToBase64(recipe.ImageData)
		}

		response[i] = gin.H{
			"id":        recipe.ID,
			"name":      recipe.Name,
			"imageData": imageBase64,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isActive": true,
		"recipes":  response,
	})
}

func GetRecipeImageForClient(c *gin.Context) {
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

	db := database.DB
	var recipe model.Recipe

	err := db.Select("id, name, image_data, image_type").
		Where("id = ? AND deleted_at IS NULL", recipeID).
		First(&recipe).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch recipe",
			"details": err.Error(),
		})
		return
	}

	// Return image data directly
	c.Data(http.StatusOK, recipe.ImageType, recipe.ImageData)
}
