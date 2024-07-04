package vectorstore

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

// Query all objects in the datastore. Respects all passed filters
func Query(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	input *QueryInput,
) (*QueryResponse, error) {
	if logger == nil {
		logger = slog.Default()
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	logger.InfoContext(ctx, "Querying vector store for general retrieval query ...")

	// get the embeddings of the input
	_, err := input.GetVectors(ctx, logger)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get vectors", err)
	}

	// get documents
	docResponse, err := QueryDocuments(ctx, logger, db, input)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the documents", err)
	}
	// get website pages
	pageResponse, err := QueryWebsitePages(ctx, logger, db, input)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the website pages", err)
	}

	// combine vectors
	vectors := make([]*queries.VectorStore, 0)
	vectors = append(vectors, docResponse.Vectors...)
	vectors = append(vectors, pageResponse.Vectors...)

	return &QueryResponse{
		Vectors:      vectors,
		Documents:    docResponse.Documents,
		WebsitePages: pageResponse.WebsitePages,
	}, nil
}

// returns the raw vectors from the database
func QueryRaw(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	input *QueryInput,
) ([]*queries.VectorStore, error) {
	if logger == nil {
		logger = slog.Default()
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	logger.InfoContext(ctx, "Querying vector store for raw vector responses ...")

	// send the request
	vector, err := input.GetVectors(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(db)
	vectors, err := model.QueryVectorStoreRaw(ctx, &queries.QueryVectorStoreRawParams{
		CustomerID: input.CustomerID,
		Limit:      int32(input.K),
		Embeddings: &vector.Embedding,
	})
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			logger.InfoContext(ctx, "The result was empty")
			return []*queries.VectorStore{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store")
	}

	logger.InfoContext(ctx, "Successfully got raw vectors")

	return vectors, nil
}

// Query documents in the datastore.
// Respects the document and folder filters
func QueryDocuments(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	input *QueryInput,
) (*QueryDocumentsResponse, error) {
	if logger == nil {
		logger = slog.Default()
	}
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("the input was not valid: %s", err)
	}

	logger.InfoContext(ctx, "Querying vector store for related documents ...")

	// send the request
	vector, err := input.GetVectors(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(db)
	response, err := model.QueryVectorStoreDocumentsScoped(ctx, &queries.QueryVectorStoreDocumentsScopedParams{
		CustomerID: input.CustomerID,
		Limit:      int32(input.K),
		Embeddings: &vector.Embedding,
		Column4:    input.DocumentIDsFilter,
		Column5:    input.FolderIDsFilter,
	})
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			logger.InfoContext(ctx, "The result was empty")
			return &QueryDocumentsResponse{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	logger.InfoContext(ctx, "Successfully found documents", "length", len(response))

	// convert to format
	vectors := make([]*queries.VectorStore, 0)
	docs := make([]*queries.Document, 0)

	for _, item := range response {
		vectors = append(vectors, &item.VectorStore)
		docs = append(docs, &item.Document)
	}

	return &QueryDocumentsResponse{
		Vectors:   vectors,
		Documents: docs,
	}, nil
}

// query website pages in the datastore
// Respects the page and website filters
func QueryWebsitePages(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	input *QueryInput,
) (*QueryWebsitePagesResponse, error) {
	if logger == nil {
		logger = slog.Default()
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}
	logger.InfoContext(ctx, "Querying vector store for related website pages ...")

	// send the request
	vector, err := input.GetVectors(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(db)
	// response, err := model.QueryVectorStoreWebsitePagesScoped(ctx, &queries.QueryVectorStoreWebsitePagesScopedParams{
	// 	CustomerID: input.CustomerID,
	// 	Limit:      int32(input.K),
	// 	Embeddings: &vector.Embedding,
	// 	Column4:    input.WebsitePageIDsFilter,
	// 	Column5:    input.WebsiteIDsFilter,
	// })
	response, err := model.QueryVectorStoreWebsitePages(ctx, &queries.QueryVectorStoreWebsitePagesParams{
		CustomerID: input.CustomerID,
		Limit:      int32(input.K),
		Embeddings: &vector.Embedding,
	})
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			logger.InfoContext(ctx, "The result was empty")
			return &QueryWebsitePagesResponse{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	logger.InfoContext(ctx, "Successfully found website pages", "length", len(response))

	// convert to format
	vectors := make([]*queries.VectorStore, 0)
	pages := make([]*queries.WebsitePage, 0)

	for _, item := range response {
		vectors = append(vectors, &item.VectorStore)
		pages = append(pages, &item.WebsitePage)
	}

	return &QueryWebsitePagesResponse{
		Vectors:      vectors,
		WebsitePages: pages,
	}, nil
}
