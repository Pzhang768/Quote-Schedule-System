package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationStore_CreateAndList(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)
	job := &models.Job{QuoteID: quote.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(3 * time.Hour), Status: models.JobStatusScheduled}
	require.NoError(t, env.db.Create(job).Error)

	n := &models.Notification{
		Type: models.NotificationTypeJobAssigned, RecipientType: models.RecipientTypeTechnician,
		RecipientID: tech.ID, JobID: job.ID, Message: "You have a job",
	}
	require.NoError(t, env.notifications.Create(n))
	assert.NotEqual(t, uuid.Nil, n.ID)

	list, err := env.notifications.List(models.RecipientTypeTechnician, tech.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "You have a job", list[0].Message)
}

func TestNotificationStore_ListSince(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote1 := seedQuote(t, env)
	quote2 := seedQuote(t, env)
	job1 := &models.Job{QuoteID: quote1.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(3 * time.Hour), Status: models.JobStatusScheduled}
	job2 := &models.Job{QuoteID: quote2.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: time.Now().Add(4 * time.Hour), EndsAt: time.Now().Add(6 * time.Hour), Status: models.JobStatusScheduled}
	require.NoError(t, env.db.Create(job1).Error)
	require.NoError(t, env.db.Create(job2).Error)

	before := time.Now()
	time.Sleep(10 * time.Millisecond)

	require.NoError(t, env.notifications.Create(&models.Notification{
		Type: models.NotificationTypeJobAssigned, RecipientType: models.RecipientTypeTechnician,
		RecipientID: tech.ID, JobID: job1.ID, Message: "Job 1",
	}))
	require.NoError(t, env.notifications.Create(&models.Notification{
		Type: models.NotificationTypeJobAssigned, RecipientType: models.RecipientTypeTechnician,
		RecipientID: tech.ID, JobID: job2.ID, Message: "Job 2",
	}))

	results, err := env.notifications.ListSince(models.RecipientTypeTechnician, tech.ID, before)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestNotificationStore_Read(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)
	job := &models.Job{QuoteID: quote.ID, TechnicianID: tech.ID, ManagerID: mgr.ID,
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(3 * time.Hour), Status: models.JobStatusScheduled}
	require.NoError(t, env.db.Create(job).Error)

	n := &models.Notification{
		Type: models.NotificationTypeJobAssigned, RecipientType: models.RecipientTypeTechnician,
		RecipientID: tech.ID, JobID: job.ID, Message: "msg",
	}
	require.NoError(t, env.notifications.Create(n))
	assert.Nil(t, n.ReadAt)

	require.NoError(t, env.notifications.Read(n.ID, tech.ID))

	var updated models.Notification
	require.NoError(t, env.db.First(&updated, "id = ?", n.ID).Error)
	assert.NotNil(t, updated.ReadAt)
}
