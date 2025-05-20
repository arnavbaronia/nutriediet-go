package api

import (
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"time"

	"github.com/cd-Ishita/nutriediet-go/controller"
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/cd-Ishita/nutriediet-go/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitiatePasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user model.UserAuth
	result := database.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset OTP has been sent"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate 4-digit OTP
	otp, err := generateOTP(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate OTP"})
		return
	}

	// Set reset token and expiration (15 minutes from now)
	expiration := time.Now().Add(15 * time.Minute)
	user.ResetToken = otp
	user.ResetTokenExp = expiration

	// Save to database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save reset token"})
		return
	}

	// Send email
	emailService := utils.NewEmailService()
	if err := emailService.SendPasswordResetEmail(input.Email, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset OTP has been sent"})
}

func CompletePasswordReset(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		OTP         string `json:"otp" binding:"required,len=4"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user with matching email and OTP
	var user model.UserAuth
	result := database.DB.Where("email = ? AND reset_token = ?", input.Email, input.OTP).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP or email"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if OTP is expired
	if time.Now().After(user.ResetTokenExp) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		return
	}

	// Hash new password (assuming you have a HashPassword function)
	hashedPassword, err := controller.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Update password and clear reset token
	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExp = time.Time{}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}
	return string(otp), nil
}
