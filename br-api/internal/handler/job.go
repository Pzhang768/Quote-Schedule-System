package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/service"
)

type JobHandler struct {
	svc *service.JobService
}

func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// @Summary     Get a job by ID
// @Tags        jobs
// @Produce     json
// @Param       id   path      string  true  "Job ID"
// @Success     200  {object}  Response[service.JobResponse]
// @Failure     400  {object}  ErrorResponse
// @Failure     404  {object}  ErrorResponse
// @Router      /jobs/{id} [get]
func (h *JobHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid job id")
		return
	}
	job, err := h.svc.GetByID(id)
	if err != nil {
		Fail(c, http.StatusNotFound, "job not found")
		return
	}
	Success(c, http.StatusOK, job)
}

type assignJobRequest struct {
	QuoteID      uuid.UUID `json:"quote_id" binding:"required"`
	TechnicianID uuid.UUID `json:"technician_id" binding:"required"`
	ManagerID    uuid.UUID `json:"manager_id" binding:"required"`
	StartsAt     time.Time `json:"starts_at" binding:"required"`
}

// @Summary     Assign a quote to a technician
// @Tags        jobs
// @Accept      json
// @Produce     json
// @Param       body  body      assignJobRequest  true  "Assignment"
// @Success     201   {object}  Response[service.JobResponse]
// @Failure     409   {object}  ErrorResponse  "Technician has a conflicting job"
// @Failure     422   {object}  ErrorResponse
// @Router      /jobs [post]
func (h *JobHandler) Assign(c *gin.Context) {
	var body assignJobRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	job, err := h.svc.AssignJob(service.AssignJobInput{
		QuoteID:      body.QuoteID,
		TechnicianID: body.TechnicianID,
		ManagerID:    body.ManagerID,
		StartsAt:     body.StartsAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConflict):
			Fail(c, http.StatusConflict, err.Error())
		case errors.Is(err, service.ErrStartsInPast), errors.Is(err, service.ErrQuoteNotUnscheduled):
			Fail(c, http.StatusUnprocessableEntity, err.Error())
		default:
			Fail(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	Success(c, http.StatusCreated, job)
}

type completeJobRequest struct {
	TechnicianID uuid.UUID `json:"technician_id" binding:"required"`
}

// @Summary     Mark a job as complete
// @Tags        jobs
// @Accept      json
// @Produce     json
// @Param       id    path      string              true  "Job ID"
// @Param       body  body      completeJobRequest  true  "Technician"
// @Success     200   {object}  Response[service.JobResponse]
// @Failure     403   {object}  ErrorResponse
// @Failure     422   {object}  ErrorResponse
// @Router      /jobs/{id}/complete [patch]
func (h *JobHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid job id")
		return
	}

	var body completeJobRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		Fail(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	job, err := h.svc.CompleteJob(id, body.TechnicianID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorised):
			Fail(c, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrJobNotScheduled):
			Fail(c, http.StatusUnprocessableEntity, err.Error())
		default:
			Fail(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	Success(c, http.StatusOK, job)
}
