package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecipientType string

const (
	RecipientTypeTechnician RecipientType = "technician"
	RecipientTypeManager    RecipientType = "manager"
)

type Notification struct {
	ID            uuid.UUID        `gorm:"type:char(36);primaryKey"`
	Type          NotificationType `gorm:"not null"`
	RecipientType RecipientType    `gorm:"not null;index:idx_notifications_recipient"`
	RecipientID   uuid.UUID        `gorm:"type:char(36);not null;index:idx_notifications_recipient"`
	JobID         uuid.UUID        `gorm:"type:char(36);not null"`
	Job           Job
	Message       string    `gorm:"not null"`
	
	ReadAt        *time.Time
	CreatedAt     time.Time `gorm:"index:idx_notifications_recipient"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.ID = uuid.New()
	return nil
}
