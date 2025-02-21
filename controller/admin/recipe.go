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
	"strconv"
	"strings"
	"time"
)

func GetRecipeByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	recipe := model.Recipe{}
	if err := db.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Errorf("error: GetRecipeByID | recipe does not exist with id: %d", c.Param("id"))
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		fmt.Errorf("error: GetRecipeByID could not fetch recipe with id %d | err: %v", c.Param("id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ingredientsList := strings.Split(recipe.Ingredients, ";")
	prepList := strings.Split(recipe.Preparation, ";")

	res := model.GetRecipeResponse{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Ingredients: ingredientsList,
		Preparation: prepList,
	}

	c.JSON(http.StatusOK, gin.H{"recipe": res})
	return
}

func UpdateRecipeByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var recipeReq model.UpdateRecipeRequest
	if err := c.BindJSON(&recipeReq); err != nil {
		fmt.Errorf("error: UpdateRecipeByID | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredients := strings.Join(recipeReq.Ingredients, ";")
	steps := strings.Join(recipeReq.Preparation, ";")

	recipe := model.Recipe{
		ID:          recipeReq.ID,
		Name:        recipeReq.Name,
		Ingredients: ingredients,
		Preparation: steps,
	}
	db := database.DB
	if err := db.Model(&model.Recipe{}).Where("id = ?", c.Param("id")).Select("name", "ingredients", "preparation").Updates(&recipe).Error; err != nil {
		fmt.Errorf("error: UpdateRecipeByID | could not save recipe %v | err: %v", recipe, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func CreateRecipe(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var recipeReq model.CreateRecipeRequest
	if err := c.BindJSON(&recipeReq); err != nil {
		fmt.Errorf("error: CreateRecipe | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredients := ""
	for _, ingredient := range recipeReq.Ingredients {
		if ingredients == "" {
			ingredients = ingredient
		} else {
			ingredients = ingredients + ";" + ingredient
		}
	}

	steps := ""
	for _, prep := range recipeReq.Preparation {
		if steps == "" {
			steps = prep
		} else {
			steps = steps + ";" + prep
		}
	}

	recipe := model.Recipe{
		Name:        recipeReq.Name,
		Ingredients: ingredients,
		Preparation: steps,
	}
	db := database.DB
	if err := db.Save(&recipe).Error; err != nil {
		fmt.Errorf("error: CreateRecipe | could not save recipe %v | err: %v", recipe, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func DeleteRecipeByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meal ID"})
		return
	}

	db := database.DB
	if err := db.Model(&model.Recipe{}).Where("id = ?", recipeID).Update("deleted_at", time.Now()).Error; err != nil {
		fmt.Errorf("error: DeleteRecipeByMealID | could not delete recipe with id: %v | err: %v", c.Param("id"), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func GetListOfRecipes(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}
	db := database.DB
	var recipes []model.Recipe
	if err := db.Where("deleted_at IS NULL").Find(&recipes).Error; err != nil {
		fmt.Errorf("error: GetListOfRecipes | could not find recipes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := []model.GetListOfRecipesResponse{}
	for _, recipe := range recipes {
		res = append(res, model.GetListOfRecipesResponse{
			Name:     recipe.Name,
			RecipeID: recipe.ID,
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "list": res})
	return
}
