package client

import (
	"fmt"
	"net/http"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/middleware"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
)

func GetExercisesForClient(c *gin.Context) {
	clientEmail := c.GetString("email")
	clientID := c.Param("client_id")
	isAllowed, isActive := middleware.ClientAuthentication(clientEmail, clientID)
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": clientEmail, "requestClientID": clientID})
		return
	}

	if !isActive {
		c.JSON(http.StatusOK, gin.H{"isActive": false})
		return
	}

	db := database.DB

	// Get all exercises
	var exercises []model.Exercise
	if err := db.Table("exercises").Where("deleted_at IS NULL").Find(&exercises).Error; err != nil {
		fmt.Errorf("error: fetching all exercises: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get favorite exercise IDs for this client
	var favoriteExerciseIDs []uint
	if err := db.Table("favorite_exercises").
		Where("client_id = ?", clientID).
		Pluck("exercise_id", &favoriteExerciseIDs).Error; err != nil {
		fmt.Errorf("error: fetching favorite exercises: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a map for quick lookup
	favoritesMap := make(map[uint]bool)
	for _, id := range favoriteExerciseIDs {
		favoritesMap[id] = true
	}

	// Mark exercises as favorite
	type ExerciseResponse struct {
		model.Exercise
		IsFavorite bool `json:"is_favorite"`
	}

	var response []ExerciseResponse
	for _, exercise := range exercises {
		response = append(response, ExerciseResponse{
			Exercise:   exercise,
			IsFavorite: favoritesMap[exercise.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{"isActive": true, "exercises": response})
	return
}

func ToggleFavoriteExercise(c *gin.Context) {
	clientEmail := c.GetString("email")
	clientID := c.Param("client_id") // This is a string
	isAllowed, _ := middleware.ClientAuthentication(clientEmail, clientID)
	if !isAllowed {
		c.JSON(http.StatusUnauthorized, gin.H{"clientEmail": clientEmail, "requestClientID": clientID})
		return
	}

	var req struct {
		ExerciseID uint `json:"exercise_id"`
		IsFavorite bool `json:"is_favorite"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	if req.IsFavorite {
		// Add to favorites
		favorite := model.FavoriteExercise{
			ClientID:   clientID, // Now using string directly
			ExerciseID: req.ExerciseID,
		}
		if err := db.Create(&favorite).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Remove from favorites
		if err := db.Where("client_id = ? AND exercise_id = ?", clientID, req.ExerciseID).
			Delete(&model.FavoriteExercise{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
