package model

import "time"

type UserAuth struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement:true" json:"id"`
	FirstName     string    `gorm:"first_name" json:"first_name" validate:"required, min=2, max=100"`
	LastName      string    `gorm:"last_name" json:"last_name" validate:"required, min=2, max=100"`
	Password      string    `gorm:"password" json:"password" validate:"required, min=6"`
	Email         string    `gorm:"email" json:"email" validate:"email, required"`
	Token         string    `gorm:"token" json:"token"`
	UserType      string    `gorm:"user_type" json:"user_type" validate:"required, eq=ADMIN | eq=CLIENT"`
	RefreshToken  string    `gorm:"refresh_token" json:"refresh_token"`
	ResetToken    string    `gorm:"reset_token" json:"reset_token"`
	ResetTokenExp time.Time `gorm:"reset_token_exp" json:"reset_token_exp"`

	// USE WHEN MOVE TO MYSQL TABLES
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`

	//CreatedAt *time.Time `gorm:"column:created_at;type:timestamp not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	//UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	//DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp;default:NULL;" json:"deleted_at,omitempty"`
}
