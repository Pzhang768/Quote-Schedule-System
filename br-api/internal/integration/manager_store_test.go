package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagerStore_List_Pagination(t *testing.T) {
	env := setupDB(t)
	for i := 0; i < 4; i++ {
		mgr := models.Manager{Name: "Manager", Email: uuid.New().String() + "@test.com"}
		require.NoError(t, env.db.Create(&mgr).Error)
	}

	page1, err := env.managers.List(1, 3)
	require.NoError(t, err)
	assert.Len(t, page1, 3)

	page2, err := env.managers.List(2, 3)
	require.NoError(t, err)
	assert.Len(t, page2, 1)
}
