package client

import (
	"errors"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//func GetRecipesForClient(c *gin.Context) {
//	// For Client users, need to check if account is active
//	clientEmail := c.GetString("email")
//	fmt.Println("GetRecipesForClient", clientEmail)
//	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, c.Param("client_id"))
//	if !isAllowed {
//		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": clientEmail, "requestClientID": c.Param("client_id")})
//		return
//	}
//	if !isActive {
//		fmt.Errorf("error: GetRecipeByMealIDForClient | client inactive | clientEmail: %s", c.Param("email"))
//		c.JSON(http.StatusOK, gin.H{"isActive": false})
//		return
//	}
//
//	db := database.DB
//
//	recipes := []model.Recipe{}
//	if err := db.Model(&model.Recipe{}).Where("deleted_at IS NULL").Find(&recipes).Error; err != nil {
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
//	res := []model.GetRecipeResponse{}
//	for _, recipe := range recipes {
//		ingredientsList := strings.Split(recipe.Ingredients, ";")
//		prepList := strings.Split(recipe.Preparation, ";")
//
//		res = append(res, model.GetRecipeResponse{
//			ID:          recipe.ID,
//			Name:        recipe.Name,
//			Ingredients: ingredientsList,
//			Preparation: prepList,
//		})
//	}
//
//	c.JSON(http.StatusOK, gin.H{"recipe": res, "isActive": isActive})
//	return
//}

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

	// Fetch only non-deleted recipes with their image URLs
	err := db.Select("id, name, image_url").Where("deleted_at IS NULL").Find(&recipes).Error
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

	// Return simplified response with just the essential data
	response := make([]gin.H, len(recipes))
	for i, recipe := range recipes {
		response[i] = gin.H{
			"id":       recipe.ID,
			"name":     recipe.Name,
			"imageUrl": recipe.ImageURL,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isActive": true,
		"recipes":  response,
	})
}
