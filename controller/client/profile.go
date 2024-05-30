package client

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// this method is triggered by the client
// can a client change their email?
func SaveProfileByClient(c *gin.Context) {
	db := database.DB

	req := model.Client{}
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(req)
	client := model.Client{}

	err := db.Table("clients").Where("email = ? and deleted_at IS NULL", req.Email).First(&client).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		fmt.Println("Record not found for this email: " + req.Email + " client_id: " + c.Param("client_id"))
		c.JSON(http.StatusNotFound, gin.H{"error": "record_not_found"})
		return
	} else if err != nil {
		fmt.Println("Error while fetching record from the DB for this email: "+req.Email+" client_id: "+c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client = migrateClientProfile(client, req)
	if err = db.Save(&client).Error; err != nil {
		fmt.Println("Error while saving client record", client)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	return
}

func GetProfileForClient(c *gin.Context) {
	db := database.DB

	client := model.Client{}
	err := db.Where("id = ? and deleted_at IS NULL", c.Param("client_id")).First(&client).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		fmt.Println("Record not found for this client_id: " + c.Param("client_id"))
		c.JSON(http.StatusNotFound, gin.H{"error": "record_not_found"})
		return
	} else if err != nil {
		fmt.Println("Error while fetching record from the DB for this client_id: "+c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := migrateClientProfile(model.Client{}, client)
	res.ID = client.ID
	c.JSON(http.StatusOK, gin.H{"response": res})
}

// only the data the client is allowed to update / retrieve is added to the response
// starting weight is allowed to be updated only during first time, implement how?
func migrateClientProfile(client, req model.Client) model.Client {
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Name != "" {
		client.Name = req.Name
	}
	if req.Age != 0 {
		client.Age = req.Age
	}
	if req.City != "" {
		client.City = req.City
	}
	if req.PhoneNumber != "" {
		client.PhoneNumber = req.PhoneNumber
	}
	if req.Height != 0 {
		client.Height = req.Height
	}
	if req.StartingWeight != 0.0 {
		client.StartingWeight = req.StartingWeight
	}
	if req.DietaryPreference != "" {
		client.DietaryPreference = req.DietaryPreference
	}
	if req.MedicalHistory != "" {
		client.MedicalHistory = req.MedicalHistory
	}
	if req.Allergies != "" {
		client.Allergies = req.Allergies
	}

	if req.Stay != "" {
		client.Stay = req.Stay
	}
	if req.Exercise != "" {
		client.Exercise = req.Exercise
	}
	if req.Comments != "" {
		client.Comments = req.Comments
	}
	if req.DietRecall != "" {
		client.DietRecall = req.DietRecall
	}
	if req.Locality != "" {
		client.Locality = req.Locality
	}
	return client
}
