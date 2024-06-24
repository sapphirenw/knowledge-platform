package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/sapphirenw/ai-content-creation-api/src/jobs"
)

func RunJobs(
	ctx context.Context,
	logger *slog.Logger,
) {
	logger.Info("Initializing job runner")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				if err := jobs.VectorizeDatastoreRunner(ctx, logger); err != nil {
					logger.Error("Error running vectorize job", "error", err)
				}
			}()
		}
	}
}
