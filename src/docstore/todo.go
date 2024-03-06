package docstore

import (
	"context"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

// faux implementation to use when storing the documents is not necessary
type TODODocstore struct {
	logger *slog.Logger
}

func NewTODODocstore(logger *slog.Logger) (*TODODocstore, error) {
	l := logger.With("docstore", "TODO")

	return &TODODocstore{
		logger: l,
	}, nil
}

func (d *TODODocstore) UploadDocument(ctx context.Context, customer *queries.Customer, doc *Doc) (string, error) {
	d.logger.InfoContext(ctx, "TODO -- UploadDocument")
	return "", nil
}

func (d *TODODocstore) GetDocument(ctx context.Context, customer *queries.Customer, filename string) (*Doc, error) {
	d.logger.InfoContext(ctx, "TODO -- GetDocument")
	return nil, nil
}

func (d *TODODocstore) DeleteDocument(ctx context.Context, customer *queries.Customer, filename string) error {
	d.logger.InfoContext(ctx, "TODO -- DeleteDocument")
	return nil
}
