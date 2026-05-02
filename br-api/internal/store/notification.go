package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
)

type NotificationStore struct {
	db *gorm.DB
}

func NewNotificationStore(db *gorm.DB) *NotificationStore {
	return &NotificationStore{db: db}
}

func (s *NotificationStore) Create(n *models.Notification) error {
	result := s.db.Create(n)
	return result.Error
}

func (s *NotificationStore) Read(id, recipientID uuid.UUID) error {
	now := time.Now()
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND recipient_id = ?", id, recipientID).
		Update("read_at", now)
	return result.Error
}

func (s *NotificationStore) List(recipientType models.RecipientType, recipientID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	result := s.db.Where("recipient_type = ? AND recipient_id = ?", recipientType, recipientID).
		Order("created_at desc").Find(&notifications)
	return notifications, result.Error
}

func (s *NotificationStore) ListSince(recipientType models.RecipientType, recipientID uuid.UUID, since time.Time) ([]models.Notification, error) {
	var notifications []models.Notification
	result := s.db.Where("recipient_type = ? AND recipient_id = ? AND created_at > ?", recipientType, recipientID, since).
		Order("created_at asc").Find(&notifications)
	return notifications, result.Error
}
