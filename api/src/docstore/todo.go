package docstore

import (
	"context"
	"log/slog"
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

func (d *TODODocstore) GeneratePresignedUrl(ctx context.Context, doc *Document) (string, error) {
	return "", nil
}

func (d *TODODocstore) DownloadFile(ctx context.Context, uniqueId string) ([]byte, error) {
	d.logger.InfoContext(ctx, "TODO -- GetDocument")
	return nil, nil
}

func (d *TODODocstore) DeleteFile(ctx context.Context, uniqueId string) error {
	d.logger.InfoContext(ctx, "TODO -- DeleteDocument")
	return nil
}

func (d *TODODocstore) DeleteRoot(ctx context.Context, prefix string) error {
	d.logger.InfoContext(ctx, "TODO -- DeleteRoot")
	return nil
}

func (d *TODODocstore) GetUploadMethod() string {
	return ""
}
