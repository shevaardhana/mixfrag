package models

import "time"

type EssentialOil struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	CategorySmellId int        `json:"category_smell_id"`
	CategoryNoteId  int        `json:"category_note_id"`
	RankNote        string     `json:"rank_note"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       string     `json:"created_by"`
	ModifiedAt      *time.Time `json:"modified_at"`
	ModifiedBy      *string    `json:"modified_by"`
}
