package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Manager struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null;unique"`
	CreatedAt time.Time
}

func (m *Manager) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}