package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Technician struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null;unique"`
	CreatedAt time.Time
}

func (t *Technician) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
