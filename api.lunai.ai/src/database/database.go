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
func initializePool(connStr *string) error {
	var c string
	if connStr == nil {
		c = os.Getenv("DATABASE_URL")
	} else {
		c = *connStr
	}
	var err error
	pool, err = pgxpool.New(context.Background(), c)
	return err
}

// GetConnection provides a thread-safe way to obtain a connection pool, initializing it if necessary.
func GetPool(connStr *string) (*pgxpool.Pool, error) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if pool == nil {
		err := initializePool(connStr)
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
			err := ReinitializePool(connStr)
			if err != nil {
				return nil, fmt.Errorf("failed to recover from an error in the pool: %s", err)
			}
		}
	}

	return pool, nil
}

// ReinitializePool safely closes the existing pool and creates a new one. This can be triggered on detecting a failover or similar event.
func ReinitializePool(connStr *string) error {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if pool != nil {
		pool.Close()
	}

	return initializePool(connStr)
}
