package store

import (
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
)

type ManagerStore struct {
	db *gorm.DB
}

func NewManagerStore(db *gorm.DB) *ManagerStore {
	return &ManagerStore{db: db}
}

func (s *ManagerStore) Count() (int, error) {
	var count int64
	result := s.db.Model(&models.Manager{}).Count(&count)
	return int(count), result.Error
}

func (s *ManagerStore) List(page, pageSize int) ([]models.Manager, error) {
	var managers []models.Manager
	result := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&managers)
	return managers, result.Error
}
