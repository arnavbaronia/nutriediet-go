package model

import (
	"time"
)

type Recipe struct {
	ID        uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name      string `gorm:"column:name" json:"name,omitempty"`
	ImageData []byte `gorm:"column:image_data;type:longblob" json:"-"`
	ImageType string `gorm:"column:image_type" json:"image_type"`

	CreatedAt *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`
}
type CreateRecipeRequest struct {
	Name        string   `json:"name,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	Preparation []string `json:"preparation,omitempty"`
}

type UpdateRecipeRequest struct {
	ID          uint     `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	Preparation []string `json:"preparation,omitempty"`
}

type GetRecipeResponse struct {
	ID          uint
	Name        string
	Ingredients []string
	Preparation []string
}

type GetListOfRecipesResponse struct {
	Name     string
	RecipeID uint
}
