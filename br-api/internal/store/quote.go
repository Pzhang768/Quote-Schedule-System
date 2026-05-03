package store

import (
	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuoteStore struct {
	db *gorm.DB
}

func NewQuoteStore(db *gorm.DB) *QuoteStore {
	return &QuoteStore{db: db}
}

func (s *QuoteStore) Create(q *models.Quote) error {
	result := s.db.Create(q)
	return result.Error
}

func (s *QuoteStore) GetByID(id uuid.UUID) (*models.Quote, error) {
	var q models.Quote
	result := s.db.First(&q, "id = ?", id)
	return &q, result.Error
}

func (s *QuoteStore) GetByIDForUpdate(tx *gorm.DB, id uuid.UUID) (*models.Quote, error) {
	var q models.Quote
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&q, "id = ?", id)
	return &q, result.Error
}

func (s *QuoteStore) List(status models.QuoteStatus, page, pageSize int) ([]models.Quote, error) {
	var quotes []models.Quote
	result := s.db.Where("status = ?", status).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&quotes)
	return quotes, result.Error
}

func (s *QuoteStore) Count(status models.QuoteStatus) (int, error) {
	var count int64
	result := s.db.Model(&models.Quote{}).Where("status = ?", status).Count(&count)
	return int(count), result.Error
}

func (s *QuoteStore) UpdateStatus(tx *gorm.DB, id uuid.UUID, status models.QuoteStatus) error {
	result := tx.Model(&models.Quote{}).Where("id = ?", id).Update("status", status)
	return result.Error
}
