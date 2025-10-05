package main

import (
	"log"
	
	"github.com/cd-Ishita/nutriediet-go/database"
	"github.com/cd-Ishita/nutriediet-go/model"
	"github.com/joho/godotenv"
)

// Run this to add the new columns to password_otps table
// go run migrate/update_password_otps.go

func main() {
	// Load environment variables first
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	} else {
		log.Println("✅ Environment variables loaded")
	}
	
	log.Println("🔧 Starting password_otps table migration...")
	
	// Connect to database
	database.ConnectToDB()
	
	// Auto-migrate the PasswordOTP model
	// This will add the new columns: attempts, max_attempts, locked_until
	migrationErr := database.DB.AutoMigrate(&model.PasswordOTP{})
	if migrationErr != nil {
		log.Fatalf("❌ Migration failed: %v", migrationErr)
	}
	
	log.Println("✅ Migration completed successfully!")
	log.Println("   - Added 'attempts' column (default: 0)")
	log.Println("   - Added 'max_attempts' column (default: 5)")
	log.Println("   - Added 'locked_until' column (nullable)")
}

