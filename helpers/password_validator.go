package helpers

import (
	"errors"
	"regexp"
	"strings"
)

// ValidatePasswordStrength validates password meets security requirements
// Requirements:
// - Minimum 12 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one number
// - At least one special character
func ValidatePasswordStrength(password string) error {
	// Check minimum length
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}
	
	// Check maximum length (prevent DoS attacks with very long passwords)
	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}
	
	// Check for uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	
	// Check for lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	
	// Check for digit
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	
	// Check for special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>_\-+=\[\]\\\/;'~` + "`" + `]`).MatchString(password)
	if !hasSpecial {
		return errors.New("password must contain at least one special character (!@#$%^&*(),.?\":{}|<>_-+=[]\\\\\\//;'~`)")
	}
	
	// Check for common weak passwords
	weakPasswords := []string{
		"password123!", "admin123456!", "welcome12345!",
		"qwerty123456!", "123456789abc!", "letmein12345!",
	}
	
	passwordLower := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if strings.Contains(passwordLower, strings.ToLower(weak)) {
			return errors.New("password is too common, please choose a stronger password")
		}
	}
	
	return nil
}

// GetPasswordRequirements returns a user-friendly message about password requirements
func GetPasswordRequirements() string {
	return "Password must be at least 12 characters long and contain: uppercase letter, lowercase letter, number, and special character (!@#$%^&*(),.?\":{}|<>_-+=[]\\\\\\//;'~`)"
}

