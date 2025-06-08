package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/helpers"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request body for reset password
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ForgotPassword handles the forgot password request
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	db := database.DB

	// Check if user exists
	var user model.UserAuth
	err := db.Where("email = ?", req.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found with this email address",
		})
		return
	} else if err != nil {
		fmt.Printf("Database error while checking user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error occurred",
		})
		return
	}

	// Generate OTP
	otp, err := helpers.GenerateOTP()
	if err != nil {
		fmt.Printf("Error generating OTP: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate OTP",
		})
		return
	}

	// Hash OTP
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing OTP: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process OTP",
		})
		return
	}

	// Set expiry time (5 minutes from now)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Store OTP in database using FirstOrCreate with Assign
	passwordOTP := model.PasswordOTP{
		Email:     req.Email,
		OtpHash:   string(otpHash),
		ExpiresAt: expiresAt,
	}

	result := db.Where(model.PasswordOTP{Email: req.Email}).
		Assign(model.PasswordOTP{OtpHash: string(otpHash), ExpiresAt: expiresAt}).
		FirstOrCreate(&passwordOTP)

	if result.Error != nil {
		fmt.Printf("Error saving OTP to database: %v\n", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store OTP",
		})
		return
	}

	// Send OTP via email
	err = helpers.SendOTPEmail(req.Email, otp)
	if err != nil {
		fmt.Printf("Error sending OTP email: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send OTP email. Please check your email configuration.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully to your email address",
		"email":   req.Email,
	})
}

// ResetPassword handles the password reset with OTP verification
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	db := database.DB

	// Check if user exists
	var user model.UserAuth
	err := db.Where("email = ?", req.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found with this email address",
		})
		return
	} else if err != nil {
		fmt.Printf("Database error while checking user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error occurred",
		})
		return
	}

	// Get OTP record
	var passwordOTP model.PasswordOTP
	err = db.Where("email = ?", req.Email).First(&passwordOTP).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No OTP found for this email. Please request a new OTP.",
		})
		return
	} else if err != nil {
		fmt.Printf("Database error while fetching OTP: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error occurred",
		})
		return
	}

	// Check if OTP has expired
	if time.Now().After(passwordOTP.ExpiresAt) {
		// Clean up expired OTP
		db.Delete(&passwordOTP)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP has expired. Please request a new OTP.",
		})
		return
	}

	// Verify OTP
	err = bcrypt.CompareHashAndPassword([]byte(passwordOTP.OtpHash), []byte(req.OTP))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OTP provided",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing new password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process new password",
		})
		return
	}

	// Update user password
	err = db.Model(&user).Update("password", string(hashedPassword)).Error
	if err != nil {
		fmt.Printf("Error updating user password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update password",
		})
		return
	}

	// Clean up used OTP
	db.Delete(&passwordOTP)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
		"email":   req.Email,
	})
}
