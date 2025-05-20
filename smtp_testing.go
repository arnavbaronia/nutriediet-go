package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func main() {
	// Load environment variables manually for testing
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_FROM_EMAIL", "nutriedietplan@gmail.com")
	os.Setenv("SMTP_PASSWORD", "dddd djzr jtfn rowv") // Your actual app password

	from := os.Getenv("SMTP_FROM_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	to := "arnav.baronia@gmail.com" // Change to your real test email

	msg := []byte(
		"From: " + from + "\n" +
			"To: " + to + "\n" +
			"Subject: SMTP Test\n\n" +
			"This is a test email from NutriEdiet SMTP setup.",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from,
		[]string{to},
		msg,
	)

	if err != nil {
		fmt.Println("SMTP error:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
