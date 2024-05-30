package controller

import (
	"fmt"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetExercisesForClient(c *gin.Context) {
	db := database.DB

	exercises := []model.Exercise{}
	err := db.Table("exercises").First(&exercises).Error
	if err != nil {
		fmt.Errorf("error: fetching all exercises: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercises": exercises})
	return
}

func GetExercisesForAdmin(c *gin.Context) {
	// authentication for admin
}

// GetExercise Admin
func GetExercise(c *gin.Context) {
	db := database.DB

	exercise := model.Exercise{}
	if err := db.Table("exercises").Where("exercise_id = ", c.Param("exercise_id")).First(&exercise).Error; err != nil {
		fmt.Errorf("error: fetching exercise with id: %s | error: %v", c.Param("exercise_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercise": exercise})
	return
}

// RemoveExerciseFromList Admin
func RemoveExerciseFromList(c *gin.Context) {
	db := database.DB

	exercise := model.Exercise{}
	err := db.Where("exercise_id = ", c.Param("exercise_id")).Delete(&exercise).Error
	if err != nil {
		fmt.Errorf("error: deleting exercise with id: %s | error: %v", c.Param("exercise_id"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
	return
}

// AddExerciseFromList Admin
func AddExerciseFromList(c *gin.Context) {
	db := database.DB

	exercise := model.Exercise{}
	if err := c.BindJSON(&exercise); err != nil {
		fmt.Errorf("error: request cannot be parsed %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&exercise).Error; err != nil {
		fmt.Errorf("error: exercise cannot be inserted in DB %v | struct: %+v", err, exercise)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
	return
}

// UpdateExerciseFromList Admin
func UpdateExerciseFromList(c *gin.Context) {
	db := database.DB

	exercise := model.Exercise{}
	if err := c.BindJSON(&exercise); err != nil {
		fmt.Errorf("error: request cannot be parsed %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("exercise_id"), 10, 32)
	if err != nil {
		fmt.Errorf("error: cannot parse exercise_id %s | error: %v", c.Param("exercise_id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercise.ID = uint(id)
	if err := db.Save(&exercise).Error; err != nil {
		fmt.Errorf("error: exercise cannot be inserted in DB %v | struct: %+v", err, exercise)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
	return
}
