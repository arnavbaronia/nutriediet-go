package model

import "time"

type GetWeightHistoryForClientResponse struct {
	Weight float32   `gorm:"weight" json:"weight"`
	Date   time.Time `gorm:"date" json:"date"`
}

type WeightUpdateRequest struct {
	Weight   float32 `json:"weight"`
	Feedback string  `json:"feedback"`
}
