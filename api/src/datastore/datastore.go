package datastore

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
)

type Object interface {
	GetRaw(ctx context.Context) (*bytes.Buffer, error)
	GetCleaned(ctx context.Context) (*bytes.Buffer, error)
	GetSha256() (string, error)

	getSummary() string
	setSummary(s string) error
}

// Gets the summary object from the object or generates it
// if it does not exist or the content is out of date.
// This does NOT write the summary to the database, updates
// will have to be handled manually for the sake of making this function
// thread-safe.
func GetSummary(
	obj Object,
	ctx context.Context,
	logger *slog.Logger,
	customerId uuid.UUID,
	tokenRecords chan *tokens.TokenRecord,
	model *llm.LLM,
) (string, error) {
	summary := obj.getSummary()

	if summary == "" {
		logger.InfoContext(ctx, "Generating new summary for object ...")

		// get the cleaned data
		cleaned, err := obj.GetCleaned(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to get the cleaned data: %s", err)
		}

		// generate a new summary
		s, err := model.Summarize(ctx, logger, customerId, tokenRecords, cleaned.String())
		if err != nil {
			return "", fmt.Errorf("failed to create the summary: %s", err)
		}
		if err := obj.setSummary(s); err != nil {
			return "", fmt.Errorf("failed to set the summary: %s", err)
		}
		summary = s
	}

	return summary, nil
}
