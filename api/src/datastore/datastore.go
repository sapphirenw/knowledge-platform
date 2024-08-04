package datastore

import (
	"bytes"
	"context"
)

type Object interface {
	GetRaw(ctx context.Context) (*bytes.Buffer, error)
	GetCleaned(ctx context.Context) (*bytes.Buffer, error)
	GetChunks(ctx context.Context) ([]string, error)
	GetMetadata(ctx context.Context) (*bytes.Buffer, error)
	GetSha256() (string, error)

	getSummary() string
	setSummary(s string) error
}

// Gets the summary object from the object or generates it
// if it does not exist or the content is out of date.
// This does NOT write the summary to the database, updates
// will have to be handled manually for the sake of making this function
// thread-safe.
// func GetSummary(
// 	obj Object,
// 	ctx context.Context,
// 	logger *slog.Logger,
// 	customerId uuid.UUID,
// 	model *llm.LLM,
// ) (*llm.SummarizeResponse, error) {
// 	summary := obj.getSummary()

// 	if summary == "" {
// 		logger.InfoContext(ctx, "Generating new summary for object ...")

// 		// get the cleaned data
// 		cleaned, err := obj.GetCleaned(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get the cleaned data: %w", err)
// 		}

// 		// generate a new summary
// 		s, err := model.Summarize(ctx, logger, customerId, cleaned.String())
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create the summary: %w", err)
// 		}
// 		if err := obj.setSummary(s); err != nil {
// 			return nil, fmt.Errorf("failed to set the summary: %w", err)
// 		}
// 		summary = s
// 	}

// 	return summary, nil
// }
