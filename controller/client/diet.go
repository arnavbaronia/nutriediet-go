package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/constants"
	"gorm.io/gorm"

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
	err := db.Model(&model.Client{}).Where("id = ? and deleted_at IS NULL", c.Param("client_id")).Select("group_id").Find(&groupID).Error
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

	c.JSON(http.StatusOK, gin.H{
		"regular_diet": regularDiet,
		"detox_diet":   detoxDiet,
		"detox_water":  detoxWater,
	})
	return
}

func getDietForClient(clientId *string, group_id *int, dietType uint32) (*model.DietHistoryResponse, error) {
	db := database.DB

	var dietHistory model.DietHistoryResponse
	if dietType == constants.RegularDiet.Uint32() && clientId != nil {
		err := db.Model(&model.DietHistory{}).
			Joins("left outer join diet_templates on diet_template_id = diet_templates.id").
			Select("diet_histories.*, diet_templates.name as diet_template_name").
			Where("client_id = ? and diet_type = ? and diet_histories.deleted_at IS NULL", *clientId, dietType).
			Order("date DESC, created_at DESC").
			Limit(1).
			First(&dietHistory).
			Error
		if err != nil {
			fmt.Errorf("err: %v", err)
			return nil, err
		}
	} else if (dietType == constants.DetoxDiet.Uint32() || dietType == constants.DetoxWater.Uint32()) && group_id != nil {
		err := db.Model(&model.DietHistory{}).
			Joins("left outer join diet_templates on diet_template_id = diet_templates.id").
			Select("diet_histories.*, diet_templates.name as diet_template_name").
			Where("group_id = ? and diet_type = ? and diet_histories.deleted_at IS NULL", *group_id, dietType).
			Order("date DESC, created_at DESC").
			Limit(1).
			First(&dietHistory).
			Error
		if err != nil {
			fmt.Errorf("err: %v", err)
			return nil, err
		}
	}

	if dietHistory.DietString == nil {
		fmt.Errorf("found invalid diet for client %s or group %d of diet_type %d", clientId, group_id, dietType)
		return nil, nil
	}
	return &dietHistory, nil
}
