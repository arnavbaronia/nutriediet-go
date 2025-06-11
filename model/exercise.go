package model

import "time"

// ExerciseType = 1 if descriptive
// ExerciseType = 2 if link
type Exercise struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"column:name" json:"name,omitempty"`
	Description string `gorm:"column:description" json:"description,omitempty"`
	Link        string `gorm:"column:link" json:"link,omitempty"`

	// POSTGRES fields
	//CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	//UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	//DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp;default:NULL;" json:"deleted_at,omitempty"`

	CreatedAt *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`
}

type GetListOfExercisesResponse struct {
	Name string
	ID   uint
}

type FavoriteExercise struct {
	ClientID   string     `gorm:"primaryKey"`
	ExerciseID uint       `gorm:"primaryKey"`
	CreatedAt  *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;"`
}
