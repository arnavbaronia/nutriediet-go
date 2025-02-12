package model

import "time"

type GetWeightHistoryForClientResponse struct {
	Weight float32   `gorm:"weight" json:"weight"`
	Date   time.Time `gorm:"date" json:"date"`
}
