package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
	"gorm.io/gorm"
)

type JobService struct {
	jobs             jobStore
	quotes           quoteStore
	notifications    notificationStore
	startTransaction func(opts ...*sql.TxOptions) *gorm.DB
}

func NewJobService(db *gorm.DB, jobs *store.JobStore, quotes *store.QuoteStore, notifications *store.NotificationStore) *JobService {
	return &JobService{
		jobs:             jobs,
		quotes:           quotes,
		notifications:    notifications,
		startTransaction: db.Begin,
	}
}

func (svc *JobService) AssignJob(input AssignJobInput) (*JobResponse, error) {
	if !input.StartsAt.After(time.Now()) {
		return nil, ErrStartsInPast
	}

	quote, err := svc.quotes.GetByID(input.QuoteID)
	if err != nil {
		return nil, err
	}
	if quote.Status != models.QuoteStatusUnscheduled {
		return nil, ErrQuoteNotUnscheduled
	}

	endTime := input.StartsAt.Add(2 * time.Hour)

	transaction := svc.startTransaction()
	if transaction.Error != nil {
		return nil, transaction.Error
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

	if err := svc.quotes.UpdateStatus(input.QuoteID, models.QuoteStatusScheduled); err != nil {
		transaction.Rollback()
		return nil, err
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, err
	}

	_ = svc.notifications.Create(&models.Notification{
		Type:          models.NotificationTypeJobAssigned,
		RecipientType: models.RecipientTypeTechnician,
		RecipientID:   input.TechnicianID,
		JobID:         job.ID,
		Message:       fmt.Sprintf("You have been assigned job %s starting at %s", job.ID, job.StartsAt.Format(time.RFC3339)),
	})

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

	_ = svc.notifications.Create(&models.Notification{
		Type:          models.NotificationTypeJobCompleted,
		RecipientType: models.RecipientTypeManager,
		RecipientID:   job.ManagerID,
		JobID:         job.ID,
		Message:       fmt.Sprintf("Job %s has been completed by technician %s", job.ID, job.TechnicianID),
	})

	response := ToJobResponse(job)
	return &response, nil
}
