package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
)

// Input Dto

type AssignJobInput struct {
	QuoteID      uuid.UUID
	TechnicianID uuid.UUID
	ManagerID    uuid.UUID
	StartsAt     time.Time
}

// Response Dtos

type JobResponse struct {
	ID           uuid.UUID        `json:"id"`
	QuoteID      uuid.UUID        `json:"quote_id"`
	TechnicianID uuid.UUID        `json:"technician_id"`
	ManagerID    uuid.UUID        `json:"manager_id"`
	StartsAt     time.Time        `json:"starts_at"`
	EndsAt       time.Time        `json:"ends_at"`
	Status       models.JobStatus `json:"status"`
	CompletedAt  *time.Time       `json:"completed_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}

type JobSlotResponse struct {
	ID       uuid.UUID        `json:"id"`
	StartsAt time.Time        `json:"starts_at"`
	EndsAt   time.Time        `json:"ends_at"`
	Status   models.JobStatus `json:"status"`
}

type NotificationResponse struct {
	ID            uuid.UUID               `json:"id"`
	Type          models.NotificationType `json:"type"`
	RecipientType models.RecipientType    `json:"recipient_type"`
	RecipientID   uuid.UUID               `json:"recipient_id"`
	JobID         uuid.UUID               `json:"job_id"`
	Message       string                  `json:"message"`
	ReadAt        *time.Time              `json:"read_at,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
}

// Mapping helpers, convert models to response DTOs, remove relation fields

func ToJobResponse(j *models.Job) JobResponse {
	return JobResponse{
		ID:           j.ID,
		QuoteID:      j.QuoteID,
		TechnicianID: j.TechnicianID,
		ManagerID:    j.ManagerID,
		StartsAt:     j.StartsAt,
		EndsAt:       j.EndsAt,
		Status:       j.Status,
		CompletedAt:  j.CompletedAt,
		CreatedAt:    j.CreatedAt,
	}
}

func ToNotificationResponses(notifications []models.Notification) []NotificationResponse {
	result := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		result[i] = NotificationResponse{
			ID:            n.ID,
			Type:          n.Type,
			RecipientType: n.RecipientType,
			RecipientID:   n.RecipientID,
			JobID:         n.JobID,
			Message:       n.Message,
			ReadAt:        n.ReadAt,
			CreatedAt:     n.CreatedAt,
		}
	}
	return result
}
