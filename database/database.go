package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := "host=rain.db.elephantsql.com user=aaoylwvk password=vPrJ8niX2QnQprwhQW7NJwQS4e6CJDcz dbname=aaoylwvk port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	DB = db
}
