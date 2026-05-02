package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobStore struct {
	db *gorm.DB
}

func NewJobStore(db *gorm.DB) *JobStore {
	return &JobStore{db: db}
}

func (s *JobStore) Create(tx *gorm.DB, j *models.Job) error {
	result := tx.Create(j)
	return result.Error
}

func (s *JobStore) GetByID(id uuid.UUID) (*models.Job, error) {
	var j models.Job
	result := s.db.Preload("Technician").Preload("Quote").Preload("Manager").
		First(&j, "id = ?", id)
	return &j, result.Error
}

func (s *JobStore) ListByTechnicianAndDate(technicianID uuid.UUID, date time.Time) ([]models.Job, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var jobs []models.Job
	result := s.db.Select("id, starts_at, ends_at, status").
		Where("technician_id = ? AND starts_at >= ? AND starts_at < ?", technicianID, dayStart, dayEnd).
		Order("starts_at asc").Find(&jobs)
	return jobs, result.Error
}

func (s *JobStore) UpdateStatus(id uuid.UUID, status models.JobStatus) error {
	result := s.db.Model(&models.Job{}).Where("id = ?", id).Update("status", status)
	return result.Error
}

func (s *JobStore) ConflictCheck(tx *gorm.DB, technicianID uuid.UUID, startsAt, endsAt time.Time) ([]models.Job, error) {
	var jobs []models.Job
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("technician_id = ? AND status != ? AND starts_at < ? AND ends_at > ?",
			technicianID, models.JobStatusCancelled, endsAt, startsAt).
		Find(&jobs)
	return jobs, result.Error
}
