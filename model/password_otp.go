package model

import "time"

type PasswordOTP struct {
	Email     string    `gorm:"primaryKey;type:varchar(255)" json:"email"`
	OtpHash   string    `gorm:"not null;type:text" json:"otp_hash"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (PasswordOTP) TableName() string {
	return "password_otps"
}
