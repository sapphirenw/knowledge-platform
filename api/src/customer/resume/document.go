package resume

import (
	"context"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

func (c *Client) AttachDocument(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	documentId string,
	isResume bool,
) error {
	logger := l.With("desc", "attach document to resume")

	// parse the id
	id, err := utils.GoogleUUIDFromString(documentId)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to parse the document id", err)
	}

	// ensure the document exists
	dmodel := queries.New(db)
	if _, err := dmodel.GetDocument(ctx, id); err != nil {
		return slogger.Error(ctx, logger, "the document does not exist", err)
	}

	// create the relationship
	if _, err := dmodel.CreateResumeDocument(ctx, &queries.CreateResumeDocumentParams{
		ResumeID:   c.Resume.ID,
		DocumentID: id,
		IsResume:   isResume,
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to attach the document to the resume", err)
	}

	logger.Info("successfully added document to the resume")

	return nil
}

func (c *Client) GetDocuments(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) ([]*queries.Document, error) {
	logger := l.With("desc", "get documents of resume")
	dmodel := queries.New(db)

	docs, err := dmodel.GetResumeDocuments(ctx, c.Resume.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the documents", err)
	}

	return docs, nil
}

func (c *Client) GetResume(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) (*queries.Document, error) {
	logger := l.With("desc", "get resume of resume")
	dmodel := queries.New(db)

	doc, err := dmodel.GetResumeResume(ctx, c.Resume.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the resume", err)
	}

	return doc, nil
}
