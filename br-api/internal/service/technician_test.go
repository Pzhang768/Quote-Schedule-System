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

func TestGetSchedule(t *testing.T) {
	techID := uuid.New()
	date := time.Now()

	tests := []struct {
		name      string
		setupMock func(*mocks.TechnicianStore, *mocks.JobStore)
		wantLen   int
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(ts *mocks.TechnicianStore, js *mocks.JobStore) {
				ts.On("GetByID", techID).Return(&models.Technician{ID: techID}, nil)
				js.On("ListByTechnicianAndDate", techID, date).Return([]models.Job{{ID: uuid.New()}}, nil)
			},
			wantLen: 1,
		},
		{
			name: "technician not found",
			setupMock: func(ts *mocks.TechnicianStore, js *mocks.JobStore) {
				ts.On("GetByID", techID).Return(&models.Technician{}, gorm.ErrRecordNotFound)
			},
			wantErr: ErrTechnicianNotFound,
		},
		{
			name: "store error",
			setupMock: func(ts *mocks.TechnicianStore, js *mocks.JobStore) {
				ts.On("GetByID", techID).Return(&models.Technician{}, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := &mocks.TechnicianStore{}
			js := &mocks.JobStore{}
			tc.setupMock(ts, js)

			slots, err := (&TechnicianService{technicians: ts, jobs: js}).GetSchedule(techID, date)
			if tc.wantErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, slots, tc.wantLen)
		})
	}
}
