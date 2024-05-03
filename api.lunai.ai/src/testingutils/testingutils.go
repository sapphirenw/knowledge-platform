package testingutils

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetDatabase(t *testing.T, ctx context.Context) *pgxpool.Pool {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("ankane/pgvector"),
		postgres.WithInitScripts(filepath.Join("..", "..", "..", "database", "schema", "00_GENERATED_schema.sql")),
		postgres.WithDatabase("aicontent"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	// ensure the function is cleaned
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	// get the connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	db.DATABASE_URL = connStr

	// create a new db pool
	pool, err := db.GetPool()
	require.NoError(t, err)

	t.Cleanup(func() {
		db.ClosePool()
	})

	return pool
}

func CreateTestCustomer(t *testing.T, ctx context.Context, db queries.DBTX) *queries.Customer {
	model := queries.New(db)
	customer, err := model.CreateCustomer(ctx, "test-customer")
	require.NoError(t, err)
	return customer
}
