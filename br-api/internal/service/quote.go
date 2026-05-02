package service

import (
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
)

type QuoteService struct {
	quotes quoteStore
}

func NewQuoteService(quotes *store.QuoteStore) *QuoteService {
	return &QuoteService{quotes: quotes}
}

func (s *QuoteService) ListUnscheduled(page, pageSize int) ([]models.Quote, error) {
	return s.quotes.List(models.QuoteStatusUnscheduled, page, pageSize)
}
