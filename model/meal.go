package model

type MealAdditionalInfo struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name" json:"name,omitempty"`
	Type string `gorm:"column:type" json:"type" validate:"required, eq=QUANTITY | eq=MEAL"`
}

type CreateNewMealRequest struct {
	Name        string   `json:"name"`
	HasRecipe   bool     `json:"has_recipe"`
	Ingredients []string `json:"ingredients"`
	Preparation []string `json:"preparation"`
}
