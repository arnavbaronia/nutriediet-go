package admin

import (
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

func GetDietHistoryForClient(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	var dietHistory []model.DietHistory
	err := db.Model(&model.DietHistory{}).Where("client_id = ? and deleted_at IS NULL and week_number > 0", c.Param("client_id")).Find(&dietHistory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: diet does not exist for client_id %d", c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Errorf("error: could not fetch diet for client_id %s", c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// transform the results into given format
	var resRegularDiet []model.DietHistory
	var resDetoxDiet []model.DietHistory
	for _, diet := range dietHistory {
		if diet.DietType == 0 {
			// regular diet
			resRegularDiet = append(resRegularDiet, diet)
		} else if diet.DietType == 1 {
			// detox diet
			resDetoxDiet = append(resDetoxDiet, diet)
		}

	}

	c.JSON(http.StatusOK, gin.H{"diet_history_regular": resRegularDiet, "diet_history_detox": resDetoxDiet})
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
	//dietJSON, err := json.Marshal(schedule.Diet)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal diet to JSON"})
	//	return
	//}

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

	// fetch the week_number of the last diet sent
	var weekNumber int
	err := db.Model(&model.DietHistory{}).
		Where("client_id = ? and diet_type = ?", clientID, schedule.DietType).
		Select("week_number").
		Order("date DESC").
		Limit(1).
		Find(&weekNumber).
		Error
	if err != nil {
		fmt.Errorf("err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dietRecord := model.DietHistory{
		WeekNumber: weekNumber + 1,
		ClientID:   clientID,
		Date:       time.Now(),
		Weight:     nil,
		DietType:   schedule.DietType,
		DietString: &schedule.Diet,
	}
	if err := db.Create(&dietRecord).Error; err != nil {
		fmt.Errorf("error: SaveDietForClient | could not create empty diet_history_id %d for client_id %s | err: %v", schedule.Diet, clientID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//if err = db.Table("diet_histories").Where("id = ?", emptyDietRecord.ID).Update("diet", dietJSON).Error; err != nil {
	//	fmt.Errorf("error: SaveDietForClient | could not save diet for diet_history_id %d for client_id %s | err: %v", schedule.Diet, clientID, err.Error())
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Diet information saved successfully"})
	return
}

func EditDietForClient(c *gin.Context) {

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	// Parse the request body to extract the diet information
	var schedule model.EditDietForClientRequest
	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB
	clientID, _ := strconv.ParseUint(c.Param("client_id"), 10, 64)

	if err := db.Table("diet_histories").Where("id = ? and diet_type = ? and client_id = ?", schedule.DietID, schedule.DietType, clientID).Update("diet_string", schedule.Diet).Error; err != nil {
		fmt.Errorf("error: SaveDietForClient | could not save diet for diet_history_id %d for client_id %s | err: %v", schedule.Diet, c.Param("client_id"), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Diet information saved successfully"})
	return
}

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
	err := db.Model(model.DietHistory{}).Where("client_id = ? and diet_type = 0", clientID).Select("weight", "date").Find(&res).Error
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
	dietRecord := model.DietHistory{}
	err := db.Table("diet_histories").Where("client_id = ? and diet_type = 0", c.Param("client_id")).Order("date DESC").Select("id").First(&dietRecord).Error
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

	if err := db.Table("diet_histories").Where("id = ? and diet_type = 0", dietRecord.ID).Updates(model.DietHistory{
		Weight:   &req.Weight,
		Feedback: req.Feedback,
	}).Error; err != nil {
		fmt.Println("Error while saving client diet record", dietRecord)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	return
}

func DeleteDietForClientByAdmin(c *gin.Context) {
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

	// request contains the diet id to be deleted
	req := uint(0)
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB
	// verify that the given id exists and is the latest diet of that type
	diet := model.DietHistory{}
	err := db.Where("id = ? and client_id = ? and deleted_at IS NULL", req, clientID).Find(&diet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: could not find diet_history_id %d for client_id %s", req, clientID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Println("Could not retrieve diet record for client_id: " + c.Param("client_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// try to see if a more recent diet of that type exists
	var latestDiet model.DietHistory
	err = db.Model(&model.DietHistory{}).
		Where("client_id = ? and diet_type = ? and deleted_at IS NULL", clientID, diet.DietType).
		Order("date DESC, created_at DESC").
		Limit(1).
		First(&latestDiet).
		Error
	if err != nil {
		fmt.Errorf("error: could not find diet_history_id for client_id %s | err: %v", clientID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if latestDiet.ID != diet.ID {
		fmt.Errorf("error: trying to delete older diet, not allowed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad delete request"})
		return
	}
	timeNow := time.Now()
	diet.DeletedAt = &timeNow

	err = db.Model(&model.DietHistory{}).Where("id = ?", diet.ID).Update("deleted_at", timeNow).Error
	if err != nil {
		fmt.Errorf("error: could not delete diet_history_id for client_id %s | err: %v", clientID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func SaveCommonDietForClients(c *gin.Context) {

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	// Parse the request body to extract the diet information
	var req model.SaveCommonDietForClientsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	var clientIDs []uint64
	for _, val := range req.Groups {
		var groupClientIDs []uint64
		err := db.Model(&model.Client{}).Where("group_id = ?", val).Select("id").Find(&groupClientIDs).Error
		if err != nil {
			fmt.Errorf("error: could not find group clients for group_id %d: %s", val, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if len(groupClientIDs) == 0 {
			fmt.Errorf("error: could not find group clients for group_id %d", val)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no count"})
			return
		}
		clientIDs = append(clientIDs, groupClientIDs...)
	}

	var createDietReq []model.DietHistory
	for _, clientID := range clientIDs {
		// fetch the week_number of the last diet sent
		var weekNumber int
		err := db.Model(&model.DietHistory{}).
			Where("client_id = ? and diet_type = ?", clientID, req.DietType).
			Select("week_number").
			Order("date DESC").
			Limit(1).
			Find(&weekNumber).
			Error
		if err != nil {
			fmt.Errorf("err: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		createDietReq = append(createDietReq, model.DietHistory{
			WeekNumber: weekNumber + 1,
			ClientID:   clientID,
			Date:       time.Now(),
			Weight:     nil,
			DietType:   req.DietType,
			DietString: &req.Diet,
		})
	}

	if err := db.Create(&createDietReq).Error; err != nil {
		fmt.Errorf("error: SaveDietForClient | could not create diets %v for clients %v | err: %v", createDietReq, clientIDs, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Diet information saved successfully"})
	return
}
