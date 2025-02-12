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
	"strconv"
	"time"
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

// SaveDietForClient questions
// should the week number be input by the UI
// should the week number be calculated directly
// should the date be the date of diet updation, or the date calculated using the group etc
// edit button separately
func SaveDietForClient(c *gin.Context) {

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	// Parse the request body to extract the diet information
	var schedule model.SaveDietForClientRequest
	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	//dietHistoryRecord := model.DietHistory{}
	//err := db.Where("client_id = ?", c.Param("client_id")).Order("date DESC").First(&dietHistoryRecord).Error
	//if errors.Is(gorm.ErrRecordNotFound, err) {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//	return
	//} else if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// a new diet always creates a new record in the diet history table
	dietJSON, err := json.Marshal(schedule.Diet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal diet to JSON"})
		return
	}

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

	clientID, _ := strconv.ParseUint(c.Param("client_id"), 10, 64)
	emptyDietRecord := model.DietHistory{
		WeekNumber: schedule.WeekNumber,
		ClientID:   clientID,
		Date:       time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
		Weight:     nil,
		DietType:   schedule.DietType,
	}
	if err := db.Create(&emptyDietRecord).Error; err != nil {
		fmt.Errorf("error: SaveDietForClient | could not create empty diet_history_id %d for client_id %s | err: %v", schedule.Diet, clientID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = db.Table("diet_histories").Where("id = ?", emptyDietRecord.ID).Update("diet", dietJSON).Error; err != nil {
		fmt.Errorf("error: SaveDietForClient | could not save diet for diet_history_id %d for client_id %s | err: %v", schedule.Diet, clientID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Diet information saved successfully"})
	return
}
