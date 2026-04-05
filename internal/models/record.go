package models

import (
	"time"

	"gorm.io/gorm"
)

type RecordType string

const (
	TypeIncome  RecordType = "income"
	TypeExpense RecordType = "expense"
)

type Record struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Amount    float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Type      RecordType     `gorm:"type:varchar(10);not null" json:"type"`
	Category  string         `gorm:"type:varchar(50);not null;index" json:"category"`
	Date      time.Time      `gorm:"not null;index" json:"date"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
