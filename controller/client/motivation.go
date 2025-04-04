package client

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetActiveMotivationsForClients(c *gin.Context) {
	clientEmail := c.GetString("email")
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	db := database.DB

	var notifs []model.Notification
	err := db.Model(&model.Notification{}).Where("posting_active = TRUE and deleted_at IS NULL").Find(&notifs).Error
	if err != nil {
		fmt.Errorf("error: GetActiveMotivationsForClients could not fetch notifications | err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isActive": true, "motivations": notifs})
	return
}
