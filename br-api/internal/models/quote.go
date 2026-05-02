package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quote struct {
	ID           uuid.UUID   `gorm:"type:char(36);primaryKey"`
	CustomerName string      `gorm:"not null"`
	Address      string      `gorm:"not null"`
	Description  string
	Status       QuoteStatus `gorm:"not null;default:unscheduled"`
	
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (q *Quote) BeforeCreate(tx *gorm.DB) error {
	q.ID = uuid.New()
	return nil
}
