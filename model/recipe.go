package model

type Recipe struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"column:name" json:"name,omitempty"`
	MealID      int    `gorm:"column:food_id" json:"food_id,omitempty"`
	Ingredients string `gorm:"column:ingredients;type:text" json:"ingredients,omitempty"`
	Preparation string `gorm:"column:preparation;type:text" json:"preparation,omitempty"`
}
