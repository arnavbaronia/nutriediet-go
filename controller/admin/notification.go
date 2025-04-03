package admin

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/constants"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateNewMotivation(c *gin.Context) {
	db := database.DB

	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	req := model.CreateNotificationReq{}
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("Wrong request, cannot be extracted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notif := model.Notification{
		Type:          constants.Motivation,
		Text:          req.Text,
		PostingActive: req.PostingActive,
	}

	if err := db.Create(&notif).Error; err != nil {
		fmt.Errorf("error: CreateNewMotivation could not create notif %v | err: %v", notif, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"motivation": notif})
	return
}

func UnpostMotivation(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	postOrUnpostMotivation(false, c)
	return
}

func PostMotivation(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	postOrUnpostMotivation(true, c)
	return
}

func postOrUnpostMotivation(post bool, c *gin.Context) {

	db := database.DB
	var notif model.Notification
	if err := db.Model(model.Notification{}).Where("id = ? and type = ? and deleted_at IS NULL", c.Param("motivation_id"), constants.Motivation).Find(&notif).Error; err != nil {
		fmt.Errorf("error: UpdateMotivation | could not find motivation with id %d | err: %v", c.Param("motivation_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notif.PostingActive = post

	if err := db.Save(&notif).Error; err != nil {
		fmt.Errorf("error: UpdateMotivation could not update notif %v | err: %v", notif, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"motivation": notif})
	return
}

func GetAllMotivations(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	var notifs []model.Notification
	err := db.Model(&model.Notification{}).Where("type = ? and deleted_at IS NULL", constants.Motivation).Find(&notifs).Error
	if err != nil {
		fmt.Errorf("error: GetAllMotivations | could not fetch motivations | err: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"motivation": notifs})
	return
}
