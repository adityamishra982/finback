package records

import "time"

type CreateRecordRequest struct {
	Amount   float64   `json:"amount" binding:"required,gt=0"`
	Type     string    `json:"type" binding:"required,oneof=income expense"`
	Category string    `json:"category" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
	Notes    string    `json:"notes"`
}

type UpdateRecordRequest struct {
	Amount   *float64  `json:"amount" binding:"omitempty,gt=0"`
	Type     *string   `json:"type" binding:"omitempty,oneof=income expense"`
	Category *string   `json:"category" binding:"omitempty"`
	Date     *time.Time `json:"date" binding:"omitempty"`
	Notes    *string   `json:"notes"`
}
