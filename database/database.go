package database

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDB() {
	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	
	// Validate required environment variables
	if dbUser == "" || dbPassword == "" {
		log.Fatal("❌ Database credentials not configured. Please set DB_USER and DB_PASSWORD environment variables")
	}
	
	// Set defaults for optional variables
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "nutriediet"
	}
	
	// Determine TLS setting based on environment
	environment := os.Getenv("ENVIRONMENT")
	tlsConfig := "false"
	
	// For production on Digital Ocean or cloud databases, enable TLS
	if environment == "production" {
		// For local MySQL (same machine), TLS is not needed
		if dbHost == "localhost" || dbHost == "127.0.0.1" {
			tlsConfig = "false"
		} else {
			// For remote databases (like Aiven), use TLS
			tlsConfig = "skip-verify" // Change to "true" when you have proper certificates
		}
	}
	
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&tls=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, tlsConfig)
	
	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Only log errors in production
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}
	
	// Development mode - show SQL queries
	if environment == "development" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	
	// Connect to database
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	
	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}
	
	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)                  // Maximum idle connections
	sqlDB.SetMaxOpenConns(100)                 // Maximum open connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Connection lifetime
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Idle connection lifetime
	
	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}
	
	DB = db
	log.Printf("✅ Database connected successfully (Host: %s, Database: %s)", dbHost, dbName)
}
