package service

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/mocks"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNotificationList(t *testing.T) {
	recipientID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mocks.NotificationStore)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "returns list",
			setupMock: func(n *mocks.NotificationStore) {
				n.On("List", models.RecipientTypeTechnician, recipientID).
					Return([]models.Notification{{Message: "You have a new job"}}, nil)
			},
			wantLen: 1,
		},
		{
			name: "store error",
			setupMock: func(n *mocks.NotificationStore) {
				n.On("List", models.RecipientTypeTechnician, recipientID).
					Return([]models.Notification{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := &mocks.NotificationStore{}
			tc.setupMock(n)

			result, err := (&NotificationService{notifications: n}).List(models.RecipientTypeTechnician, recipientID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, result, tc.wantLen)
		})
	}
}

func TestNotificationRead(t *testing.T) {
	id := uuid.New()
	recipientID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mocks.NotificationStore)
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(n *mocks.NotificationStore) {
				n.On("Read", id, recipientID).Return(nil)
			},
		},
		{
			name: "store error",
			setupMock: func(n *mocks.NotificationStore) {
				n.On("Read", id, recipientID).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := &mocks.NotificationStore{}
			tc.setupMock(n)

			err := (&NotificationService{notifications: n}).Read(id, recipientID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestNotificationRead_NotFound(t *testing.T) {
	id := uuid.New()
	recipientID := uuid.New()

	n := &mocks.NotificationStore{}
	n.On("Read", id, recipientID).Return(gorm.ErrRecordNotFound)

	err := (&NotificationService{notifications: n}).Read(id, recipientID)
	assert.ErrorIs(t, err, ErrNotificationNotFound)
}

func TestNotificationListSince(t *testing.T) {
	recipientID := uuid.New()
	since := time.Now().Add(-time.Minute)

	tests := []struct {
		name      string
		setupMock func(*mocks.NotificationStore)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "returns list since timestamp",
			setupMock: func(n *mocks.NotificationStore) {
				n.On("ListSince", models.RecipientTypeTechnician, recipientID, since).
					Return([]models.Notification{{Message: "new job"}}, nil)
			},
			wantLen: 1,
		},
		{
			name: "store error",
			setupMock: func(n *mocks.NotificationStore) {
				n.On("ListSince", models.RecipientTypeTechnician, recipientID, since).
					Return([]models.Notification{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := &mocks.NotificationStore{}
			tc.setupMock(n)

			result, err := (&NotificationService{notifications: n}).ListSince(models.RecipientTypeTechnician, recipientID, since)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, result, tc.wantLen)
		})
	}
}
