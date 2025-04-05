package model

import "time"

type Payment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	ClientID    uint      `gorm:"column:client_id" json:"client_id,omitempty"`
	Date        time.Time `gorm:"column:date" json:"date,omitempty"`
	AmountPaid  int64     `gorm:"column:amount_paid" json:"amount_paid,omitempty"`
	AmountDue   int64     `gorm:"column:amount_due" json:"amount_due,omitempty"`
	TotalAmount int64     `gorm:"column:total_amount" json:"total_amount,omitempty"`
	Package     string    `gorm:"column:package" json:"package,omitempty"`

	CreatedAt *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`
}
