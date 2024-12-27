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

// GetListOfExercises Used to populate the drop-down menu
func GetListOfExercises(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}
	db := database.DB

	var exercises []model.GetListOfExercisesResponse
	if err := db.Table("exercises").Find(&exercises).Error; err != nil {
		fmt.Errorf("error: GetListOfExercises | could not find exercises: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "exercises": exercises})
	return
}

func CreateExercise(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}
	var exercise model.Exercise
	if err := c.BindJSON(&exercise); err != nil {
		fmt.Errorf("error: CreateExercise | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	if err := db.Create(&exercise).Error; err != nil {
		fmt.Errorf("error: CreateExercise | could not create exercise: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "exercise": exercise})
}

func GetExerciseByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB
	var exercise model.Exercise
	err := db.First(&exercise, c.Param("exercise_id")).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Errorf("error: GetExerciseByID | exercise with id %d does not exist", c.Param("exercise_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		fmt.Errorf("error: GetExerciseByID | could not find exercise: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "exercise": exercise})
	return
}

func UpdateExerciseByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	var exercise model.Exercise
	if err := c.BindJSON(&exercise); err != nil {
		fmt.Errorf("error: UpdateExerciseByID | could not extract request from context | err : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB
	if err := db.Save(&exercise).Error; err != nil {
		fmt.Errorf("error: UpdateExerciseByID | could not update exercise %v | err: %v", exercise, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "exercise": exercise})
	return
}

func DeleteExerciseByID(c *gin.Context) {
	if !helpers.CheckUserType(c, "ADMIN") {
		fmt.Errorf("error: client user not allowed to access")
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized access by client"})
		return
	}

	db := database.DB
	if err := db.Where("id = ?", c.Param("exercise_id")).Update("deleted_at", time.Now()).Error; err != nil {
		fmt.Errorf("error: DeleteExerciseByID | could not delete exercise with id: %v| err: %v", c.Param("exercise_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}
