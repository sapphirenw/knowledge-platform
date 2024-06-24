package customer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

// creates a request to vectorize the data
func createVectorizeRequest(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the request
	body, valid := request.Decode[createVectorRequest](w, r, c.logger)
	if !valid {
		return
	}

	// start the transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to connect to the database", "error", err)
		return
	}
	defer tx.Commit(r.Context())

	dmodel := queries.New(tx)
	response, err := dmodel.CreateVectorizeJob(r.Context(), &queries.CreateVectorizeJobParams{
		CustomerID: c.ID,
		Documents:  body.Documents,
		Websites:   body.Websites,
	})
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to create the vectorize job", "error", err)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

// get a vectorize request
func getVectorizeRequest(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	jobId, err := utils.GoogleUUIDFromString(chi.URLParam(r, "id"))
	if err != nil {
		slogger.ServerError(w, r, c.logger, 400, "failed to parse the id", "error", err)
		return
	}

	// start the transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to connect to the database", "error", err)
		return
	}
	defer tx.Commit(r.Context())

	dmodel := queries.New(tx)
	response, err := dmodel.GetVectorizeJob(r.Context(), jobId)
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to get the vectorize job", "error", err)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

// get all the vectorize requests
func getAllVectorizeRequests(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// start the transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to connect to the database", "error", err)
		return
	}
	defer tx.Commit(r.Context())

	dmodel := queries.New(tx)
	response, err := dmodel.GetCustomerVectorizeJobs(r.Context(), c.ID)
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to get the vectorize jobs", "error", err)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

func queryVectorStore(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the request
	body, valid := request.Decode[queryVectorStoreRequest](w, r, c.logger)
	if !valid {
		return
	}

	// start the transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	response, err := c.QueryVectorStore(r.Context(), tx, &body)
	if err != nil {
		tx.Rollback(r.Context())
		c.logger.Error("failed to query the vectorstore", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)

}
func (c *Customer) QueryVectorStore(ctx context.Context, db queries.DBTX, request *queryVectorStoreRequest) (*queryVectorStoreResponse, error) {
	logger := c.logger.With("request", request)
	logger.InfoContext(ctx, "Querying vectorstore ...")

	// create a vector input
	embs := c.GetEmbeddings(ctx)
	input := &vectorstore.QueryInput{
		CustomerId: c.ID,
		Embeddings: embs,
		DB:         db,
		Query:      request.Query,
		K:          request.K,
		Logger:     logger,
	}

	// get the documents
	docs, err := vectorstore.QueryDocuments(ctx, &vectorstore.QueryDocstoreInput{QueryInput: input})
	if err != nil {
		return nil, fmt.Errorf("failed to query for documents: %s", err)
	}

	// report the usage
	if err := utils.ReportUsage(ctx, logger, db, c.ID, embs.GetUsageRecords(), nil); err != nil {
		logger.ErrorContext(ctx, "Failed to log vector usage: %s", err)
	}

	// get the website pages
	pages, err := vectorstore.QueryWebsitePages(ctx, &vectorstore.QueryWebsitePagesInput{QueryInput: input})
	if err != nil {
		return nil, fmt.Errorf("failed to query for website pages: %s", err)
	}

	// report the usage
	if err := utils.ReportUsage(ctx, logger, db, c.ID, embs.GetUsageRecords(), nil); err != nil {
		logger.ErrorContext(ctx, "Failed to log vector usage: %s", err)
	}

	return &queryVectorStoreResponse{
		Documents:    docs,
		WebsitePages: pages,
	}, nil
}
