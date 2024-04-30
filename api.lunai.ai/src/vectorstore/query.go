package vectorstore

import (
	"context"
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/webscrape"
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
		return nil, fmt.Errorf("error querying the vector store")
	}

	input.Logger.InfoContext(ctx, "Successfully got raw vectors")

	return vectors, nil
}

func QueryDocuments(ctx context.Context, input *QueryInput, include bool) ([]*DocumentResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	input.Logger.InfoContext(ctx, "Querying vector store for related documents ...")

	// send the request
	vector, err := input.GetVectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %s", err)
	}

	// send the request to the database
	model := queries.New(input.DB)
	rawDocs, err := model.QueryVectorStoreDocuments(ctx, &queries.QueryVectorStoreDocumentsParams{
		CustomerID: input.CustomerId,
		Limit:      int32(input.K),
		Embeddings: vector.Embedding,
	})
	if err != nil {
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	// convert to internal type with content
	docs := make([]*DocumentResponse, 0)
	docMap := make(map[int64]bool, 0)
	for _, item := range rawDocs {
		// skip if doc already used
		if _, exists := docMap[item.ID]; exists {
			continue
		}

		doc, err := docstore.NewDocument(input.CustomerId, item)
		if err != nil {
			return nil, fmt.Errorf("error creating document: %s", err)
		}

		var content string

		if include {
			content, err = doc.GetCleanedContents(ctx, input.Docstore)
			if err != nil {
				return nil, fmt.Errorf("error getting the cleaned contents: %s", err)
			}
		}

		docs = append(docs, &DocumentResponse{Document: doc, Content: content})
		docMap[item.ID] = true
	}

	input.Logger.InfoContext(ctx, "Successfully found documents", "length", len(docs))

	return docs, nil
}

func QueryWebsitePages(ctx context.Context, input *QueryInput, include bool) ([]*WebsitePageResonse, error) {
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
	pagesRaw, err := model.QueryVectorStoreWebsitePages(ctx, &queries.QueryVectorStoreWebsitePagesParams{
		CustomerID: input.CustomerId,
		Limit:      int32(input.K),
		Embeddings: vector.Embedding,
	})
	if err != nil {
		return nil, fmt.Errorf("error querying the vector store: %s", err)
	}

	// query the website page for the content
	pages := make([]*WebsitePageResonse, 0)
	webMap := make(map[int64]bool, 0)
	for _, item := range pagesRaw {
		// skip if the website already has been used
		if _, exists := webMap[item.ID]; exists {
			continue
		}

		var content []byte
		if include {
			content, err = webscrape.ScrapeSingle(ctx, input.Logger, item)
			if err != nil {
				return nil, fmt.Errorf("failed to scrape the website: %s", err)
			}
		}

		pages = append(pages, &WebsitePageResonse{
			WebsitePage: item,
			Content:     string(content),
		})
		webMap[item.ID] = true
	}

	input.Logger.InfoContext(ctx, "Successfully found website pages", "length", len(pages))

	return pages, nil
}
