package model

import "time"

type Client struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name              string    `gorm:"column:name" json:"name,omitempty"`
	Age               int64     `gorm:"column:age" json:"age,omitempty"`
	City              string    `gorm:"column:city" json:"city,omitempty"`
	PhoneNumber       string    `gorm:"column:phone_number" json:"phone_number,omitempty"`
	DateOfJoining     time.Time `gorm:"column:date_of_joining" json:"date_of_joining,omitempty"`
	Package           string    `gorm:"column:package" json:"package,omitempty"`
	AmountPaid        int64     `gorm:"column:amount_paid" json:"amount_paid,omitempty"`
	LastPaymentDate   time.Time `gorm:"column:last_payment_date" json:"last_payment_date,omitempty"`
	NextPaymentDate   time.Time `gorm:"column:next_payment_date" json:"next_payment_date,omitempty"`
	Remarks           string    `gorm:"column:remarks" json:"remarks,omitempty"`
	DietitianId       int       `gorm:"column:dietitian_id" json:"dietitian_id,omitempty"`
	Group             int       `gorm:"column:\"group\"" json:"group,omitempty"`
	Email             string    `gorm:"column:email" json:"email,omitempty"`
	Height            int       `gorm:"column:height" json:"height,omitempty"`
	StartingWeight    float32   `gorm:"column:starting_weight" json:"starting_weight,omitempty"`
	DietaryPreference string    `gorm:"column:dietary_preference" json:"dietary_preference,omitempty"`
	MedicalHistory    string    `gorm:"column:medical_history" json:"medical_history,omitempty"`
	Allergies         string    `gorm:"column:allergies" json:"allergies,omitempty"`
	Stay              string    `gorm:"column:stay" json:"stay,omitempty"`
	Exercise          string    `gorm:"column:exercise" json:"exercise,omitempty"`
	Comments          string    `gorm:"column:comments" json:"comments,omitempty"`
	DietRecall        string    `gorm:"column:diet_recall" json:"diet_recall,omitempty"`
	IsActive          bool      `gorm:"column:is_active" json:"is_active,omitempty"`
	Locality          string    `gorm:"column:locality" json:"locality,omitempty"`
	// USE WHEN MOVE TO MYSQL TABLES
	//CreatedAt         *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	//UpdatedAt         *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	//DeletedAt         *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`

	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp;default:NULL;" json:"deleted_at,omitempty"`
}

type ClientMiniInfo struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name            string    `gorm:"column:name" json:"name,omitempty"`
	DietitianId     int       `gorm:"column:dietitian_id" json:"dietitian_id,omitempty"`
	Group           int       `gorm:"column:\"group\"" json:"group,omitempty"`
	Email           string    `gorm:"column:email" json:"email,omitempty"`
	IsActive        bool      `gorm:"column:is_active" json:"is_active,omitempty"`
	NextPaymentDate time.Time `gorm:"column:next_payment_date" json:"next_payment_date,omitempty"`
	LastDietDate    time.Time `json:"last_diet_date,omitempty"`
}
