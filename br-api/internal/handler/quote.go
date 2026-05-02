package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/service"
)

type QuoteHandler struct {
	svc *service.QuoteService
}

func NewQuoteHandler(svc *service.QuoteService) *QuoteHandler {
	return &QuoteHandler{svc: svc}
}

func (h *QuoteHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	quotes, err := h.svc.ListUnscheduled(page, pageSize)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, http.StatusOK, quotes)
}

func (h *QuoteHandler) Create(c *gin.Context) {
	var body models.Quote
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := h.svc.Create(&body); err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, http.StatusCreated, body)
}
