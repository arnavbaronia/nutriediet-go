package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetDietByDietHistoryID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var dietHistoryID uint64
	if err := c.BindJSON(&dietHistoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	var dietJSON []string
	err := db.Model(&model.DietHistory{}).Where("client_id = ? and id = ?", c.Param("client_id"), dietHistoryID).Pluck("diet", &dietJSON).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: diet does not exist with diet history id %d", dietHistoryID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Errorf("error: could not fetch diet with diet_history_id %d for client_id %s", dietHistoryID, c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dietFinal := model.DietSchedule{}
	err = json.Unmarshal([]byte(dietJSON[0]), &dietFinal)
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"diet": dietFinal})
	return
}
