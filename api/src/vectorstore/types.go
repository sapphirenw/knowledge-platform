package vectorstore

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/textsplitter"
)

type QueryInput struct {
	CustomerID uuid.UUID
	Embeddings gollm.Embeddings
	Query      string
	K          int

	// filters

	FolderIDsFilter   []uuid.UUID
	DocumentIDsFilter []uuid.UUID

	WebsiteIDsFilter     []uuid.UUID
	WebsitePageIDsFilter []uuid.UUID

	// can inbed the vector incase the input is re-used, or user already embedded content
	Vector *ltypes.EmbeddingsData
}

func (input *QueryInput) Validate() error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}
	if input.Embeddings == nil {
		return fmt.Errorf("embeddings cannot be empty")
	}
	if input.Query == "" {
		return fmt.Errorf("no query provided")
	}
	if input.K == 0 || input.K > 5 {
		return fmt.Errorf("k must be between 1 and 5")
	}

	// ensure the length is not too large for the embeddings
	// TODO -- support other embeddings
	if len(input.Query) > gollm.OPENAI_EMBEDDINGS_INPUT_MAX {
		return fmt.Errorf("the query is too long: %d characters", len(input.Query))
	}

	return nil
}

func (input *QueryInput) GetVectors(ctx context.Context, logger *slog.Logger) (*ltypes.EmbeddingsData, error) {
	if input.Vector == nil {
		logger.InfoContext(ctx, "Vectors not present, creating new vectors from the input")

		// get a text splitter
		splitter := textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(8191),
			textsplitter.WithChunkOverlap(0),
		)

		response, err := input.Embeddings.Embed(ctx, logger, &gollm.EmbedArgs{
			Input:            input.Query,
			ChunkingFunction: splitter.SplitText,
		})
		if err != nil {
			return nil, fmt.Errorf("error sending the embedding request: %s", err)
		}
		if len(response.Embeddings) == 0 {
			return nil, fmt.Errorf("there were no embeddings returned")
		}

		input.Vector = response.Embeddings[0]
	}

	return input.Vector, nil
}

type QueryResponse struct {
	Vectors      []*queries.VectorStore
	Documents    []*queries.Document
	WebsitePages []*queries.WebsitePage
}

type QueryDocumentsResponse struct {
	Vectors   []*queries.VectorStore
	Documents []*queries.Document
}

type QueryWebsitePagesResponse struct {
	Vectors      []*queries.VectorStore
	WebsitePages []*queries.WebsitePage
}
