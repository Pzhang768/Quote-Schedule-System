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

type PagedQuotes struct {
	Items    []models.Quote
	Total    int
}

func (s *QuoteService) ListUnscheduled(page, pageSize int) (PagedQuotes, error) {
	items, err := s.quotes.List(models.QuoteStatusUnscheduled, page, pageSize)
	if err != nil {
		return PagedQuotes{}, err
	}
	total, err := s.quotes.Count(models.QuoteStatusUnscheduled)
	if err != nil {
		return PagedQuotes{}, err
	}
	return PagedQuotes{Items: items, Total: total}, nil
}

func (s *QuoteService) Create(q *models.Quote) error {
	return s.quotes.Create(q)
}
