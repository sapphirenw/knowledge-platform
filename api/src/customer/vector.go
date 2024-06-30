package customer

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
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
	embs := llm.GetEmbeddings(logger, c.Customer)
	input := &vectorstore.QueryInput{
		CustomerID: c.ID,
		Embeddings: embs,
		Query:      request.Query,
		K:          request.K,
	}

	// run the general response
	response, err := vectorstore.Query(ctx, logger, db, input)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to query the vectorstore", err)
	}

	return &queryVectorStoreResponse{
		Documents:    response.Documents,
		WebsitePages: response.WebsitePages,
	}, nil
}

func queryVectorStoreDocuments(
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
		slogger.ServerError(w, r, c.logger, 500, "failed to start transaction", "error", err)
		return
	}
	defer tx.Commit(r.Context())

	embs := llm.GetEmbeddings(c.logger, c.Customer)
	response, err := vectorstore.QueryDocuments(r.Context(), c.logger, pool, &vectorstore.QueryInput{
		CustomerID: c.ID,
		Embeddings: embs,
		Query:      body.Query,
		K:          body.K,
	})
	if err != nil {
		slogger.ServerError(w, r, c.logger, 500, "failed to query the vectorstore", "error", err)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}
