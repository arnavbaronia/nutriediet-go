package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// GenerateOTP generates a 6-digit numeric OTP
func GenerateOTP() (string, error) {
	// Generate a random number between 100000 and 999999
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	
	otp := n.Int64() + 100000
	return fmt.Sprintf("%06d", otp), nil
}

// SendOTPEmail sends OTP via SMTP
func SendOTPEmail(toEmail, otp string) error {
	// Get SMTP configuration from environment variables
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	
	// Default values if not set
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Default to Gmail
	}
	if smtpPortStr == "" {
		smtpPortStr = "587" // Default SMTP port
	}
	
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}
	
	if smtpEmail == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP credentials not configured. Please set SMTP_EMAIL and SMTP_PASSWORD environment variables")
	}
	
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your OTP for Password Reset")
	
	body := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Your OTP for password reset is: <strong>%s</strong></p>
		<p>This OTP will expire in 5 minutes.</p>
		<p>If you didn't request this password reset, please ignore this email.</p>
	`, otp)
	
	m.SetBody("text/html", body)
	
	// Send the email
	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)
	
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	
	return nil
} 