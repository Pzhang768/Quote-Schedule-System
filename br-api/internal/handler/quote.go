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

// @Summary     List unscheduled quotes
// @Tags        quotes
// @Produce     json
// @Param       page      query  int  false  "Page number"       default(1)
// @Param       page_size query  int  false  "Items per page"    default(20)
// @Success     200  {object}  Response[[]models.Quote]
// @Router      /quotes [get]
func (h *QuoteHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	result, err := h.svc.ListUnscheduled(page, pageSize)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessPaged(c, http.StatusOK, result.Items, result.Total, page, pageSize)
}

// @Summary     Create a quote
// @Tags        quotes
// @Accept      json
// @Produce     json
// @Param       body  body      models.Quote  true  "Quote"
// @Success     201   {object}  Response[models.Quote]
// @Failure     422   {object}  ErrorResponse
// @Router      /quotes [post]
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
