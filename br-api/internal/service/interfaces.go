package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
)

type jobStore interface {
	Create(tx *gorm.DB, j *models.Job) error
	GetByID(id uuid.UUID) (*models.Job, error)
	ListByTechnicianAndDate(technicianID uuid.UUID, date time.Time) ([]models.Job, error)
	UpdateStatus(id uuid.UUID, status models.JobStatus) error
	ConflictCheck(tx *gorm.DB, technicianID uuid.UUID, startsAt, endsAt time.Time) ([]models.Job, error)
}

type quoteStore interface {
	Create(q *models.Quote) error
	List(status models.QuoteStatus, page, pageSize int) ([]models.Quote, error)
	Count(status models.QuoteStatus) (int, error)
	GetByID(id uuid.UUID) (*models.Quote, error)
	UpdateStatus(tx *gorm.DB, id uuid.UUID, status models.QuoteStatus) error
}

type notificationStore interface {
	Create(n *models.Notification) error
	List(recipientType models.RecipientType, recipientID uuid.UUID) ([]models.Notification, error)
	ListSince(recipientType models.RecipientType, recipientID uuid.UUID, since time.Time) ([]models.Notification, error)
	Read(id, recipientID uuid.UUID) error
}

type technicianStore interface {
	List(page, pageSize int) ([]models.Technician, error)
	Count() (int, error)
	GetByID(id uuid.UUID) (*models.Technician, error)
}

type managerStore interface {
	List(page, pageSize int) ([]models.Manager, error)
	Count() (int, error)
}
