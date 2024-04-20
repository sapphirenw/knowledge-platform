package docstore

import (
	"context"

	"github.com/sapphirenw/ai-content-creation-api/src/document"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type Docstore interface {
	// Uploads a document and returns the url of the document
	UploadDocument(ctx context.Context, customer *queries.Customer, doc *document.Doc) (string, error)
	GetDocument(ctx context.Context, customer *queries.Customer, parentId *int64, filename string) (*document.Doc, error)
	DeleteDocument(ctx context.Context, customer *queries.Customer, parentId *int64, filename string) error
	DeleteRoot(ctx context.Context, customer *queries.Customer) error

	// Requests a pre-signed url a client can use to upload documents
	GeneratePresignedUrl(ctx context.Context, customer *queries.Customer, input *UploadUrlInput) (string, error)

	// returns the method the client should use for the pre-signed url
	GetUploadMethod() string
}

type UploadUrlInput struct {
	ParentId  *int64 `json:"parentId"`
	Filename  string `json:"filename"`
	Mime      string `json:"mime"`
	Signature string `json:"signature"`
	Size      int64  `json:"size"`
}
