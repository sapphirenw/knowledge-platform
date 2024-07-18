package datastore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/textsplitter"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Document struct {
	*queries.Document

	// cached data to reduce compute if needed
	raw      *bytes.Buffer // raw data
	metadata *bytes.Buffer // optional metadata
	cleaned  *bytes.Buffer // data but cleaned
	logger   *slog.Logger
}

func GetDocument(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	docId uuid.UUID,
) (*Document, error) {
	dmodel := queries.New(db)
	document, err := dmodel.GetDocument(ctx, docId)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %s", err)
	}
	return NewDocumentFromDocument(ctx, logger, document)
}

func NewDocumentFromData(
	ctx context.Context,
	logger *slog.Logger,
	customerId uuid.UUID,
	datastoreType string,
	filename string,
	data io.Reader,
) (*Document, error) {
	// parse the filetype
	filetype, err := ParseFileType(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the filetype: %s", err)
	}

	// create an id for the document
	docId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to create a uuid: %s", err)
	}

	// read the data
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read from the buffer: %s", err)
	}
	raw := buf.Bytes()

	return &Document{
		Document: &queries.Document{
			ID:            docId,
			CustomerID:    customerId,
			Filename:      filename,
			Type:          string(filetype),
			SizeBytes:     int64(len(raw)),
			Sha256:        utils.GenerateFingerprint(raw),
			Validated:     false,
			DatastoreType: datastoreType,
			DatastoreID:   fmt.Sprintf("%s/%s", customerId.String(), uuid.New().String()),
		},
		raw: buf,
	}, nil
}

func NewDocumentFromDocument(
	ctx context.Context,
	logger *slog.Logger,
	document *queries.Document,
) (*Document, error) {
	return &Document{Document: document, logger: logger}, nil
}

func (d *Document) GetDocstore(ctx context.Context) (docstore.RemoteDocstore, error) {
	switch d.DatastoreType {
	case "s3":
		return docstore.NewS3Docstore(ctx, docstore.S3_BUCKET, d.logger)
	default:
		return nil, fmt.Errorf("invalid docstore: %s", d.DatastoreType)
	}
}

func (d *Document) GetRaw(ctx context.Context) (*bytes.Buffer, error) {
	if d.raw == nil {
		// fetch the file from the remote datastore
		dstore, err := d.GetDocstore(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get docstore: %s", err)
		}

		// download the file
		raw, err := dstore.DownloadFile(ctx, d.DatastoreID)
		if err != nil {
			return nil, fmt.Errorf("failed to download the file: %s", err)
		}

		// create the buffer
		buf := bytes.NewBuffer(raw)
		if buf == nil {
			return nil, fmt.Errorf("failed to create the buffer")
		}
		d.raw = buf

	}
	return d.raw, nil
}

func (d *Document) GetCleaned(ctx context.Context) (*bytes.Buffer, error) {
	if d.cleaned != nil {
		return d.cleaned, nil
	}

	raw, err := d.GetRaw(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the raw data: %s", err)
	}

	filetype, err := ParseFileType(d.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the filetype: %s", err)
	}

	// clean the contents
	var content string

	// parse the contents based on the filetype
	switch filetype {
	case FT_html:
		content, err = ParseHTML(raw.Bytes())
	default:
		// use an auto-content detection parser
		content, err = ParseDynamic(raw.Bytes(), filetype)
		// TODO -- handle errors. May potentially need to just use the raw content here
	}

	if err != nil {
		return nil, fmt.Errorf("there was an issue parsing the document: %v", err)
	}

	// clean the string
	cleaned := utils.CleanInput(content)

	// create a new buffer with this cleaned content
	buf := new(bytes.Buffer)
	_, err = buf.WriteString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("failed to write to the buffer: %s", err)
	}
	d.cleaned = buf

	return d.cleaned, nil
}

func (d *Document) GetChunks(ctx context.Context) ([]string, error) {

	filetype, err := ParseFileType(d.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the filetype: %s", err)
	}

	content, err := d.GetCleaned(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the cleaned content")
	}

	// chunk the content based on what type of file
	var chunks []string
	switch filetype {
	case FT_html:
	case FT_md:
		splitter := textsplitter.NewMarkdownTextSplitter(
			textsplitter.WithChunkSize(2000),
			textsplitter.WithChunkOverlap(200),
		)
		chunks, err = splitter.SplitText(content.String())
	default:
		splitter := textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(2000),
			textsplitter.WithChunkOverlap(200),
		)
		chunks, err = splitter.SplitText(content.String())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to split the text")
	}

	return chunks, nil
}

// TODO -- implement
func (d *Document) GetMetadata(ctx context.Context) (*bytes.Buffer, error) {
	return new(bytes.Buffer), nil
}

func (d *Document) GetSha256() (string, error) {
	return d.Sha256, nil
}

func (d *Document) getSummary() string {
	if d.Summary == "" || d.Sha256 != d.SummarySha256 {
		return ""
	}
	return d.Summary
}

func (d *Document) setSummary(s string) error {
	d.Summary = s
	d.SummarySha256 = d.Sha256
	return nil
}
