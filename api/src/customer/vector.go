package customer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

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

	store, err := c.GetDocstore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docstore: %s", err)
	}

	// create a vector input
	embs := c.GetEmbeddings(ctx)
	input := &vectorstore.QueryInput{
		CustomerId: c.ID,
		Docstore:   store,
		Embeddings: embs,
		DB:         db,
		Query:      request.Query,
		K:          request.K,
		Logger:     logger,
	}

	// get the documents
	docs, err := vectorstore.QueryDocuments(ctx, input, request.IncludeContent)
	if err != nil {
		return nil, fmt.Errorf("failed to query for documents: %s", err)
	}

	// report the usage
	if err := utils.ReportUsage(ctx, logger, db, c.ID, embs.GetTokenRecords(), nil); err != nil {
		logger.ErrorContext(ctx, "Failed to log vector usage: %s", err)
	}

	// get the website pages
	pages, err := vectorstore.QueryWebsitePages(ctx, input, request.IncludeContent)
	if err != nil {
		return nil, fmt.Errorf("failed to query for website pages: %s", err)
	}

	// report the usage
	if err := utils.ReportUsage(ctx, logger, db, c.ID, embs.GetTokenRecords(), nil); err != nil {
		logger.ErrorContext(ctx, "Failed to log vector usage: %s", err)
	}

	return &queryVectorStoreResponse{
		Documents:    docs,
		WebsitePages: pages,
	}, nil

}
