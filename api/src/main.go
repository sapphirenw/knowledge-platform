package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/logger"
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
	logger := logger.NewLogger()
	srv := NewServer(logger)

	// set the database url
	db.DATABASE_URL = getenv("DATABASE_URL")

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(getenv("SERVER_HOST"), getenv("SERVER_PORT")),
		Handler: srv,
	}
	go func() {
		fmt.Fprintf(stdout, "listening on %s\n", httpServer.Addr)
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

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
