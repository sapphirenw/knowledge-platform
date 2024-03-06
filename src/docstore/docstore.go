package docstore

import (
	"context"

	"github.com/sapphirenw/ai-content-creation-api/src/document"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type Docstore interface {
	// Uploads a document and returns the url of the document
	UploadDocument(ctx context.Context, customer *queries.Customer, doc *document.Doc) (string, error)
	GetDocument(ctx context.Context, customer *queries.Customer, filename string) (*document.Doc, error)
	DeleteDocument(ctx context.Context, customer *queries.Customer, filename string) error
}
