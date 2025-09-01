package models

import "time"

type Parfume struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	TotalMl           string     `json:"total_ml"`
	CategoryParfumeId int        `json:"category_parfume_id"`
	TotalOilDrop      int        `json:"total_oil_drop"`
	TotalOil          int        `json:"total_oil"`
	TotalParfumeBase  int        `json:"total_parfume_base"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         string     `json:"created_by"`
	ModifiedAt        *time.Time `json:"modified_at"`
	ModifiedBy        *string    `json:"modified_by"`
}
