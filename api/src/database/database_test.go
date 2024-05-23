package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	// test working
	pool, err := GetPool()
	require.Nil(t, err)
	t.Cleanup(func() {
		ClosePool()
	})

	require.NotNil(t, pool)
	err = pool.Ping(context.TODO())
	require.Nil(t, err)
	ClosePool() // close the pool

	// test not working
	DATABASE_URL = "postgres://invalid"
	_, err = GetPool()
	require.Error(t, err)
}
