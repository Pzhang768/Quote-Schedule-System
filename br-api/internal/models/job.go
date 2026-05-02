package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey"`
	TechnicianID uuid.UUID  `gorm:"type:char(36);not null;index:idx_jobs_conflict"`
	Technician   Technician
	QuoteID      uuid.UUID  `gorm:"type:char(36);not null;unique"`
	Quote        Quote
	ManagerID    uuid.UUID  `gorm:"type:char(36);not null;index"`
	Manager      Manager
	
	StartsAt     time.Time  `gorm:"not null;index:idx_jobs_conflict"`
	EndsAt       time.Time  `gorm:"not null;index:idx_jobs_conflict"`
	Status       JobStatus  `gorm:"not null;index:idx_jobs_conflict"`
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	j.ID = uuid.New()
	return nil
}
