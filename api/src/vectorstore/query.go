package vectorstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/webparse"
)

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
		Embeddings: vector.Embedding,
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
) ([]*DocumentResponse, error) {
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
			Embeddings: vector.Embedding,
		})
	} else {
		// query scoped to the folders and/or documents
		rawDocs, err = model.QueryVectorStoreDocumentsScoped(ctx, &queries.QueryVectorStoreDocumentsScopedParams{
			CustomerID: input.CustomerId,
			Limit:      int32(input.K),
			Embeddings: vector.Embedding,
			Column4:    input.FolderIds,
			Column5:    input.DocumentIds,
		})
	}
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			input.Logger.InfoContext(ctx, "The result was empty")
			return []*DocumentResponse{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	// convert to internal type with content
	docs := make([]*DocumentResponse, 0)
	docMap := make(map[string]bool, 0)
	for _, item := range rawDocs {
		// skip if doc already used
		if _, exists := docMap[item.ID.String()]; exists {
			continue
		}

		doc, err := docstore.NewDocumentFromDocument(item)
		if err != nil {
			return nil, fmt.Errorf("error creating document: %s", err)
		}

		var content string

		if input.IncludeContent {
			content, err = doc.GetCleanedContents(ctx, input.Docstore)
			if err != nil {
				return nil, fmt.Errorf("error getting the cleaned contents: %s", err)
			}
		}

		docs = append(docs, &DocumentResponse{Document: doc, Content: content})
		docMap[item.ID.String()] = true
	}

	input.Logger.InfoContext(ctx, "Successfully found documents", "length", len(docs))

	return docs, nil
}

func QueryWebsitePages(
	ctx context.Context,
	input *QueryWebsitePagesInput,
) ([]*WebsitePageResponse, error) {
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
			Embeddings: vector.Embedding,
		})
	} else {
		// scope the response to the website and/or pages
		pagesRaw, err = model.QueryVectorStoreWebsitePagesScoped(ctx, &queries.QueryVectorStoreWebsitePagesScopedParams{
			CustomerID: input.CustomerId,
			Limit:      int32(input.K),
			Embeddings: vector.Embedding,
			Column4:    input.WebsiteIds,
			Column5:    input.WebsitePageIds,
		})
	}
	if err != nil {
		if strings.Contains(err.Error(), "db cannot be empty") {
			input.Logger.InfoContext(ctx, "The result was empty")
			return []*WebsitePageResponse{}, nil
		}
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	// query the website page for the content
	pages := make([]*WebsitePageResponse, 0)
	webMap := make(map[string]bool, 0)
	for _, item := range pagesRaw {
		// skip if the website already has been used
		if _, exists := webMap[item.ID.String()]; exists {
			continue
		}

		var content string
		if input.IncludeContent {
			response, err := webparse.ScrapeSingle(ctx, input.Logger, item)
			if err != nil {
				return nil, fmt.Errorf("failed to scrape the website: %s", err)
			}
			content = response.Content
		}

		pages = append(pages, &WebsitePageResponse{
			WebsitePage: item,
			Content:     content,
		})
		webMap[item.ID.String()] = true
	}

	input.Logger.InfoContext(ctx, "Successfully found website pages", "length", len(pages))

	return pages, nil
}
