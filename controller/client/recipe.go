package client

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetRecipeByMealIDForClient(c *gin.Context) {
	// For Client users, need to check if account is active
	clientEmail := c.GetString("email")
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}
	if !isActive {
		fmt.Errorf("error: GetRecipeByMealIDForClient | client inactive | clientEmail: %s", c.Param("email"))
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	db := database.DB

	recipe := model.Recipe{}
	if err := db.Where("meal_id = ?", c.Param("meal_id")).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Errorf("error: GetRecipeByMealIDForClient | recipe does not exist with meal_id: %d", c.Param("meal_id"))
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		fmt.Errorf("error: GetRecipeByMealIDForClient could not fetch recipe with meal_id %d | err: %v", c.Param("meal_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe, "isActive": isActive})
	return
}
