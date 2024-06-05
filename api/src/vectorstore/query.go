package vectorstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/datastore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func Query(ctx context.Context, input *QueryAllInput) (*queries.QueryVectorStoreResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Logger.InfoContext(ctx, "Querying vector store for general retrieval query ...")

	// get the embeddings of the input
	vector, err := input.GetVectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	dmodel := queries.New(input.DB)
	response, err := dmodel.CUSTOMQueryVectorStore(ctx, &queries.QueryVectorStoreParams{
		CustomerID: input.CustomerId,
		Limit:      int32(input.K),
		Embeddings: &vector.Embedding,
		Column4:    input.DocumentIds,
		Column5:    input.FolderIds,
		Column6:    input.WebsitePageIds,
		Column7:    input.WebsiteIds,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query the vectorstore: %s", err)
	}

	return response, nil
}

// returns the raw vectors from the database
func QueryRaw(ctx context.Context, input *QueryInput) ([]*queries.VectorStore, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Logger.InfoContext(ctx, "Querying vector store for raw vector responses ...")

	// send the request
	vector, err := input.GetVectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(input.DB)
	vectors, err := model.QueryVectorStoreRaw(ctx, &queries.QueryVectorStoreRawParams{
		CustomerID: input.CustomerId,
		Limit:      int32(input.K),
		Embeddings: &vector.Embedding,
	})
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			input.Logger.InfoContext(ctx, "The result was empty")
			return []*queries.VectorStore{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store")
	}

	input.Logger.InfoContext(ctx, "Successfully got raw vectors")

	return vectors, nil
}

func QueryDocuments(
	ctx context.Context,
	input *QueryDocstoreInput,
) ([]*datastore.Document, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("the input was not valid: %s", err)
	}

	input.Logger.InfoContext(ctx, "Querying vector store for related documents ...")

	// send the request
	vector, err := input.GetVectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(input.DB)
	var rawDocs []*queries.Document
	if len(input.FolderIds) == 0 && len(input.DocumentIds) == 0 {
		rawDocs, err = model.QueryVectorStoreDocuments(ctx, &queries.QueryVectorStoreDocumentsParams{
			CustomerID: input.CustomerId,
			Limit:      int32(input.K),
			Embeddings: &vector.Embedding,
		})
	} else {
		// query scoped to the folders and/or documents
		rawDocs, err = model.QueryVectorStoreDocumentsScoped(ctx, &queries.QueryVectorStoreDocumentsScopedParams{
			CustomerID: input.CustomerId,
			Limit:      int32(input.K),
			Embeddings: &vector.Embedding,
			Column4:    input.FolderIds,
			Column5:    input.DocumentIds,
		})
	}
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			input.Logger.InfoContext(ctx, "The result was empty")
			return []*datastore.Document{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	// convert to internal type with content
	docs := make([]*datastore.Document, 0)
	docMap := make(map[string]bool, 0)
	for _, item := range rawDocs {
		// skip if doc already used
		if _, exists := docMap[item.ID.String()]; exists {
			continue
		}

		doc, err := datastore.NewDocumentFromDocument(ctx, input.Logger, item)
		if err != nil {
			return nil, fmt.Errorf("error creating document: %s", err)
		}

		docs = append(docs, doc)
		docMap[item.ID.String()] = true
	}

	input.Logger.InfoContext(ctx, "Successfully found documents", "length", len(docs))

	return docs, nil
}

func QueryWebsitePages(
	ctx context.Context,
	input *QueryWebsitePagesInput,
) ([]*datastore.WebsitePage, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.Logger.InfoContext(ctx, "Querying vector store for related website pages ...")

	// send the request
	vector, err := input.GetVectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(input.DB)
	var pagesRaw []*queries.WebsitePage
	if len(input.WebsiteIds) == 0 && len(input.WebsitePageIds) == 0 {
		pagesRaw, err = model.QueryVectorStoreWebsitePages(ctx, &queries.QueryVectorStoreWebsitePagesParams{
			CustomerID: input.CustomerId,
			Limit:      int32(input.K),
			Embeddings: &vector.Embedding,
		})
	} else {
		// scope the response to the website and/or pages
		pagesRaw, err = model.QueryVectorStoreWebsitePagesScoped(ctx, &queries.QueryVectorStoreWebsitePagesScopedParams{
			CustomerID: input.CustomerId,
			Limit:      int32(input.K),
			Embeddings: &vector.Embedding,
			Column4:    input.WebsiteIds,
			Column5:    input.WebsitePageIds,
		})
	}
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			input.Logger.InfoContext(ctx, "The result was empty")
			return []*datastore.WebsitePage{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	// query the website page for the content
	pages := make([]*datastore.WebsitePage, 0)
	webMap := make(map[string]bool, 0)
	for _, item := range pagesRaw {
		// skip if the website already has been used
		if _, exists := webMap[item.ID.String()]; exists {
			continue
		}

		page, err := datastore.NewWebsitePageFromWebsitePage(ctx, input.Logger, item)
		if err != nil {
			return nil, fmt.Errorf("error creating document: %s", err)
		}

		pages = append(pages, page)
		webMap[item.ID.String()] = true
	}

	input.Logger.InfoContext(ctx, "Successfully found website pages", "length", len(pages))

	return pages, nil
}
