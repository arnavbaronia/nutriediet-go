package admin

import (
	"errors"
	"fmt"
	"net/http"

	// "strconv"
	// "time"

	"github.com/cd-Ishita/nutriediet-go/constants"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetWeightHistoryForClient(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	clientID := c.Param("client_id")
	if clientID == "" || clientID == "0" {
		fmt.Errorf("error: client_id cannot be empty string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id cannot be empty string"})
		return
	}

	db := database.DB
	var res []model.GetWeightHistoryForClientResponse
	err := db.Model(model.DietHistory{}).
		Where("client_id = ? and diet_type = ?", clientID, constants.RegularDiet.Uint32()).
		Select("weight", "date").
		Find(&res).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: could not find diet_history_id %d for client_id %s", clientID, c.Param("client_id"))
		c.JSON(http.StatusOK, gin.H{"response": nil})
		return
	} else if err != nil {
		fmt.Errorf("error: cannot fetch weights for client with id: %d | err: %v", clientID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res})
	return
}

func UpdateWeightForClientByAdmin(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	clientID := c.Param("client_id")
	if clientID == "" || clientID == "0" {
		fmt.Errorf("error: client_id cannot be empty string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id cannot be empty string"})
		return
	}

	db := database.DB

	var req model.WeightUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dietRecord := model.DietHistory{}
	err := db.Table("diet_histories").Where("client_id = ? and diet_type = ? and week_number = ?", c.Param("client_id"), constants.RegularDiet.Uint32(), req.WeekNumber).Order("date DESC").Select("id").First(&dietRecord).Error
	if err != nil {
		fmt.Println("Could not retrieve diet record for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Table("diet_histories").Where("id = ? and diet_type = ?", dietRecord.ID, constants.RegularDiet.Uint32()).Updates(model.DietHistory{
		Weight:   &req.Weight,
		Feedback: req.Feedback,
	}).Error; err != nil {
		fmt.Println("Error while saving client diet record", dietRecord)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Weight updated successfully",
		"success": true,
	})
}
