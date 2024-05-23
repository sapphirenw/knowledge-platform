package docstore

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Document struct {
	*queries.Document

	// created on initialization
	Filetype Filetype `json:"filetype"`
	UniqueID string   `json:"uniqueId"`

	// internal field used to cache remote request results
	data    []byte
	cleaned string
}

func NewDocument(ctx context.Context, db queries.DBTX, docId uuid.UUID) (*Document, error) {
	model := queries.New(db)
	document, err := model.GetDocument(ctx, docId)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %s", err)
	}
	return NewDocumentFromDocument(document)
}

func NewDocumentFromDocument(doc *queries.Document) (*Document, error) {
	if doc == nil {
		return nil, fmt.Errorf("doc cannot be nil")
	}

	// parse the filetype
	filetype, err := ParseFileType(doc.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the filetype: %s", err)
	}

	// create a uniqueID for the document
	id := fileIdFromDoc(doc)

	return &Document{Document: doc, Filetype: filetype, UniqueID: id}, nil
}

// Creates a document from raw file information. Used mostly in tests
func NewDocumentFromRaw(customer *queries.Customer, filename string, data []byte) (*Document, error) {
	// parse the filetype
	filetype, err := ParseFileType(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the filetype: %s", err)
	}

	// create a uniqueID for the document
	documentId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create a uuid: %s", err)
	}
	id := fileIdFromRaw(customer.ID, documentId, filename)

	return &Document{
		Document: &queries.Document{
			ID:         documentId,
			CustomerID: customer.ID,
			Filename:   filename,
			Type:       string(filetype),
			SizeBytes:  int64(len(data)),
			Sha256:     utils.GenerateFingerprint(data),
			Validated:  false,
		},
		Filetype: filetype,
		UniqueID: id,
		data:     data,
	}, nil
}

func (d *Document) GetRawData(ctx context.Context, store RemoteDocstore) ([]byte, error) {
	if d.data != nil {
		return d.data, nil
	}

	// donwload the file using the provided document store
	data, err := store.DownloadFile(ctx, d.UniqueID)
	if err != nil {
		return nil, fmt.Errorf("there was an issue downloading the data: %s", err)
	}
	d.data = data

	return d.data, nil
}

func (d *Document) GetCleanedContents(ctx context.Context, store RemoteDocstore) (string, error) {
	if d.cleaned != "" {
		return d.cleaned, nil
	}

	if d.data == nil {
		// get the raw data
		if _, err := d.GetRawData(ctx, store); err != nil {
			return "", fmt.Errorf("there was an error getting the raw content: %s", err)
		}
	}

	// clean the contents
	var content string
	var err error

	// parse the contents based on the filetype
	switch d.Filetype {
	case FT_html:
		content, err = ParseHTML(d.data)
	default:
		// use an auto-content detection parser
		content, err = ParseDynamic(d.data, string(d.Filetype))
		// TODO -- handle errors. May potentially need to just use the raw content here
	}

	if err != nil {
		return "", fmt.Errorf("there was an issue parsing the document: %v", err)
	}

	// clean the string
	d.cleaned = utils.CleanInput(content)

	return d.cleaned, nil
}
