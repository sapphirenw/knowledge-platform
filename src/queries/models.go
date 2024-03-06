// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package queries

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type Customer struct {
	ID        int64
	Name      string
	Datastore string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Document struct {
	ID         int64
	ParentID   int64
	CustomerID int64
	Filename   string
	Type       string
	SizeBytes  int64
	Sha256     string
	CreatedAt  pgtype.Timestamp
}

type Folder struct {
	ID         int64
	ParentID   int64
	CustomerID int64
	Title      string
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

type TokenUsage struct {
	ID           pgtype.UUID
	CustomerID   int64
	Model        string
	InputTokens  int32
	OutputTokens int32
	TotalTokens  int32
	CreatedAt    pgtype.Timestamp
}

type VectorStore struct {
	ID         int64
	Raw        string
	Embeddings pgvector.Vector
	CustomerID int64
	DocumentID int64
	Index      int32
	CreatedAt  pgtype.Timestamp
}
