package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobStore_CreateAndGetByID(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)

	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	job := &models.Job{
		QuoteID:      quote.ID,
		TechnicianID: tech.ID,
		ManagerID:    mgr.ID,
		StartsAt:     startsAt,
		EndsAt:       startsAt.Add(2 * time.Hour),
		Status:       models.JobStatusScheduled,
	}
	require.NoError(t, env.jobs.Create(env.db, job))
	assert.NotEqual(t, uuid.Nil, job.ID)

	got, err := env.jobs.GetByID(job.ID)
	require.NoError(t, err)
	assert.Equal(t, models.JobStatusScheduled, got.Status)
	assert.Equal(t, tech.ID, got.TechnicianID)
}

func TestJobStore_ListByTechnicianAndDate(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote1 := seedQuote(t, env)
	quote2 := seedQuote(t, env)

	today := time.Now().UTC().Truncate(time.Second)
	tomorrow := today.Add(24 * time.Hour)

	job1 := &models.Job{QuoteID: quote1.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: today.Add(1 * time.Hour), EndsAt: today.Add(3 * time.Hour), Status: models.JobStatusScheduled}
	job2 := &models.Job{QuoteID: quote2.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: tomorrow.Add(1 * time.Hour), EndsAt: tomorrow.Add(3 * time.Hour), Status: models.JobStatusScheduled}
	require.NoError(t, env.jobs.Create(env.db, job1))
	require.NoError(t, env.jobs.Create(env.db, job2))

	results, err := env.jobs.ListByTechnicianAndDate(tech.ID, today)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, job1.ID, results[0].ID)
}

func TestJobStore_ConflictCheck(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)

	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	endsAt := startsAt.Add(2 * time.Hour)

	job := &models.Job{QuoteID: quote.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: startsAt, EndsAt: endsAt, Status: models.JobStatusScheduled}
	require.NoError(t, env.jobs.Create(env.db, job))

	tx := env.db.Begin()
	conflicts, err := env.jobs.ConflictCheck(tx, tech.ID, startsAt.Add(30*time.Minute), endsAt.Add(30*time.Minute))
	tx.Rollback()
	require.NoError(t, err)
	assert.Len(t, conflicts, 1)

	tx2 := env.db.Begin()
	none, err := env.jobs.ConflictCheck(tx2, tech.ID, endsAt.Add(1*time.Hour), endsAt.Add(3*time.Hour))
	tx2.Rollback()
	require.NoError(t, err)
	assert.Len(t, none, 0)
}

func TestJobStore_UpdateStatus(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)

	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	job := &models.Job{QuoteID: quote.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: startsAt, EndsAt: startsAt.Add(2 * time.Hour), Status: models.JobStatusScheduled}
	require.NoError(t, env.jobs.Create(env.db, job))

	require.NoError(t, env.jobs.UpdateStatus(job.ID, models.JobStatusCompleted))

	got, err := env.jobs.GetByID(job.ID)
	require.NoError(t, err)
	assert.Equal(t, models.JobStatusCompleted, got.Status)
}
