package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/service"
)

type TechnicianHandler struct {
	svc *service.TechnicianService
}

func NewTechnicianHandler(svc *service.TechnicianService) *TechnicianHandler {
	return &TechnicianHandler{svc: svc}
}

// @Summary     List technicians
// @Tags        technicians
// @Produce     json
// @Param       page      query  int  false  "Page number"     default(1)
// @Param       page_size query  int  false  "Items per page"  default(20)
// @Success     200  {object}  Response[[]models.Technician]
// @Router      /technicians [get]
func (h *TechnicianHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	result, err := h.svc.List(page, pageSize)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessPaged(c, http.StatusOK, result.Items, result.Total, page, pageSize)
}

// @Summary     Get a technician's job schedule
// @Tags        technicians
// @Produce     json
// @Param       id    path   string  true   "Technician ID"
// @Param       date  query  string  false  "Date (YYYY-MM-DD), defaults to today"
// @Success     200  {object}  Response[[]service.JobSlotResponse]
// @Failure     400  {object}  ErrorResponse
// @Router      /technicians/{id}/jobs [get]
func (h *TechnicianHandler) GetSchedule(c *gin.Context) {
	technicianID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid technician id")
		return
	}

	date := time.Now().UTC()
	if d := c.Query("date"); d != "" {
		parsed, parseErr := time.ParseInLocation("2006-01-02", d, time.UTC)
		if parseErr != nil {
			Fail(c, http.StatusBadRequest, "date must be YYYY-MM-DD")
			return
		}
		date = parsed
	}

	slots, err := h.svc.GetSchedule(technicianID, date)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, http.StatusOK, slots)
}

