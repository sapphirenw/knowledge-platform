package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout, stderr io.Writer,
) error {
	logger := slogger.NewLogger()
	srv := NewServer(logger)

	// set the database url
	db.DATABASE_URL = getenv("DATABASE_URL")

	// ensure the database can be reached
	if _, err := db.GetPool(); err != nil {
		logger.Warn("Failed to connect to database on first pass, waiting ...")
		retries := 0
		for {
			if retries > 3 {
				panic("Could not connect to the database!")
			}
			time.Sleep(time.Second * 10)
			if _, err := db.GetPool(); err == nil {
				logger.Info("Successfully connected to database")
				break
			}
			retries += 1
			logger.Warn("Failed attempt to connect to database", "attempt", retries)
		}
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(getenv("SERVER_HOST"), getenv("SERVER_PORT")),
		Handler: srv,
	}

	// run the server on a thread
	go func() {
		fmt.Fprintf(stdout, "listening on %s\n", httpServer.Addr)
		logger.Info("Starting server", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	// run the jobs on a thread
	go RunJobs(ctx, logger.Logger)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown (thanks Alessandro Rosetti)
		shutdownCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
