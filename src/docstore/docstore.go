package docstore

import (
	"context"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type Filetype string

const (
	FT_none = ""

	FT_txt  = "txt"
	FT_md   = "md"
	FT_html = "html"
	FT_xml  = "xml"
	FT_csv  = "csv"
	FT_tsv  = "tsv"
	FT_pdf  = "pdf"

	FT_unknown = "unknown"
)

type Docstore interface {
	// Uploads a document and returns the url of the document
	UploadDocument(ctx context.Context, customer *queries.Customer, doc *Doc) (string, error)
	GetDocument(ctx context.Context, customer *queries.Customer, filename string) (*Doc, error)
	DeleteDocument(ctx context.Context, customer *queries.Customer, filename string) error
}

type Doc struct {
	Filename string
	Filetype Filetype
	Data     []byte
}

func NewDoc(filename string, data []byte) (*Doc, error) {
	// parse the filetype
	filetype, err := parseFileType(filename)
	if err != nil {
		return nil, err
	}

	return &Doc{
		Filename: filename,
		Filetype: filetype,
		Data:     data,
	}, nil
}
