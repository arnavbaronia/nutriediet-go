package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
)

const (
	TypeDiet = "0"

	TypeDetoxDiet = "2"
)

func GetRegularDietForClient(c *gin.Context) {
	clientID := c.Param("client_id")

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

	diet, err := getDietForClient(clientID, TypeDiet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch diet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isActive": true, "diet": diet})
}

func GetDetoxDietForClient(c *gin.Context) {
	isAllowed, isActive := middleware.ClientAuthentication(c.Param("email"), c.Param("client_id"))
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": c.Param("email"), "requestClientID": c.Param("client_id")})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	diet, err := getDietForClient(c.Param("client_id"), TypeDetoxDiet)
	if err != nil {
		fmt.Errorf("error finding diet for client_id: %s", c.Param("client_id"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"isActive": true, "diet": diet})
	return
}

func getDietForClient(clientId, dietType string) (model.DietSchedule, error) {
	// Assuming you have a DB instance initialized elsewhere
	db := database.DB

	// Retrieve the latest diet history record for the client

	var dietJSON []string
	err := db.Model(&model.DietHistory{}).
		Where("client_id = ? and diet_type = ?", clientId, dietType).
		Order("date DESC").
		Pluck("diet", &dietJSON).
		Error
	if err != nil {
		return model.DietSchedule{}, err
	}

	dietFinal := model.DietSchedule{}
	err = json.Unmarshal([]byte(dietJSON[0]), &dietFinal)
	if err != nil {
		fmt.Println("Error:", err)
		return model.DietSchedule{}, err
	}

	//err := db.Where("client_id = ? and diet_type = ?", clientId, dietType).Order("date DESC").First(&dietHistory).Error
	//if errors.Is(gorm.ErrRecordNotFound, err) {
	//	fmt.Errorf("Error: RecordNotFound for client_id: " + clientId + " diet_type: " + dietType)
	//	return model.DietSchedule{}, err
	//} else {
	//	fmt.Errorf("Error: Error finding diet record for client_id: " + clientId + " diet_type: " + dietType)
	//	return model.DietSchedule{}, err
	//}

	// Extract the schedule from the diet history record
	return dietFinal, nil
}

// FUTURE: check if meal id is applicable for client before fetching
func GetRecipeForMealID(c *gin.Context) {
	// pull the recipe
	//db := database.DB
	//
	//err := db.Where("")
}
