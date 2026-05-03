package service

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/hub"
	"github.com/melfish/br-api/internal/mocks"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmock "github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// stubDB returns a gorm.DB backed by sqlmock. Tests that don't care about
// transaction behaviour can ignore the mock — it's only here because startTransaction
// requires a real *gorm.DB to call Begin() on.
func stubDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return gormDB, mock
}

// noTx is a startTransaction stub for tests that never reach the transaction path.
func noTx(_ ...*sql.TxOptions) *gorm.DB { return nil }

type jobTestEnv struct {
	svc    *JobService
	jobs   *mocks.JobStore
	quotes *mocks.QuoteStore
	notifs *mocks.NotificationStore
}

func setupJobTest(t *testing.T, startTransaction func(...*sql.TxOptions) *gorm.DB) *jobTestEnv {
	t.Helper()
	jobs := &mocks.JobStore{}
	quotes := &mocks.QuoteStore{}
	notifs := &mocks.NotificationStore{}
	svc := &JobService{jobs: jobs, quotes: quotes, notifications: notifs, publisher: hub.New(), startTransaction: startTransaction}
	return &jobTestEnv{svc: svc, jobs: jobs, quotes: quotes, notifs: notifs}
}

func TestAssignJob(t *testing.T) {
	quoteID := uuid.New()
	techID := uuid.New()
	managerID := uuid.New()
	startsAt := time.Now().Add(1 * time.Hour)
	validInput := AssignJobInput{QuoteID: quoteID, TechnicianID: techID, ManagerID: managerID, StartsAt: startsAt}

	t.Run("assigns job successfully", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectCommit()
		env := setupJobTest(t, db.Begin)

		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return(&models.Quote{ID: quoteID, Status: models.QuoteStatusUnscheduled}, nil)
		env.jobs.On("ConflictCheck", tmock.AnythingOfType("*gorm.DB"), techID, startsAt, startsAt.Add(2*time.Hour)).Return([]models.Job{}, nil)
		env.jobs.On("Create", tmock.AnythingOfType("*gorm.DB"), tmock.AnythingOfType("*models.Job")).Return(nil)
		env.quotes.On("UpdateStatus", tmock.AnythingOfType("*gorm.DB"), quoteID, models.QuoteStatusScheduled).Return(nil)
		env.notifs.On("Create", tmock.AnythingOfType("*models.Notification")).Return(nil)

		response, err := env.svc.AssignJob(validInput)
		assert.NoError(t, err)
		assert.Equal(t, models.JobStatusScheduled, response.Status)
	})

	t.Run("rejects when starts_at is in the past", func(t *testing.T) {
		env := setupJobTest(t, noTx)
		_, err := env.svc.AssignJob(AssignJobInput{StartsAt: time.Now().Add(-1 * time.Hour)})
		assert.ErrorIs(t, err, ErrStartsInPast)
	})

	t.Run("rejects when quote is already scheduled", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()
		env := setupJobTest(t, db.Begin)
		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return(&models.Quote{ID: quoteID, Status: models.QuoteStatusScheduled}, nil)

		_, err := env.svc.AssignJob(validInput)
		assert.ErrorIs(t, err, ErrQuoteNotUnscheduled)
	})

	t.Run("rejects when technician has a conflicting job", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()
		env := setupJobTest(t, db.Begin)

		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return(&models.Quote{ID: quoteID, Status: models.QuoteStatusUnscheduled}, nil)
		env.jobs.On("ConflictCheck", tmock.AnythingOfType("*gorm.DB"), techID, startsAt, startsAt.Add(2*time.Hour)).
			Return([]models.Job{{ID: uuid.New()}}, nil)

		_, err := env.svc.AssignJob(validInput)
		assert.ErrorIs(t, err, ErrConflict)
	})

	t.Run("returns error when quote lookup fails", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()
		env := setupJobTest(t, db.Begin)
		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return((*models.Quote)(nil), errors.New("db error"))

		_, err := env.svc.AssignJob(validInput)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("returns error when conflict check fails", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()
		env := setupJobTest(t, db.Begin)

		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return(&models.Quote{ID: quoteID, Status: models.QuoteStatusUnscheduled}, nil)
		env.jobs.On("ConflictCheck", tmock.AnythingOfType("*gorm.DB"), techID, startsAt, startsAt.Add(2*time.Hour)).
			Return([]models.Job{}, errors.New("db error"))

		_, err := env.svc.AssignJob(validInput)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("returns error when job create fails", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()
		env := setupJobTest(t, db.Begin)

		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return(&models.Quote{ID: quoteID, Status: models.QuoteStatusUnscheduled}, nil)
		env.jobs.On("ConflictCheck", tmock.AnythingOfType("*gorm.DB"), techID, startsAt, startsAt.Add(2*time.Hour)).Return([]models.Job{}, nil)
		env.jobs.On("Create", tmock.AnythingOfType("*gorm.DB"), tmock.AnythingOfType("*models.Job")).Return(errors.New("db error"))

		_, err := env.svc.AssignJob(validInput)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("returns error when quote status update fails", func(t *testing.T) {
		db, mock := stubDB(t)
		mock.ExpectBegin()
		mock.ExpectRollback()
		env := setupJobTest(t, db.Begin)

		env.quotes.On("GetByIDForUpdate", tmock.AnythingOfType("*gorm.DB"), quoteID).Return(&models.Quote{ID: quoteID, Status: models.QuoteStatusUnscheduled}, nil)
		env.jobs.On("ConflictCheck", tmock.AnythingOfType("*gorm.DB"), techID, startsAt, startsAt.Add(2*time.Hour)).Return([]models.Job{}, nil)
		env.jobs.On("Create", tmock.AnythingOfType("*gorm.DB"), tmock.AnythingOfType("*models.Job")).Return(nil)
		env.quotes.On("UpdateStatus", tmock.AnythingOfType("*gorm.DB"), quoteID, models.QuoteStatusScheduled).Return(errors.New("db error"))

		_, err := env.svc.AssignJob(validInput)
		assert.ErrorContains(t, err, "db error")
	})
}

func TestCompleteJob(t *testing.T) {
	jobID := uuid.New()
	techID := uuid.New()
	managerID := uuid.New()

	t.Run("completes job successfully", func(t *testing.T) {
		env := setupJobTest(t, noTx)
		env.jobs.On("GetByID", jobID).Return(&models.Job{ID: jobID, TechnicianID: techID, ManagerID: managerID, Status: models.JobStatusScheduled}, nil)
		env.jobs.On("UpdateStatus", jobID, models.JobStatusCompleted).Return(nil)
		env.notifs.On("Create", tmock.AnythingOfType("*models.Notification")).Return(nil)

		response, err := env.svc.CompleteJob(jobID, techID)
		assert.NoError(t, err)
		assert.Equal(t, models.JobStatusCompleted, response.Status)
	})

	t.Run("rejects when technician does not own the job", func(t *testing.T) {
		env := setupJobTest(t, noTx)
		env.jobs.On("GetByID", jobID).Return(&models.Job{ID: jobID, TechnicianID: techID, Status: models.JobStatusScheduled}, nil)

		_, err := env.svc.CompleteJob(jobID, uuid.New())
		assert.ErrorIs(t, err, ErrUnauthorised)
	})

	t.Run("rejects when job is not in scheduled state", func(t *testing.T) {
		env := setupJobTest(t, noTx)
		env.jobs.On("GetByID", jobID).Return(&models.Job{ID: jobID, TechnicianID: techID, Status: models.JobStatusCompleted}, nil)

		_, err := env.svc.CompleteJob(jobID, techID)
		assert.ErrorIs(t, err, ErrJobNotScheduled)
	})

	t.Run("returns error when job lookup fails", func(t *testing.T) {
		env := setupJobTest(t, noTx)
		env.jobs.On("GetByID", jobID).Return((*models.Job)(nil), errors.New("db error"))

		_, err := env.svc.CompleteJob(jobID, techID)
		assert.ErrorContains(t, err, "db error")
	})

	t.Run("returns error when status update fails", func(t *testing.T) {
		env := setupJobTest(t, noTx)
		env.jobs.On("GetByID", jobID).Return(&models.Job{ID: jobID, TechnicianID: techID, ManagerID: managerID, Status: models.JobStatusScheduled}, nil)
		env.jobs.On("UpdateStatus", jobID, models.JobStatusCompleted).Return(errors.New("db error"))

		_, err := env.svc.CompleteJob(jobID, techID)
		assert.ErrorContains(t, err, "db error")
	})
}
