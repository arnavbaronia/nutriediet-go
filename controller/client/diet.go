package client

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/constants"
	"gorm.io/gorm"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
)

func GetDietsForClient(c *gin.Context) {
	db := database.DB

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

	clientID := c.Param("client_id")

	var groupID int
	err := db.Model(&model.Client{}).Where("id = ? and deleted_at IS NULL", c.Param("client_id")).Find(&groupID).Error
	if err != nil {
		fmt.Errorf("error: GetDetoxDietForClient | could not fetch client information %d | err: %v", c.Param("client_id"), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	regularDiet, err := getDietForClient(&clientID, nil, constants.RegularDiet.Uint32())
	if !errors.Is(gorm.ErrRecordNotFound, err) && err != nil {
		fmt.Errorf("error: GetDietsForClient | could not fetch regular diet for client %s | err: %v", c.Param("client_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detoxDiet, err := getDietForClient(nil, &groupID, constants.DetoxDiet.Uint32())
	if !errors.Is(gorm.ErrRecordNotFound, err) && err != nil {
		fmt.Errorf("error: GetDietsForClient | could not fetch detox diet for client %s and groupID %d | err: %v", c.Param("client_id"), groupID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detoxWater, err := getDietForClient(nil, &groupID, constants.DetoxWater.Uint32())
	if !errors.Is(gorm.ErrRecordNotFound, err) && err != nil {
		fmt.Errorf("error: GetDietsForClient | could not fetch detox water for client %s and groupID %d | err: %v", c.Param("client_id"), groupID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"regular_diet": regularDiet, "detox_diet": detoxDiet, "detox_water": detoxWater})
	return
}

func getDietForClient(clientId *string, group_id *int, dietType uint32) (*string, error) {
	// Assuming you have a DB instance initialized elsewhere
	db := database.DB

	// Retrieve the latest diet history record for the client
	var diet sql.NullString
	if dietType == constants.RegularDiet.Uint32() && clientId != nil {
		err := db.Model(&model.DietHistory{}).
			Where("client_id = ? and diet_type = ?", *clientId, dietType).
			Order("date DESC, created_at DESC").
			Limit(1).
			Pluck("diet_string", &diet).
			Error
		if err != nil {
			fmt.Errorf("err: %v", err)
			return nil, err
		}
	} else if (dietType == constants.DetoxDiet.Uint32() || dietType == constants.DetoxWater.Uint32()) && group_id != nil {
		err := db.Model(&model.DietHistory{}).
			Where("group_id = ? and diet_type = ?", *group_id, dietType).
			Order("date DESC, created_at DESC").
			Limit(1).
			Pluck("diet_string", &diet).
			Error
		if err != nil {
			fmt.Errorf("err: %v", err)
			return nil, err
		}
	}

	if !diet.Valid {
		return nil, errors.New("no diet")
	}
	return &diet.String, nil
}
