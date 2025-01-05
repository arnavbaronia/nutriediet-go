package admin

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/constants"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetAllClients(c *gin.Context) {
	db := database.DB

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	clients := []model.ClientMiniInfo{}
	err := db.Table("clients").Find(&clients).Error
	if err != nil {
		fmt.Errorf("error: could not find all clients | %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	clientIDs := []uint64{}
	clientIDMap := map[uint64]int{}

	for index, client := range clients {
		clientIDs = append(clientIDs, client.ID)
		clientIDMap[client.ID] = index
	}

	lastDietDates := []model.DietHistory{}

	err = db.Table("diet_histories AS d").
		Select("d.client_id, d.date").
		Joins("JOIN (SELECT client_id, MAX(date) as max_date FROM diet_histories WHERE client_id IN (?) GROUP BY client_id) AS sub ON d.client_id = sub.client_id AND d.date = sub.max_date", clientIDs).
		Find(&lastDietDates).Error
	if err != nil {
		fmt.Errorf("error: could not find the last diet dates client_id: %v | %v", clientIDs, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for _, res := range lastDietDates {
		index := clientIDMap[res.ClientID]
		clients[index].LastDietDate = res.Date
	}

	c.JSON(http.StatusOK, gin.H{"clients": clients})
	return
}

func GetClientInfo(c *gin.Context) {
	db := database.DB

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Println("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	client := model.Client{}
	err := db.Table("clients").Where("id = ?", c.Param("client_id")).First(&client).Error
	if err != nil {
		fmt.Errorf("error: could not fetch client with id %s | %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	dietHistories := []model.DietHistory{}
	err = db.Table("diet_histories").Where("client_id = ?", c.Param("client_id")).Select("id", "week_number").Find(&dietHistories).Error
	if err != nil {
		fmt.Errorf("error: could not fetch number of rows for client_id %s | %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client": client, "diets": dietHistories})
	return
}

func UpdateClientInfo(c *gin.Context) {
	db := database.DB

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Println("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	req := model.Client{}
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := model.Client{}
	err := db.Table("clients").Where("id = ?", c.Param("client_id")).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: client does not exist with id %s", c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		fmt.Errorf("error: could not fetch client with id %s | %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	upsertedClient := migrateClientInfoForAdmin(req, client)
	err = db.Save(&upsertedClient).Error
	if err != nil {
		fmt.Errorf("error: could not save client information | client_info: %v | err: %v", upsertedClient, err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"client": upsertedClient})
	return
}

func migrateClientInfoForAdmin(updatedInfo model.Client, existingInfo model.Client) model.Client {
	// TODO: do we want admin to be able to update the starting weight in cases where client comes back
	client := model.Client{
		ID:                existingInfo.ID,
		Name:              updatedInfo.Name,
		Age:               updatedInfo.Age,
		City:              updatedInfo.City,
		PhoneNumber:       updatedInfo.PhoneNumber,
		DateOfJoining:     updatedInfo.DateOfJoining,
		Package:           updatedInfo.Package,
		AmountPaid:        updatedInfo.AmountPaid,
		LastPaymentDate:   updatedInfo.LastPaymentDate,
		NextPaymentDate:   updatedInfo.NextPaymentDate, // should be computed field
		Remarks:           updatedInfo.Remarks,
		DietitianId:       updatedInfo.DietitianId,
		Group:             updatedInfo.Group,
		Email:             existingInfo.Email,
		Height:            updatedInfo.Height,
		StartingWeight:    existingInfo.StartingWeight,
		DietaryPreference: updatedInfo.DietaryPreference,
		MedicalHistory:    updatedInfo.MedicalHistory,
		Allergies:         updatedInfo.Allergies,
		Stay:              updatedInfo.Stay,
		Exercise:          updatedInfo.Exercise,
		Comments:          updatedInfo.Comments,
		DietRecall:        updatedInfo.DietRecall,
		IsActive:          updatedInfo.IsActive,
		Locality:          updatedInfo.Locality,
		CreatedAt:         updatedInfo.CreatedAt,
	}
	client.NextPaymentDate = client.LastPaymentDate.AddDate(0, 0, constants.PackageDayMap[updatedInfo.Package])
	return client
}

// deactivation of client account handled by a separate API
func ActivateOrDeactivateClientAccount(c *gin.Context) {
	db := database.DB

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Println("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	// Check if user exists
	client := model.Client{}
	err := db.Table("clients").Where("id = ?", c.Param("client_id")).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: client does not exist with id %s", c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Errorf("error: could not fetch client with id %s | %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	err = db.Where("client_id = ?", c.Param("client_id")).UpdateColumn("is_active", !client.IsActive).Error
	if err != nil {
		fmt.Errorf("error: could not update activation value for client with id %s | err: %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}
