package main

import (
	database "github.com/cd-Ishita/nutriediet-go/database"
)

func init() {
	database.ConnectToDB()
}

func main() {
	//database.DB.AutoMigrate(&model.Client{})
	//database.DB.AutoMigrate(&model.DietHistory{})
	//database.DB.AutoMigrate(&model.DietTemplate{})
	//database.DB.AutoMigrate(&model.Client{})
	//database.DB.AutoMigrate(&model.Exercise{})
	//database.DB.AutoMigrate(&model.UserAuth{})
	//database.DB.AutoMigrate(&model.Notification{})
	//database.DB.AutoMigrate(&model.Payment{})
	//dummyData1()
}

//func dummyData1() {
//	userAuth := []model.UserAuth{
//		{
//			FirstName: "Swapna",
//			LastName: "Boyapalli",
//
//		},
//		{},
//		{},
//		{},
//	}
//}

//func dummyData() {
//	loc, _ := time.LoadLocation("Asia/Kolkata")
//	client1 := model.Client{
//		Name:              "Yedla Pranavi Reddy",
//		Age:               17,
//		City:              "Hyderabad",
//		PhoneNumber:       "+918897315213",
//		DateOfJoining:     time.Date(2024, 7, 9, 0, 0, 0, 0, loc),
//		Package:           "1_MONTH",
//		AmountPaid:        3000,
//		LastPaymentDate:   time.Date(2024, 6, 11, 0, 0, 0, 0, loc),
//		NextPaymentDate:   time.Date(2024, 7, 9, 0, 0, 0, 0, loc),
//		Remarks:           "Wrote Neet exam so no physical activity",
//		DietitianId:       2,
//		Group:             5,
//		Email:             "yedlapranavireddy222@gmail.com",
//		Height:            165,
//		StartingWeight:    110,
//		DietaryPreference: "Non Vegetarian",
//		MedicalHistory:    "No",
//		Allergies:         "No",
//		Stay:              "Home",
//		Exercise:          "No",
//		Comments:          "",
//		DietRecall:        "South Indian Food",
//		IsActive:          true,
//		Locality:          "Miyapur",
//	}
//
//	db := database.DB
//	err := db.Table("clients").Create(&client1).Error
//	if err != nil {
//		fmt.Errorf("error: ", err)
//	}
//}
