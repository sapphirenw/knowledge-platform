package vectorstore

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/embeddings"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type QueryInput struct {
	CustomerId uuid.UUID
	Docstore   docstore.RemoteDocstore
	Embeddings embeddings.Embeddings
	DB         queries.DBTX
	Query      string
	K          int
	Logger     *slog.Logger

	// can inbed the vector incase the input is re-used, or user already embedded content
	Vector *embeddings.EmbeddingsData
}

func (input *QueryInput) Validate() error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}
	if input.Docstore == nil {
		return fmt.Errorf("docstore cannot be nil")
	}
	if input.Embeddings == nil {
		return fmt.Errorf("embeddings cannot be empty")
	}
	if input.DB == nil {
		return fmt.Errorf("db cannot be empty")
	}
	if input.Query == "" {
		return fmt.Errorf("no query provided")
	}
	if input.K == 0 || input.K > 5 {
		return fmt.Errorf("k must be between 1 and 5")
	}
	if input.Logger == nil {
		input.Logger = utils.DefaultLogger()
	}

	// ensure the length is not too large for the embeddings
	// TODO -- support other embeddings
	if len(input.Query) > embeddings.OPENAI_EMBEDDINGS_INPUT_MAX {
		return fmt.Errorf("the query is too long: %d characters", len(input.Query))
	}

	return nil
}

func (input *QueryInput) GetVectors(ctx context.Context) (*embeddings.EmbeddingsData, error) {
	var vector *embeddings.EmbeddingsData
	if input.Vector == nil {
		input.Logger.InfoContext(ctx, "Vectors not present, creating new vectors from the input")
		response, err := input.Embeddings.Embed(ctx, input.Query)
		if err != nil {
			return nil, fmt.Errorf("error sending the embedding request: %s", err)
		}
		if len(response) == 0 {
			return nil, fmt.Errorf("there were no embeddings returned")
		}

		input.Vector = response[0]
		vector = response[0]
	} else {
		vector = input.Vector
	}

	return vector, nil
}
