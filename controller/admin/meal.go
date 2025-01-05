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
	"strings"
)

func GetMealList(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	var mealList []model.MealAdditionalInfo
	if err := db.Where("type = MEAL").Find(&mealList).Error; err != nil {
		fmt.Errorf("error: GetMealList failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mealList})
	return
}

func GetQuantityList(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	var quantityList []model.MealAdditionalInfo
	if err := db.Where("type = QUANTITY").Find(&quantityList).Error; err != nil {
		fmt.Errorf("error: GetQuantityList failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": quantityList})
	return
}

func CreateNewMeal(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	mealReq := model.CreateNewMealRequest{}
	if err := c.BindJSON(&mealReq); err != nil {
		fmt.Errorf("error: CreateRecipeByMealID | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meal := model.MealAdditionalInfo{
		Name: strings.ToLower(mealReq.Name),
		Type: "MEAL",
	}

	db := database.DB
	if err := db.Where("name = ? and type = ?", strings.ToUpper(mealReq.Name), "MEAL").First(&model.MealAdditionalInfo{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: CreateNewMeal | meal with same name already exists %v | err : %v", mealReq, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := db.Create(&meal).Error; err != nil {
		fmt.Errorf("error: CreateNewMeal | error saving meal %v | err : %v", meal, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if mealReq.HasRecipe {
		ingredients, steps := "", ""
		for _, ingredient := range mealReq.Ingredients {
			ingredients = ingredients + ";" + ingredient
		}

		for _, prep := range mealReq.Preparation {
			steps = steps + ";" + prep
		}

		if err := db.Create(&model.Recipe{
			MealID:      int(meal.ID),
			Ingredients: ingredients,
			Preparation: steps,
		}).Error; err != nil {
			fmt.Errorf("error: CreateNewMeal | error saving recipe %v | err : %v", meal, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": mealReq})
	return
}
