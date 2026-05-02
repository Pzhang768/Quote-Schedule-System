package service

import (
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
)

type ManagerService struct {
	managers managerStore
}

func NewManagerService(managers *store.ManagerStore) *ManagerService {
	return &ManagerService{managers: managers}
}

func (s *ManagerService) List(page, pageSize int) ([]models.Manager, error) {
	return s.managers.List(page, pageSize)
}
