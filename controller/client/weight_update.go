package client

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// UpdateWeightForClient ...
// if there is an update from client side, it means
// no updations allowed after sending diet for 5 days
func UpdateWeightForClient(c *gin.Context) {
	db := database.DB

	status, err := IsWeightUpdationAllowed(c.Param("client_id"))
	if errors.Is(gorm.ErrRecordNotFound, err) {
		fmt.Println("Record not found for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Println("Error fetching weight updation allowed for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !status {
		fmt.Println("Weight updation not allowed for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusNotAcceptable, gin.H{"status": "not_allowed"})
		return
	}

	dietRecord := model.DietHistory{}
	err = db.Where("client_id = ?", c.Param("client_id")).Order("date DESC").First(&dietRecord).Error
	if err != nil {
		fmt.Println("Could not retrieve diet record for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req := float32(0)
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dietRecord.Weight = req
	if err = db.Save(&dietRecord).Error; err != nil {
		fmt.Println("Error while saving client diet record", dietRecord)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	return
}

func WeightUpdationStatus(c *gin.Context) {
	status, err := IsWeightUpdationAllowed(c.Param("client_id"))
	if errors.Is(gorm.ErrRecordNotFound, err) {
		fmt.Println("Record not found for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Println("Error fetching weight updation allowed for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if status {
		fmt.Println("Weight updation allowed for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusOK, gin.H{"status": "allowed"})
		return
	} else {
		fmt.Println("Weight updation not allowed for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusNotAcceptable, gin.H{"status": "not_allowed"})
		return
	}
}

// IsWeightUpdationAllowed show the component to update weight only if this value comes true
func IsWeightUpdationAllowed(clientId string) (bool, error) {

	db := database.DB

	date := time.Time{}
	err := db.Where("client_id = ?", clientId).Order("date DESC").Select("date").First(&date).Error
	if err != nil {
		return false, err
	}

	allowedUpdationDate := date.Add(time.Hour * 24 * 4)
	currentDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, date.Location())
	if currentDate.Before(allowedUpdationDate) {
		return false, nil
	}

	return true, nil
}

// logic behind weight updation and diet submit
// 1. diet submit always creates new rows in the table
// 2. the weight gets added to the latest record of the table
