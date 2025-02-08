package client

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// this method is triggered by the client
// TODO: decide is client should be allowed to update profile in deactivated account
func UpdateProfileByClient(c *gin.Context) {
	clientEmail := c.GetString("email")
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	db := database.DB

	req := model.Client{}
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	c.JSON(http.StatusOK, gin.H{"isActive": true})
	return
}

func GetProfileForClient(c *gin.Context) {
	clientEmail := c.GetString("email")
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

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
	c.JSON(http.StatusOK, gin.H{"isActive": true, "response": res})
}

// only the data the client is allowed to update / retrieve is added to the response
// starting weight is allowed to be updated only during first time, implement how?
func migrateClientProfile(client, req model.Client) model.Client {

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

func CreateProfileByClient(c *gin.Context) {
	db := database.DB

	req := model.Client{}
	if err := c.BindJSON(&req); err != nil {
		fmt.Errorf("Wrong request, cannot be extracted. For client_id: " + c.Param("client_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email != c.Param("email") {
		fmt.Errorf("error: Email does not match | request: %v | headers: %v", req, c.Param("email"))
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email does not match")})
		return
	}

	user := model.UserAuth{}
	err := db.Where("email = ?", c.Param("email")).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: UserAuth record not found for client with email: %s", c.Param("email"))
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New("user has not signed up before profile creation")})
		return
	} else if err != nil {
		fmt.Errorf("error: could not fetch user auth information for client with email %s | err: %v", c.Param("email"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	client := migrateClientProfileByClientUpdate(req)

	client.Email = user.Email
	client.IsActive = false
	timeNow := time.Now()
	client.DateOfJoining = &timeNow
	if err = db.Save(&client).Error; err != nil {
		fmt.Errorf("error: could not save the client's profile information in database | email: %s | err: %v", c.Param("email"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isActive": false, "client": client})
	return
}

func migrateClientProfileByClientUpdate(updatedInfo model.Client) model.Client {
	client := model.Client{
		Name:              updatedInfo.Name,
		Age:               updatedInfo.Age,
		City:              updatedInfo.City,
		PhoneNumber:       updatedInfo.PhoneNumber,
		Remarks:           updatedInfo.Remarks,
		Height:            updatedInfo.Height,
		StartingWeight:    updatedInfo.StartingWeight,
		DietaryPreference: updatedInfo.DietaryPreference,
		MedicalHistory:    updatedInfo.MedicalHistory,
		Allergies:         updatedInfo.Allergies,
		Stay:              updatedInfo.Stay,
		Exercise:          updatedInfo.Exercise,
		Comments:          updatedInfo.Comments,
		DietRecall:        updatedInfo.DietRecall,
		Locality:          updatedInfo.Locality,
	}
	return client
}

func HasClientCreatedProfile(c *gin.Context) {
	db := database.DB

	client := model.Client{}
	err := db.Where("id = ?", c.Param("client_id")).First(client).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		c.JSON(http.StatusOK, gin.H{"profile_created": false, "is_active": false})
		return
	} else if err != nil {
		fmt.Errorf("error: could not fetch client's profile information in database client_id: %s | err: %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"profile_created": true, "is_active": client.IsActive})
	return
}
