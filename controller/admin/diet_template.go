package admin

import (
	"encoding/json"
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

func GetDietTemplatesList(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	dietTemplates := []model.DietTemplate{}
	err := db.Where("deleted_at IS NULL").Select("id", "name").Find(&dietTemplates).Error
	if err != nil {
		fmt.Errorf("error: could not fetch diet templates for GetDietTemplatesList API | err: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": dietTemplates})
	return
}

func GetDietTemplateByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	var dietTemplateJSON []string
	err := db.Model(&model.DietTemplate{}).Where("id = ? and deleted_at IS NULL", c.Param("diet_template_id")).Pluck("diet", &dietTemplateJSON).Error
	if err != nil {
		fmt.Errorf("error: could not fetch dietTemplate with id: %s for GetDietTemplateByID | err: %v", c.Param("diet_template_id"), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	dietTemplate := model.DietSchedule{}
	err = json.Unmarshal([]byte(dietTemplateJSON[0]), &dietTemplate)
	if err != nil {
		fmt.Errorf("error: could not unmarshal diet template with id: %s for GetDietTemplateByID | err: %v", c.Param("diet_template_id"), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": dietTemplate})
}

func CreateDietTemplate(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var template model.CreateDietTemplateRequest
	if err := c.BindJSON(&template); err != nil {
		fmt.Errorf("error: could not extract request from context for CreateDietTemplate | err: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	err := db.Table("diet_templates").Where("deleted_at IS NULL and name = ?", template.Name).First(&model.DietTemplate{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// continue
	} else if err != nil {
		fmt.Errorf("error: CreateDietTemplate | could not check for existing diet template with name | err: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		fmt.Errorf("error: CreateDietTemplate | already exists diet template with name: %s", template.Name)
		c.JSON(http.StatusConflict, gin.H{"error": "diet template already exists"})
		return
	}

	dietTemplate := model.DietTemplate{
		Name:       template.Name,
		DietString: &template.Diet,
	}
	err = db.Table("diet_templates").Save(&dietTemplate).Error
	if err != nil {
		fmt.Errorf("error: could not save dietTemplate %v for CreateDietTemplate | err: %v", template, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//dietJSON, err := json.Marshal(template.Diet)
	//if err != nil {
	//	fmt.Errorf("error: could not marshal diet to json for CreateDietTemplate | err: %v", err.Error())
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal diet to JSON"})
	//	return
	//}
	//
	//if err := db.Table("diet_templates").Where("id = ?", dietTemplate.ID).Update("diet", dietJSON).Error; err != nil {
	//	fmt.Errorf("error: could not save template for CreateDietTemplate | err: %v", err.Error())
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save template for CreateDietTemplate"})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func UpdateDietTemplate(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var template model.UpdateDietTemplateRequest
	if err := c.BindJSON(&template); err != nil {
		fmt.Errorf("error: could not extract request from context for UpdateDietTemplateByID | err: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dietTemplate := model.DietTemplate{
		Name:       template.Name,
		DietString: &template.Diet,
		ID:         template.ID,
	}
	db := database.DB
	err := db.Table("diet_templates").Where("id = ?", c.Param("diet_template_id")).Select("name", "diet").Updates(&dietTemplate).Error
	if err != nil {
		fmt.Errorf("error: could not update dietTemplate %v for UpdateDietTemplateByID | err: %v", template, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func DeleteDietTemplateByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB

	err := db.Table("diet_templates").Where("id = ?", c.Param("diet_template_id")).Update("deleted_at", time.Now()).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: diet template with id %s does not exist in DeleteDietTemplateByID", c.Param("diet_template_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"err": "record not found"})
		return
	} else if err != nil {
		fmt.Errorf("error: diet template with id %s could not be marked deleted in DeleteDietTemplateByID", c.Param("diet_template_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
