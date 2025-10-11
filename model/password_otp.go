package model

import "time"

type PasswordOTP struct {
	Email        string     `gorm:"primaryKey;type:varchar(255)" json:"email"`
	OtpHash      string     `gorm:"not null;type:text" json:"otp_hash"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	Attempts     int        `gorm:"not null;default:0" json:"attempts"`
	MaxAttempts  int        `gorm:"not null;default:5" json:"max_attempts"`
	LockedUntil  *time.Time `gorm:"null" json:"locked_until,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
}

func (PasswordOTP) TableName() string {
	return "password_otps"
}
