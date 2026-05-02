package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/db"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/store"
	"github.com/stretchr/testify/require"
	mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

type testEnv struct {
	db            *gorm.DB
	jobs          *store.JobStore
	quotes        *store.QuoteStore
	notifications *store.NotificationStore
	technicians   *store.TechnicianStore
	managers      *store.ManagerStore
}

func setupDB(t *testing.T) *testEnv {
	t.Helper()
	ctx := context.Background()

	container, err := mysql.Run(ctx,
		"mysql:8",
		mysql.WithDatabase("brix"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/brix?parseTime=true", host, port.Port())
	database, err := db.Connect(dsn)
	require.NoError(t, err)

	return &testEnv{
		db:            database,
		jobs:          store.NewJobStore(database),
		quotes:        store.NewQuoteStore(database),
		notifications: store.NewNotificationStore(database),
		technicians:   store.NewTechnicianStore(database),
		managers:      store.NewManagerStore(database),
	}
}

func seedTechnicianAndManager(t *testing.T, env *testEnv) (models.Technician, models.Manager) {
	t.Helper()
	tech := models.Technician{Name: "Test Tech", Email: uuid.New().String() + "@test.com"}
	require.NoError(t, env.db.Create(&tech).Error)
	mgr := models.Manager{Name: "Test Manager", Email: uuid.New().String() + "@test.com"}
	require.NoError(t, env.db.Create(&mgr).Error)
	return tech, mgr
}

func seedQuote(t *testing.T, env *testEnv) models.Quote {
	t.Helper()
	q := models.Quote{CustomerName: "Test Customer", Address: "1 Test St", Description: "Test job"}
	require.NoError(t, env.db.Create(&q).Error)
	return q
}
