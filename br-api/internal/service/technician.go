package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
	"gorm.io/gorm"
)

type TechnicianService struct {
	technicians technicianStore
	jobs        jobStore
}

func NewTechnicianService(technicians *store.TechnicianStore, jobs *store.JobStore) *TechnicianService {
	return &TechnicianService{technicians: technicians, jobs: jobs}
}

func (s *TechnicianService) List(page, pageSize int) ([]models.Technician, error) {
	return s.technicians.List(page, pageSize)
}

func (s *TechnicianService) GetSchedule(technicianID uuid.UUID, date time.Time) ([]JobSlotResponse, error) {
	_, err := s.technicians.GetByID(technicianID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTechnicianNotFound
		}
		return nil, err
	}

	jobs, err := s.jobs.ListByTechnicianAndDate(technicianID, date)
	if err != nil {
		return nil, err
	}

	slots := make([]JobSlotResponse, len(jobs))
	for i, j := range jobs {
		slots[i] = JobSlotResponse{ID: j.ID, StartsAt: j.StartsAt, EndsAt: j.EndsAt, Status: j.Status}
	}
	return slots, nil
}
