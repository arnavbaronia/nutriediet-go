package main

import (
	database "github.com/cd-Ishita/nutriediet-go/database"
	model "github.com/cd-Ishita/nutriediet-go/model"
)

func init() {
	database.ConnectToDB()
}

func main() {
	//database.DB.AutoMigrate(&model.Client{})
	//database.DB.AutoMigrate(&model.DietHistory{})
	//database.DB.AutoMigrate(&model.DietTemplate{})
	//database.DB.AutoMigrate(&model.Recipe{})
	database.DB.AutoMigrate(&model.Exercise{})
	database.DB.AutoMigrate(&model.UserAuth{})
}
