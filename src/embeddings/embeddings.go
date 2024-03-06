package embeddings

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type Embeddings interface {
	// Create the embdeddings using the provider
	Embed(ctx context.Context, input string) ([]*EmbeddingsData, error)

	// report the usage, if any, consumed by the request by the user
	ReportUsage(ctx context.Context, txn pgx.Tx, customer *queries.Customer) error
}

type EmbeddingsData struct {
	Raw       string
	Embedding pgvector.Vector
}
