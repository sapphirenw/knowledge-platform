package jobs

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/customer"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

// query all the active vectorize requests and vectorize the datastore
func VectorizeDatastoreRunner(
	ctx context.Context,
	logger *slog.Logger,
) error {
	pool, err := db.GetPool()
	if err != nil {
		return slogger.Error(ctx, logger, "failed to get the database pool", err)
	}
	dmodel := queries.New(pool)

	// get the waiting jobs
	jobs, err := dmodel.GetVectorizeJobsStatus(ctx, queries.VectorizeStatusWaiting)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to get the vectorize jobs", err)
	}

	// process all jobs
	for _, job := range jobs {
		logger.InfoContext(ctx, "Processing job", "job", *job)

		// get the customer
		c, err := customer.NewCustomer(ctx, logger, job.CustomerID, pool)
		if err != nil {
			// set the job status as rejected
			if _, err := dmodel.UpdateVectorizeJobStatus(ctx, &queries.UpdateVectorizeJobStatusParams{
				ID:      job.ID,
				Status:  queries.VectorizeStatusRejected,
				Message: "There is no customer with this id",
			}); err != nil {
				slogger.Error(ctx, logger, "failed to update the job", err)
			}
			slogger.Error(ctx, logger, "failed to get the customer", err)
			continue
		}

		// process the request
		if err := c.VectorizeDatastore(ctx, pool, job); err != nil {
			// update the status
			if _, err := dmodel.UpdateVectorizeJobStatus(ctx, &queries.UpdateVectorizeJobStatusParams{
				ID:      job.ID,
				Status:  queries.VectorizeStatusError,
				Message: fmt.Sprintf("There was an issue running the vectorization request: %s", err),
			}); err != nil {
				slogger.Error(ctx, logger, "failed to update the job", err)
			}
			slogger.Error(ctx, logger, "failed to process the vectorization request", err)
			continue
		}

		// update the status to complete
		if _, err := dmodel.UpdateVectorizeJobStatus(ctx, &queries.UpdateVectorizeJobStatusParams{
			ID:      job.ID,
			Status:  queries.VectorizeStatusComplete,
			Message: "Successfully vectorized the datastore",
		}); err != nil {
			slogger.Error(ctx, logger, "failed to update the job", err)
			continue
		}
	}

	return nil
}
