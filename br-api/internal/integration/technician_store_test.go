package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTechnicianStore_ListAndGetByID(t *testing.T) {
	env := setupDB(t)
	t1 := models.Technician{Name: "Tom", Email: "tom@test.com"}
	t2 := models.Technician{Name: "Lisa", Email: "lisa@test.com"}
	require.NoError(t, env.db.Create(&t1).Error)
	require.NoError(t, env.db.Create(&t2).Error)

	list, err := env.technicians.List(1, 10)
	require.NoError(t, err)
	assert.Len(t, list, 2)

	got, err := env.technicians.GetByID(t1.ID)
	require.NoError(t, err)
	assert.Equal(t, "Tom", got.Name)
}

func TestTechnicianStore_GetByID_NotFound(t *testing.T) {
	env := setupDB(t)
	_, err := env.technicians.GetByID(uuid.New())
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestTechnicianStore_List_Pagination(t *testing.T) {
	env := setupDB(t)
	for i := 0; i < 5; i++ {
		tech := models.Technician{Name: "Tech", Email: uuid.New().String() + "@test.com"}
		require.NoError(t, env.db.Create(&tech).Error)
	}

	page1, err := env.technicians.List(1, 3)
	require.NoError(t, err)
	assert.Len(t, page1, 3)

	page2, err := env.technicians.List(2, 3)
	require.NoError(t, err)
	assert.Len(t, page2, 2)
}
