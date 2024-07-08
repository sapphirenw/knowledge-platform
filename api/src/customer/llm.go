package customer

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

func getAvailableLLMs(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "getAvailableLLMs")

	dmodel := queries.New(pool)

	response, err := dmodel.GetLLMsByCustomerAvailable(r.Context(), utils.GoogleUUIDToPGXUUID(c.ID))
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to query the database", err)
		return
	}

	request.Encode(w, r, logger, 200, response)
}

func updateCustomerLLMConfigurations(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "updateCustomerLLMConfigurations")
	// read the request body
	// parse the request
	body, valid := request.Decode[updateCustomerLLMConfigsRequest](w, r, c.logger)
	if !valid {
		return
	}

	// get a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start a transaction", err)
		return
	}
	defer tx.Commit(r.Context())
	dmodel := queries.New(tx)

	// update based on body
	if body.SummaryLLMID != "" {
		if _, err := dmodel.UpdateCustomerSummaryLLM(r.Context(), &queries.UpdateCustomerSummaryLLMParams{
			CustomerID:   c.ID,
			SummaryLlmID: utils.PGXUUIDFromString(body.SummaryLLMID),
		}); err != nil {
			slogger.ServerError(w, logger, 500, "failed to update the summary llm", err)
			return
		}
	}
	if body.ChatLLMID != "" {
		if _, err := dmodel.UpdateCustomerChatLLM(r.Context(), &queries.UpdateCustomerChatLLMParams{
			CustomerID: c.ID,
			ChatLlmID:  utils.PGXUUIDFromString(body.ChatLLMID),
		}); err != nil {
			slogger.ServerError(w, logger, 500, "failed to update the chat llm", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Customer) GetSummaryLLM(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
) (*llm.LLM, error) {
	dmodel := queries.New(db)

	response, err := dmodel.GetCustomerSummaryLLM(ctx, c.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the summary llm", err)
	}

	return llm.FromObjects(&response.Llm, &response.AvailableModel), nil
}

func (c *Customer) GetChatLLM(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
) (*llm.LLM, error) {
	dmodel := queries.New(db)

	response, err := dmodel.GetCustomerChatLLM(ctx, c.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the summary llm", err)
	}

	return llm.FromObjects(&response.Llm, &response.AvailableModel), nil
}
