package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cd-Ishita/nutriediet-go/constants"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	var activeClients []model.ClientMiniInfo
	var inactiveClients []model.ClientMiniInfo
	for _, res := range lastDietDates {
		index := clientIDMap[res.ClientID]
		clients[index].LastDietDate = res.Date

		if clients[index].IsActive {
			activeClients = append(activeClients, clients[index])
		} else {
			inactiveClients = append(inactiveClients, clients[index])
		}
	}

	c.JSON(http.StatusOK, gin.H{"active_clients": activeClients, "inactive_clients": inactiveClients})
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
	err = db.Table("diet_histories").Where("client_id = ?", c.Param("client_id")).Find(&dietHistories).Error
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

	isSuperAdmin := c.GetString("email") == constants.SuperAdminEmail
	upsertedClient, err := migrateClientInfoForAdmin(req, client, isSuperAdmin)
	if err != nil {
		fmt.Errorf("error: failed to migrate client info | %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.Save(&upsertedClient).Error
	if err != nil {
		fmt.Errorf("error: could not save client information | client_info: %v | err: %v", upsertedClient, err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"client": upsertedClient})
	return
}

func migrateClientInfoForAdmin(updatedInfo model.Client, existingInfo model.Client, isSuperAdmin bool) (model.Client, error) {
	// TODO: do we want admin to be able to update the starting weight in cases where client comes back

	if isSuperAdmin {
		if _, exists := constants.PackageDayMap[updatedInfo.Package]; exists && updatedInfo.Package != "" {
			existingInfo.Package = updatedInfo.Package
		}
		// Add validation for empty package
		if updatedInfo.Package == "" {
			return model.Client{}, errors.New("package cannot be empty")
		}

		if updatedInfo.TotalAmount != 0 {
			existingInfo.TotalAmount = updatedInfo.TotalAmount
		}
		if updatedInfo.AmountPaid != 0 {
			existingInfo.AmountPaid = updatedInfo.AmountPaid
		}
		existingInfo.AmountDue = existingInfo.TotalAmount - existingInfo.AmountPaid

		// Update lastPaymentDate if provided
		if updatedInfo.LastPaymentDate != nil {
			existingInfo.LastPaymentDate = updatedInfo.LastPaymentDate
			// Only calculate the next payment date if not explicitly provided
			if updatedInfo.NextPaymentDate == nil && existingInfo.Package != "" {
				nextPaymentDate := existingInfo.LastPaymentDate.AddDate(0, 0, constants.PackageDayMap[existingInfo.Package])
				existingInfo.NextPaymentDate = &nextPaymentDate
			}
		}

		// Update nextPaymentDate if explicitly provided
		if updatedInfo.NextPaymentDate != nil {
			existingInfo.NextPaymentDate = updatedInfo.NextPaymentDate
		}

		// Update DateOfJoining if provided
		if updatedInfo.DateOfJoining != nil {
			existingInfo.DateOfJoining = updatedInfo.DateOfJoining
		}
	}

	if updatedInfo.Name != "" {
		existingInfo.Name = updatedInfo.Name
	}
	if updatedInfo.Age != 0 {
		existingInfo.Age = updatedInfo.Age
	}
	if updatedInfo.City != "" {
		existingInfo.City = updatedInfo.City
	}
	if updatedInfo.PhoneNumber != "" {
		existingInfo.PhoneNumber = updatedInfo.PhoneNumber
	}
	if updatedInfo.Remarks != "" {
		existingInfo.Remarks = updatedInfo.Remarks
	}
	if updatedInfo.DietitianId != 0 {
		existingInfo.DietitianId = updatedInfo.DietitianId
	}
	if updatedInfo.GroupID != 0 {
		existingInfo.GroupID = updatedInfo.GroupID
	}
	if updatedInfo.Height != 0 {
		existingInfo.Height = updatedInfo.Height
	}
	if updatedInfo.StartingWeight != 0 {
		existingInfo.StartingWeight = updatedInfo.StartingWeight
	}
	if updatedInfo.DietaryPreference != "" {
		existingInfo.DietaryPreference = updatedInfo.DietaryPreference
	}
	if updatedInfo.MedicalHistory != "" {
		existingInfo.MedicalHistory = updatedInfo.MedicalHistory
	}
	if updatedInfo.Allergies != "" {
		existingInfo.Allergies = updatedInfo.Allergies
	}
	if updatedInfo.Stay != "" {
		existingInfo.Stay = updatedInfo.Stay
	}
	if updatedInfo.Exercise != "" {
		existingInfo.Exercise = updatedInfo.Exercise
	}
	if updatedInfo.Comments != "" {
		existingInfo.Comments = updatedInfo.Comments
	}
	if updatedInfo.DietRecall != "" {
		existingInfo.DietRecall = updatedInfo.DietRecall
	}
	if updatedInfo.Locality != "" {
		existingInfo.Locality = updatedInfo.Locality
	}

	return existingInfo, nil
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
	clientID, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	// Check if the client exists
	client := model.Client{}
	err = db.Table("clients").Where("id = ?", clientID).First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: client does not exist with id %d", clientID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "client not found"})
		return
	} else if err != nil {
		fmt.Errorf("error: could not fetch client with id %d | %v", clientID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	err = db.Table("clients").Where("id = ?", clientID).UpdateColumn("is_active", !client.IsActive).Error
	if err != nil {
		fmt.Errorf("error: could not update activation value for client with id %d | err: %v", clientID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
