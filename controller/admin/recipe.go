package admin

import (
	"errors"
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
)

//func GetRecipeByID(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	db := database.DB
//
//	recipe := model.Recipe{}
//	if err := db.Where("id = ?", c.Param("id")).First(&recipe).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			fmt.Errorf("error: GetRecipeByID | recipe does not exist with id: %d", c.Param("id"))
//			c.JSON(http.StatusNotFound, gin.H{"error": err})
//			return
//		}
//		fmt.Errorf("error: GetRecipeByID could not fetch recipe with id %d | err: %v", c.Param("id"), err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	ingredientsList := strings.Split(recipe.Ingredients, ";")
//	prepList := strings.Split(recipe.Preparation, ";")
//
//	res := model.GetRecipeResponse{
//		ID:          recipe.ID,
//		Name:        recipe.Name,
//		Ingredients: ingredientsList,
//		Preparation: prepList,
//	}
//
//	c.JSON(http.StatusOK, gin.H{"recipe": res})
//	return
//}
//
//func UpdateRecipeByID(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	var recipeReq model.UpdateRecipeRequest
//	if err := c.BindJSON(&recipeReq); err != nil {
//		fmt.Errorf("error: UpdateRecipeByID | could not extract request from context | err : %v", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	ingredients := strings.Join(recipeReq.Ingredients, ";")
//	steps := strings.Join(recipeReq.Preparation, ";")
//
//	recipe := model.Recipe{
//		ID:          recipeReq.ID,
//		Name:        recipeReq.Name,
//		Ingredients: ingredients,
//		Preparation: steps,
//	}
//	db := database.DB
//	if err := db.Model(&model.Recipe{}).Where("id = ?", c.Param("id")).Select("name", "ingredients", "preparation").Updates(&recipe).Error; err != nil {
//		fmt.Errorf("error: UpdateRecipeByID | could not save recipe %v | err: %v", recipe, err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true})
//	return
//}
//
//func CreateRecipe(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	var recipeReq model.CreateRecipeRequest
//	if err := c.BindJSON(&recipeReq); err != nil {
//		fmt.Errorf("error: CreateRecipe | could not extract request from context | err : %v", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	ingredients := ""
//	for _, ingredient := range recipeReq.Ingredients {
//		if ingredients == "" {
//			ingredients = ingredient
//		} else {
//			ingredients = ingredients + ";" + ingredient
//		}
//	}
//
//	steps := ""
//	for _, prep := range recipeReq.Preparation {
//		if steps == "" {
//			steps = prep
//		} else {
//			steps = steps + ";" + prep
//		}
//	}
//
//	recipe := model.Recipe{
//		Name:        recipeReq.Name,
//		Ingredients: ingredients,
//		Preparation: steps,
//	}
//	db := database.DB
//	if err := db.Save(&recipe).Error; err != nil {
//		fmt.Errorf("error: CreateRecipe | could not save recipe %v | err: %v", recipe, err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true})
//	return
//}
//
//func DeleteRecipeByID(c *gin.Context) {
//	if !helpers.CheckUserType(c, "ADMIN") {
//		fmt.Errorf("error: client user not allowed to access")
//		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
//		return
//	}
//
//	recipeID, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meal ID"})
//		return
//	}
//
//	db := database.DB
//	if err := db.Model(&model.Recipe{}).Where("id = ?", recipeID).Update("deleted_at", time.Now()).Error; err != nil {
//		fmt.Errorf("error: DeleteRecipeByMealID | could not delete recipe with id: %v | err: %v", c.Param("id"), err.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true})
//	return
//}

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

func UploadRecipeImage(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
		return
	}

	imageName := c.PostForm("name") // this is from the "name" field
	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image name"})
		return
	}

	// Ensure "images" folder exists
	os.MkdirAll("images", os.ModePerm)

	// Save file inside "images" directory
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	savePath := fmt.Sprintf("images/%s", filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return URL to access image later
	imageURL := fmt.Sprintf("/images/%s", file.Filename)

	recipe := model.Recipe{
		Name:     imageName,
		ImageURL: imageURL,
	}

	db := database.DB

	if err := db.Table("recipes").Save(&recipe).Error; err != nil {
		fmt.Errorf("error: UploadRecipeImage failed, could not save recipe to DB | err: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload successful",
		"url":     imageURL,
	})
}

func GetRecipeImageForAdmin(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	var recipe model.Recipe
	err := db.Table("recipes").Where("id = ? and deleted_at IS NULL", c.Param("recipe_id")).First(&recipe).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		fmt.Errorf("error: GetRecipeImageForAdmin recipe with id %s does not exist", c.Param("recipe_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"err": "no recipe found"})
		return
	} else if err != nil {
		fmt.Errorf("error: GetRecipeImageForAdmin could not fetch recipe with id %s | err: %v", c.Param("recipe_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
	return
}
