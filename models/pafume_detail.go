package models

import "time"

type ParfumeDetail struct {
	ID          int        `json:"id"`
	ParfumeId   int        `json:"parfume_id"`
	OilId       string     `json:"oil_id"`
	TotalDrop   string     `json:"total_drop"`
	Rasio       float64    `json:"rasio"`
	WeightTotal int        `json:"weight_total"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	ModifiedAt  *time.Time `json:"modified_at"`
	ModifiedBy  *string    `json:"modified_by"`
}
