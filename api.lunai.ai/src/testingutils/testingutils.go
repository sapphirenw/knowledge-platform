package testingutils

import (
	"context"
	"fmt"
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

const TEST_CUSTOMER_ID = 7

func CreateTestCustomer(db queries.DBTX) (*queries.Customer, error) {
	model := queries.New(db)
	c, err := model.GetCustomer(context.TODO(), TEST_CUSTOMER_ID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the test customer: %v", err)
	}
	return c, err
}

func GetDatabase(t *testing.T, ctx context.Context) *pgxpool.Pool {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("ankane/pgvector"),
		postgres.WithInitScripts(filepath.Join("..", "..", "schema", "schema.sql")),
		postgres.WithInitScripts(filepath.Join("..", "..", "schema", "triggers.sql")),
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

	// create a new db pool
	pool, err := db.GetPool(&connStr)
	require.NoError(t, err)

	return pool
}
