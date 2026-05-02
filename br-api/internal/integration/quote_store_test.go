package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuoteStore_CreateAndGetByID(t *testing.T) {
	env := setupDB(t)
	q := models.Quote{CustomerName: "Alice", Address: "1 St", Description: "Paint walls"}
	require.NoError(t, env.quotes.Create(&q))
	assert.NotEqual(t, uuid.Nil, q.ID)

	got, err := env.quotes.GetByID(q.ID)
	require.NoError(t, err)
	assert.Equal(t, "Alice", got.CustomerName)
	assert.Equal(t, models.QuoteStatusUnscheduled, got.Status)
}

func TestQuoteStore_List_FiltersByStatus(t *testing.T) {
	env := setupDB(t)
	q1 := models.Quote{CustomerName: "A", Address: "1 St", Status: models.QuoteStatusUnscheduled}
	q2 := models.Quote{CustomerName: "B", Address: "2 St", Status: models.QuoteStatusUnscheduled}
	q3 := models.Quote{CustomerName: "C", Address: "3 St", Status: models.QuoteStatusScheduled}
	require.NoError(t, env.quotes.Create(&q1))
	require.NoError(t, env.quotes.Create(&q2))
	require.NoError(t, env.quotes.Create(&q3))

	results, err := env.quotes.List(models.QuoteStatusUnscheduled, 1, 20)
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestQuoteStore_List_Pagination(t *testing.T) {
	env := setupDB(t)
	for i := 0; i < 5; i++ {
		q := models.Quote{CustomerName: "Customer", Address: "1 St"}
		require.NoError(t, env.quotes.Create(&q))
	}

	page1, err := env.quotes.List(models.QuoteStatusUnscheduled, 1, 2)
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page3, err := env.quotes.List(models.QuoteStatusUnscheduled, 3, 2)
	require.NoError(t, err)
	assert.Len(t, page3, 1)
}

func TestQuoteStore_UpdateStatus(t *testing.T) {
	env := setupDB(t)
	q := models.Quote{CustomerName: "Bob", Address: "2 St"}
	require.NoError(t, env.quotes.Create(&q))

	require.NoError(t, env.quotes.UpdateStatus(env.db, q.ID, models.QuoteStatusScheduled))

	got, err := env.quotes.GetByID(q.ID)
	require.NoError(t, err)
	assert.Equal(t, models.QuoteStatusScheduled, got.Status)
}
