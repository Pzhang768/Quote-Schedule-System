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

func (h *TechnicianHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	technicians, err := h.svc.List(page, pageSize)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, http.StatusOK, technicians)
}

func (h *TechnicianHandler) GetSchedule(c *gin.Context) {
	technicianID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid technician id")
		return
	}

	date := time.Now()
	if d := c.Query("date"); d != "" {
		date, err = time.Parse("2006-01-02", d)
		if err != nil {
			Fail(c, http.StatusBadRequest, "date must be YYYY-MM-DD")
			return
		}
	}

	slots, err := h.svc.GetSchedule(technicianID, date)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, http.StatusOK, slots)
}
