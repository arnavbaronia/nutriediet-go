package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
	"time"
)

type EmailService interface {
	SendPasswordResetEmail(email, otp string) error
}

type SMTPService struct{}

func (s *SMTPService) SendPasswordResetEmail(to, otp string) error {
	// Email configuration
	fmt.Println("Attempting to send email to:", to)
	fmt.Println("Using SMTP server:", os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"))
	from := os.Getenv("SMTP_FROM_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Template data
	data := struct {
		OTP  string
		Year int
	}{
		OTP:  otp,
		Year: time.Now().Year(),
	}

	// Parse template
	t, err := template.ParseFiles("templates/password_reset.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Execute template
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Email headers
	headers := make(map[string]string)
	headers["From"] = os.Getenv("SMTP_FROM_NAME") + " <" + from + ">"
	headers["To"] = to
	headers["Subject"] = "Password Reset Request"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Send email
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		fmt.Println("SMTP Error Details:", err)
		return fmt.Errorf("SMTP error: %w", err)
	}
	return nil
}

func NewEmailService() EmailService {
	return &SMTPService{}
}
