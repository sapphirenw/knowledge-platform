package db

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool      *pgxpool.Pool
	poolMutex sync.Mutex
)

// create the connection on startup
func init() {
	if os.Getenv("DATABASE_URL") == "" {
		panic("the env variable `DATABASE_URL` is required")
	}
	err := initializePool()
	if err != nil {
		panic(err)
	}
}

// initializePool safely initializes the pool, ensuring it's only done once.
func initializePool() error {
	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	return err
}

// GetConnection provides a thread-safe way to obtain a connection pool, initializing it if necessary.
func GetPool() (*pgxpool.Pool, error) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if pool == nil {
		err := initializePool()
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}

// ReinitializePool safely closes the existing pool and creates a new one. This can be triggered on detecting a failover or similar event.
func ReinitializePool() error {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if pool != nil {
		pool.Close()
	}

	return initializePool()
}
