package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
	"gorm.io/gorm"
)

type NotificationService struct {
	notifications notificationStore
}

func NewNotificationService(notifications *store.NotificationStore) *NotificationService {
	return &NotificationService{notifications: notifications}
}

func (s *NotificationService) List(recipientType models.RecipientType, recipientID uuid.UUID) ([]NotificationResponse, error) {
	notifications, err := s.notifications.List(recipientType, recipientID)
	if err != nil {
		return nil, err
	}
	return ToNotificationResponses(notifications), nil
}

func (s *NotificationService) Read(id, recipientID uuid.UUID) error {
	err := s.notifications.Read(id, recipientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}
	return nil
}
