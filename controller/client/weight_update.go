package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cd-Ishita/nutriediet-go/constants"
	"github.com/cd-Ishita/nutriediet-go/middleware"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateWeightForClient ...
// if there is an update from client side, it means
// no updations allowed after sending diet for 5 days
func UpdateWeightForClient(c *gin.Context) {
	isAllowed, isActive := middleware.ClientAuthentication(c.GetString("email"), c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

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
	err = db.Table("diet_histories").Where("client_id = ? and diet_type = ?", c.Param("client_id"), constants.RegularDiet.Uint32()).Order("date DESC").Select("id").First(&dietRecord).Error
	if err != nil {
		fmt.Println("Could not retrieve diet record for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req model.WeightUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = db.Table("diet_histories").Where("id = ? and diet_type = ?", dietRecord.ID, constants.RegularDiet.Uint32()).Update("weight", req.Weight).Update("feedback", req.Feedback).Error; err != nil {
		fmt.Println("Error while saving client diet record", dietRecord)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"isActive": true})
	return
}

func WeightUpdationStatus(c *gin.Context) {
	isAllowed, isActive := middleware.ClientAuthentication(c.GetString("email"), c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

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
		c.JSON(http.StatusOK, gin.H{"isActive": true, "status": "allowed"})
		return
	} else {
		fmt.Println("Weight updation not allowed for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusOK, gin.H{"status": "not_allowed"})
		return
	}
}

// IsWeightUpdationAllowed show the component to update weight only if this value comes true
func IsWeightUpdationAllowed(clientId string) (bool, error) {
	db := database.DB

	var date time.Time

	err := db.Table("diet_histories").Select("date").Where("client_id = ? and diet_type = ? and deleted_at IS NULL", clientId, constants.RegularDiet.Uint32()).Order("date DESC").Limit(1).Find(&date).Error
	if err != nil {
		return false, err
	}

	// Weight updation only allowed after 4 days of latest diet given
	// COMMENT OUT - for local testing
	allowedUpdationDate := date.Add(time.Hour * 24 * 4)
	allowedUpdationDate = time.Date(allowedUpdationDate.Year(), allowedUpdationDate.Month(), allowedUpdationDate.Day(), 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, date.Location())
	if currentDate.Before(allowedUpdationDate) {
		return false, nil
	}

	return true, nil
}

func GetWeightHistoryForClient(c *gin.Context) {
	db := database.DB

	clientEmail := c.GetString("email")
	isAllowed, _ := middleware.ClientAuthentication(clientEmail, c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
		return
	}

	var weightHistory []struct {
		Date   time.Time `json:"date"`
		Weight float32   `json:"weight"`
	}

	err := db.Model(&model.DietHistory{}).
		Where("client_id = ? and weight IS NOT NULL and deleted_at IS NULL", c.Param("client_id")).
		Order("date ASC").
		Select("date, weight").
		Find(&weightHistory).
		Error

	if err != nil {
		fmt.Errorf("error: GetWeightHistoryForClient | could not fetch weight history for client %s | err: %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weight_history": weightHistory})
}

// logic behind weight updation and diet submit
// 1. diet submit always creates new rows in the table
// 2. the weight gets added to the latest record of the table
