package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/logger"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
	"gorm.io/gorm"
)

type notificationPublisher interface {
	Publish(n *models.Notification)
}

type JobService struct {
	jobs             jobStore
	quotes           quoteStore
	notifications    notificationStore
	publisher        notificationPublisher
	startTransaction func(opts ...*sql.TxOptions) *gorm.DB
}

func NewJobService(db *gorm.DB, jobs *store.JobStore, quotes *store.QuoteStore, notifications *store.NotificationStore, publisher notificationPublisher) *JobService {
	return &JobService{
		jobs:             jobs,
		quotes:           quotes,
		notifications:    notifications,
		publisher:        publisher,
		startTransaction: db.Begin,
	}
}

func (svc *JobService) GetByID(id uuid.UUID) (*JobResponse, error) {
	job, err := svc.jobs.GetByID(id)
	if err != nil {
		return nil, err
	}
	response := ToJobResponse(job)
	return &response, nil
}

func (svc *JobService) AssignJob(input AssignJobInput) (*JobResponse, error) {
	if !input.StartsAt.After(time.Now()) {
		return nil, ErrStartsInPast
	}

	endTime := input.StartsAt.Add(2 * time.Hour)

	transaction := svc.startTransaction()
	if transaction.Error != nil {
		return nil, transaction.Error
	}

	quote, err := svc.quotes.GetByIDForUpdate(transaction, input.QuoteID)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}
	if quote.Status != models.QuoteStatusUnscheduled {
		transaction.Rollback()
		return nil, ErrQuoteNotUnscheduled
	}

	conflicts, err := svc.jobs.ConflictCheck(transaction, input.TechnicianID, input.StartsAt, endTime)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}
	if len(conflicts) > 0 {
		transaction.Rollback()
		return nil, ErrConflict
	}

	job := &models.Job{
		QuoteID:      input.QuoteID,
		TechnicianID: input.TechnicianID,
		ManagerID:    input.ManagerID,
		StartsAt:     input.StartsAt,
		EndsAt:       endTime,
		Status:       models.JobStatusScheduled,
	}
	if err := svc.jobs.Create(transaction, job); err != nil {
		transaction.Rollback()
		return nil, err
	}

	if err := svc.quotes.UpdateStatus(transaction, input.QuoteID, models.QuoteStatusScheduled); err != nil {
		transaction.Rollback()
		return nil, err
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, err
	}

	n := &models.Notification{
		Type:          models.NotificationTypeJobAssigned,
		RecipientType: models.RecipientTypeTechnician,
		RecipientID:   input.TechnicianID,
		JobID:         job.ID,
		Message:       fmt.Sprintf("You have been assigned job %s starting at %s", job.ID, job.StartsAt.Format(time.RFC3339)),
	}
	if err := svc.notifications.Create(n); err != nil {
		logger.Log.Error("failed to create notification", "error", err, "job_id", job.ID)
	}
	svc.publisher.Publish(n)

	response := ToJobResponse(job)
	return &response, nil
}

func (svc *JobService) CompleteJob(jobID, technicianID uuid.UUID) (*JobResponse, error) {
	job, err := svc.jobs.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	if job.TechnicianID != technicianID {
		return nil, ErrUnauthorised
	}

	if job.Status != models.JobStatusScheduled {
		return nil, ErrJobNotScheduled
	}

	if err := svc.jobs.UpdateStatus(jobID, models.JobStatusCompleted); err != nil {
		return nil, err
	}

	job.Status = models.JobStatusCompleted

	m := &models.Notification{
		Type:          models.NotificationTypeJobCompleted,
		RecipientType: models.RecipientTypeManager,
		RecipientID:   job.ManagerID,
		JobID:         job.ID,
		Message:       fmt.Sprintf("Job %s has been completed by technician %s", job.ID, job.TechnicianID),
	}
	if err := svc.notifications.Create(m); err != nil {
		logger.Log.Error("failed to create notification", "error", err, "job_id", job.ID)
	}
	svc.publisher.Publish(m)

	response := ToJobResponse(job)
	return &response, nil
}
