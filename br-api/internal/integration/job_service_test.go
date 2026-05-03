package integration

import (
	"sync"
	"testing"
	"time"

	"github.com/melfish/br-api/internal/hub"
	"github.com/melfish/br-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssignJob_Integration(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)

	svc := service.NewJobService(env.db, env.jobs, env.quotes, env.notifications, hub.New())
	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	resp, err := svc.AssignJob(service.AssignJobInput{
		QuoteID:      quote.ID,
		TechnicianID: tech.ID,
		ManagerID:    mgr.ID,
		StartsAt:     startsAt,
	})
	require.NoError(t, err)
	assert.Equal(t, "scheduled", string(resp.Status))
	assert.Equal(t, quote.ID, resp.QuoteID)

	// quote status updated to scheduled
	updatedQuote, err := env.quotes.GetByID(quote.ID)
	require.NoError(t, err)
	assert.Equal(t, "scheduled", string(updatedQuote.Status))

	// technician notification created
	notifications, err := env.notifications.List("technician", tech.ID)
	require.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, "job_assigned", string(notifications[0].Type))
}

func TestAssignJob_Conflict_Integration(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote1 := seedQuote(t, env)
	quote2 := seedQuote(t, env)

	svc := service.NewJobService(env.db, env.jobs, env.quotes, env.notifications, hub.New())
	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	_, err := svc.AssignJob(service.AssignJobInput{
		QuoteID:      quote1.ID,
		TechnicianID: tech.ID,
		ManagerID:    mgr.ID,
		StartsAt:     startsAt,
	})
	require.NoError(t, err)

	// same technician, overlapping window
	_, err = svc.AssignJob(service.AssignJobInput{
		QuoteID:      quote2.ID,
		TechnicianID: tech.ID,
		ManagerID:    mgr.ID,
		StartsAt:     startsAt.Add(30 * time.Minute),
	})
	assert.ErrorIs(t, err, service.ErrConflict)
}

func TestAssignJob_ConcurrentConflict_Integration(t *testing.T) {
	// Two goroutines race to book the same technician. SELECT FOR UPDATE ensures only one wins.
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote1 := seedQuote(t, env)
	quote2 := seedQuote(t, env)

	svc := service.NewJobService(env.db, env.jobs, env.quotes, env.notifications, hub.New())
	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, errs[0] = svc.AssignJob(service.AssignJobInput{
			QuoteID: quote1.ID, TechnicianID: tech.ID, ManagerID: mgr.ID, StartsAt: startsAt,
		})
	}()
	go func() {
		defer wg.Done()
		_, errs[1] = svc.AssignJob(service.AssignJobInput{
			QuoteID: quote2.ID, TechnicianID: tech.ID, ManagerID: mgr.ID, StartsAt: startsAt,
		})
	}()
	wg.Wait()

	// exactly one should succeed, one should get ErrConflict
	successes := 0
	for _, err := range errs {
		if err == nil {
			successes++
		}
	}
	assert.Equal(t, 1, successes)
}

func TestCompleteJob_Integration(t *testing.T) {
	env := setupDB(t)
	tech, mgr := seedTechnicianAndManager(t, env)
	quote := seedQuote(t, env)

	svc := service.NewJobService(env.db, env.jobs, env.quotes, env.notifications, hub.New())
	startsAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	assigned, err := svc.AssignJob(service.AssignJobInput{
		QuoteID:      quote.ID,
		TechnicianID: tech.ID,
		ManagerID:    mgr.ID,
		StartsAt:     startsAt,
	})
	require.NoError(t, err)

	completed, err := svc.CompleteJob(assigned.ID, tech.ID)
	require.NoError(t, err)
	assert.Equal(t, "completed", string(completed.Status))

	// manager notification created
	notifications, err := env.notifications.List("manager", mgr.ID)
	require.NoError(t, err)
	assert.Len(t, notifications, 1)
	assert.Equal(t, "job_completed", string(notifications[0].Type))
}
