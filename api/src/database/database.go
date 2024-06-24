package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool      *pgxpool.Pool
	poolMutex sync.Mutex
)

var DATABASE_URL = ""

// create the connection on startup
// func init() {
// 	if os.Getenv("DATABASE_URL") == "" {
// 		panic("the env variable `DATABASE_URL` is required")
// 	}
// 	err := initializePool()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// initializePool safely initializes the pool, ensuring it's only done once.
func initializePool() error {
	var err error
	if DATABASE_URL == "" {
		fmt.Println("WARNING -- using default env variable for database")
		DATABASE_URL = os.Getenv("DATABASE_URL")
	}
	pool, err = pgxpool.New(context.Background(), DATABASE_URL)
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

	// make sure the database is functional
	if err := pool.Ping(context.Background()); err != nil {
		fmt.Printf("there was an issue pinging the db pool: %s\n", err)
		if strings.Contains(err.Error(), "close") {
			fmt.Println("Attemting to re-open the connection ...")
			// re-init the pool
			err := ReinitializePool()
			if err != nil {
				return nil, fmt.Errorf("failed to recover from an error in the pool: %s", err)
			}
		} else {
			return nil, err
		}
	}

	return pool, nil
}

// ReinitializePool safely closes the existing pool and creates a new one. This can be triggered on detecting a failover or similar event.
func ReinitializePool() error {
	if pool != nil {
		pool.Close()
		pool = nil
	}

	return initializePool()
}

func ClosePool() {
	if pool != nil {
		pool.Close()
		pool = nil
	}
}
