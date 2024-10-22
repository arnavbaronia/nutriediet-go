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
	"time"
)

func GetRecipeByMealID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	recipe := model.Recipe{}
	if err := db.Where("meal_id = ?", c.Param("meal_id")).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Errorf("error: GetRecipeByMealIDForClient | recipe does not exist with meal_id: %d", c.Param("meal_id"))
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		fmt.Errorf("error: GetRecipeByMealIDForClient could not fetch recipe with meal_id %d | err: %v", c.Param("meal_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
	return
}

func UpdateRecipeByMealID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var recipeReq model.UpdateRecipeRequest
	if err := c.BindJSON(&recipeReq); err != nil {
		fmt.Errorf("error: UpdateRecipeByMealID | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredients := ""
	for _, ingredient := range recipeReq.Ingredients {
		ingredients = ingredients + ";" + ingredient
	}

	steps := ""
	for _, prep := range recipeReq.Preparation {
		steps = steps + ";" + prep
	}

	recipe := model.Recipe{
		ID:          recipeReq.ID,
		Name:        recipeReq.Name,
		Ingredients: ingredients,
		Preparation: steps,
	}
	db := database.DB
	if err := db.Save(&recipe).Error; err != nil {
		fmt.Errorf("error: UpdateRecipeByMealID | could not save recipe %v | err: %v", recipe, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func CreateRecipeByMealID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var recipeReq model.CreateRecipeRequest
	if err := c.BindJSON(&recipeReq); err != nil {
		fmt.Errorf("error: CreateRecipeByMealID | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredients := ""
	for _, ingredient := range recipeReq.Ingredients {
		ingredients = ingredients + ";" + ingredient
	}

	steps := ""
	for _, prep := range recipeReq.Preparation {
		steps = steps + ";" + prep
	}

	recipe := model.Recipe{
		Name:        recipeReq.Name,
		Ingredients: ingredients,
		Preparation: steps,
	}
	db := database.DB
	if err := db.Save(&recipe).Error; err != nil {
		fmt.Errorf("error: CreateRecipeByMealID | could not save recipe %v | err: %v", recipe, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func DeleteRecipeByMealID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB
	if err := db.Where("meal_id = ?", c.Param("meal_id")).Update("deleted_at", time.Now()).Error; err != nil {
		fmt.Errorf("error: DeleteRecipeByMealID | could not delete recipe with meal_id: %v | err: %v", c.Param("meal_id"), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}
