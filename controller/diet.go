package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// SaveDietForClient questions
// should the week number be input by the UI
// should the week number be calculated directly
// should the date be the date of diet updation, or the date calculated using the group etc
// edit button separately
func SaveDietForClient(c *gin.Context) {
	// Parse the request body to extract the diet information
	var schedule model.SaveDietForClientRequest
	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB
	// this is necessary because the client has already updated the weight and their feedback, the diet has to be uploaded in that record only
	// however, what if client has not updated weight yet?
	dietHistoryRecord := model.DietHistory{}
	err := db.Where("client_id = ?", c.Param("client_id")).Order("date DESC").First(&dietHistoryRecord).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dietJSON, err := json.Marshal(schedule.Diet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal diet to JSON"})
		return
	}

	//clientID, _ := strconv.ParseUint(c.Param("client_id"), 10, 64)
	//dietHistory := model.DietHistory{
	//	ClientID:   clientID,
	//	WeekNumber: schedule.WeekNumber,
	//	Date:       time.Now(),
	//}
	//
	//// Save the diet history record to the database
	//if err := db.Save(&dietHistory).Error; err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	fmt.Println("dietJSON", dietJSON)

	if err = db.Table("diet_histories").Where("id = ?", 6).Update("diet", dietJSON).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Diet information saved successfully"})
	return
}
