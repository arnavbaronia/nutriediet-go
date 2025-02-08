package client

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetExercisesForClient(c *gin.Context) {
	clientEmail, exists := c.Get("email")
	if !exists {
		fmt.Errorf("error: UpdateProfileByClient called but no email found: client_id: %s", c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"err": "email not found"})
		return
	}
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail.(string), c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	db := database.DB

	exercises := []model.Exercise{}
	err := db.Table("exercises").First(&exercises).Error
	if err != nil {
		fmt.Errorf("error: fetching all exercises: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isActive": true, "exercises": exercises})
	return
}
