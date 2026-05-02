package store

import (
	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
)

type TechnicianStore struct {
	db *gorm.DB
}

func NewTechnicianStore(db *gorm.DB) *TechnicianStore {
	return &TechnicianStore{db: db}
}

func (s *TechnicianStore) List(page, pageSize int) ([]models.Technician, error) {
	var technicians []models.Technician
	result := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&technicians)
	return technicians, result.Error
}

func (s *TechnicianStore) Count() (int, error) {
	var count int64
	result := s.db.Model(&models.Technician{}).Count(&count)
	return int(count), result.Error
}

func (s *TechnicianStore) GetByID(id uuid.UUID) (*models.Technician, error) {
	var t models.Technician
	result := s.db.First(&t, "id = ?", id)
	return &t, result.Error
}

