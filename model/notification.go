package model

import "time"

// currently only used for motivations but can be extended
type Notification struct {
	ID            uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Type          string `gorm:"column:type" json:"type,omitempty"`
	Text          string `gorm:"column:text" json:"text,omitempty"`
	PostingActive bool   `gorm:"column:posting_active" json:"posting_active,omitempty"`

	CreatedAt *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`
}

type CreateNotificationReq struct {
	Text          string `json:"text"`
	PostingActive bool   `json:"posting_active,omitempty"`
}

type UpdateNotificationReq struct {
	ID            uint   `json:"id"`
	Text          string `json:"text"`
	PostingActive bool   `json:"posting_active"`
}
