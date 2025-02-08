package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectToDB() {
	// ssh verification is skipped for now, remove before deployment
	dsn := "avnadmin:AVNS_7QDxgZDlRhQXAx3QV4z@tcp(nutriediet-mysql-ishitagupta-5564.f.aivencloud.com:22013)/defaultdb?tls=skip-verify&parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect to database")
	}

	DB = db
}
