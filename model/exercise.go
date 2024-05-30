package model

// ExerciseType = 1 if descriptive
// ExerciseType = 2 if link
type Exercise struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name        string `gorm:"column:name" json:"name,omitempty"`
	Description string `gorm:"column:description" json:"description,omitempty"`
	Link        string `gorm:"column:link" json:"link,omitempty"`
}
