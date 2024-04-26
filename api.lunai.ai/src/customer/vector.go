package customer

import (
	"context"
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

func (c *Customer) QueryVectorStore(ctx context.Context, db queries.DBTX, request *queryVectorStoreRequest) (*queryVectorStoreResponse, error) {
	logger := c.logger.With("request", request)
	logger.InfoContext(ctx, "Querying vectorstore ...")

	store, err := c.GetDocstore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docstore: %s", err)
	}

	// create a vector input
	input := &vectorstore.QueryInput{
		CustomerId: c.ID,
		Docstore:   store,
		Embeddings: c.GetEmbeddings(ctx),
		DB:         db,
		Query:      request.Query,
		K:          request.K,
		Logger:     logger,
	}

	// get the documents
	docs, err := vectorstore.QueryDocuments(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query for documents: %s", err)
	}

	// get the website pages
	pages, err := vectorstore.QueryWebsitePages(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query for website pages: %s", err)
	}

	return &queryVectorStoreResponse{
		Documents:    docs,
		WebsitePages: pages,
	}, nil

}
