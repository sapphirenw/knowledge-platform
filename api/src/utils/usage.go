package utils

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

// reports the llm usage on the customer
func ReportUsage(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	customerId uuid.UUID,
	records []*tokens.UsageRecord,
	conversation *queries.Conversation, // optional conversation to tie the usage to
) error {
	logger.InfoContext(ctx, "Reporting usage", "records", len(records))
	model := queries.New(db)

	// insert all internal token records
	for idx, item := range records {
		logger.InfoContext(ctx, "Posting to database ...", "index", idx)
		var convId pgtype.UUID
		if conversation != nil {
			convId = GoogleUUIDToPGXUUID(conversation.ID)
		}
		_, err := model.CreateTokenUsage(ctx, &queries.CreateTokenUsageParams{
			ID:             item.ID,
			CustomerID:     customerId,
			ConversationID: convId,
			Model:          item.Model,
			InputTokens:    int32(item.InputTokens),
			OutputTokens:   int32(item.OutputTokens),
			TotalTokens:    int32(item.TotalTokens),
		})
		if err != nil {
			return err
		}
		logger.InfoContext(ctx, "Done.", "index", idx)
	}

	logger.InfoContext(ctx, "Successfully reported usage")

	return nil
}
