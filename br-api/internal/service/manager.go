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

type PagedManagers struct {
	Items []models.Manager
	Total int
}

func (s *ManagerService) List(page, pageSize int) (PagedManagers, error) {
	items, err := s.managers.List(page, pageSize)
	if err != nil {
		return PagedManagers{}, err
	}
	total, err := s.managers.Count()
	if err != nil {
		return PagedManagers{}, err
	}
	return PagedManagers{Items: items, Total: total}, nil
}
