package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type JobStore struct{ mock.Mock }

func (m *JobStore) Create(tx *gorm.DB, j *models.Job) error {
	return m.Called(tx, j).Error(0)
}
func (m *JobStore) GetByID(id uuid.UUID) (*models.Job, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Job), args.Error(1)
}
func (m *JobStore) ListByTechnicianAndDate(technicianID uuid.UUID, date time.Time) ([]models.Job, error) {
	args := m.Called(technicianID, date)
	return args.Get(0).([]models.Job), args.Error(1)
}
func (m *JobStore) UpdateStatus(id uuid.UUID, status models.JobStatus) error {
	return m.Called(id, status).Error(0)
}
func (m *JobStore) ConflictCheck(tx *gorm.DB, technicianID uuid.UUID, startsAt, endsAt time.Time) ([]models.Job, error) {
	args := m.Called(tx, technicianID, startsAt, endsAt)
	return args.Get(0).([]models.Job), args.Error(1)
}

type QuoteStore struct{ mock.Mock }

func (m *QuoteStore) Create(q *models.Quote) error {
	return m.Called(q).Error(0)
}
func (m *QuoteStore) GetByID(id uuid.UUID) (*models.Quote, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Quote), args.Error(1)
}
func (m *QuoteStore) List(status models.QuoteStatus, page, pageSize int) ([]models.Quote, error) {
	args := m.Called(status, page, pageSize)
	return args.Get(0).([]models.Quote), args.Error(1)
}
func (m *QuoteStore) UpdateStatus(tx *gorm.DB, id uuid.UUID, status models.QuoteStatus) error {
	return m.Called(tx, id, status).Error(0)
}

type NotificationStore struct{ mock.Mock }

func (m *NotificationStore) Create(n *models.Notification) error {
	return m.Called(n).Error(0)
}
func (m *NotificationStore) List(recipientType models.RecipientType, recipientID uuid.UUID) ([]models.Notification, error) {
	args := m.Called(recipientType, recipientID)
	return args.Get(0).([]models.Notification), args.Error(1)
}
func (m *NotificationStore) ListSince(recipientType models.RecipientType, recipientID uuid.UUID, since time.Time) ([]models.Notification, error) {
	args := m.Called(recipientType, recipientID, since)
	return args.Get(0).([]models.Notification), args.Error(1)
}
func (m *NotificationStore) Read(id, recipientID uuid.UUID) error {
	return m.Called(id, recipientID).Error(0)
}

type ManagerStore struct{ mock.Mock }

func (m *ManagerStore) List(page, pageSize int) ([]models.Manager, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]models.Manager), args.Error(1)
}

type TechnicianStore struct{ mock.Mock }

func (m *TechnicianStore) List(page, pageSize int) ([]models.Technician, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]models.Technician), args.Error(1)
}
func (m *TechnicianStore) GetByID(id uuid.UUID) (*models.Technician, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Technician), args.Error(1)
}
