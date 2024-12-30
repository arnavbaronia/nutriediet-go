package admin

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMealList(c *gin.Context) {
	db := database.DB

	mealList := []model.MealAdditionalInfo{}
	if err := db.Where("type = MEAL").Find(&mealList).Error; err != nil {
		fmt.Errorf("error: GetMealList failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mealList})
	return
}

func GetQuantityList(c *gin.Context) {
	db := database.DB

	quantityList := []model.MealAdditionalInfo{}
	if err := db.Where("type = QUANTITY").Find(&quantityList).Error; err != nil {
		fmt.Errorf("error: GetQuantityList failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": quantityList})
	return
}
