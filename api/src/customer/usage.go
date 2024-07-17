package customer

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func getUsage(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "getUsage")

	// get the page
	var page int
	var err error
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
	}

	if err != nil {
		slogger.ServerError(w, logger, 400, "failed to parse the page", err)
		return
	}

	// get the provider (optional)
	model := r.URL.Query().Get("model")

	// get the batch size (optional)
	batchSize, err := strconv.Atoi(r.URL.Query().Get("batchSize"))
	if err != nil || batchSize > 20 || batchSize < 0 {
		batchSize = 20
	}

	// query the database
	dmodel := queries.New(pool)
	records, err := dmodel.GetCustomerTokenUsage(r.Context(), &queries.GetCustomerTokenUsageParams{
		CustomerID: c.ID,
		Model:      model,
		Limit:      int32(batchSize),
		Column4:    int32(page),
	})
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to run the query", err)
		return
	}

	// get the number of pages
	pageCount, err := dmodel.GetCustomerTokensPageCount(r.Context(), &queries.GetCustomerTokensPageCountParams{
		CustomerID: c.ID,
		Model:      model,
		Column3:    int32(batchSize),
	})
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get number of pages", err)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, map[string]any{
		"metadata": map[string]any{
			"pageCount": int(pageCount),
		},
		"records": records,
	})
}

func getUsageGrouped(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "getUsageGrouped")

	dmodel := queries.New(pool)
	response, err := dmodel.GetCustomerUsageGrouped(r.Context(), c.ID)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to query the db", err)
		return
	}

	// compose some caluclations based on the retrieved data
	list := make([]map[string]any, 0)
	for _, item := range response {
		inc, _ := item.InputCostPerMillionTokens.Float64Value()
		outc, _ := item.OutputCostPerMillionTokens.Float64Value()
		list = append(list, map[string]any{
			"model":                      item.Model,
			"inputTokensSum":             item.InputTokensSum,
			"outputTokensSum":            item.OutputTokensSum,
			"totalTokensSum":             item.TotalTokensSum,
			"inputCostPerMillionTokens":  inc.Float64,
			"outputCostPerMillionTokens": outc.Float64,
			"inputCostCalculated":        (float64(item.InputTokensSum) / 1000000) * inc.Float64,
			"outputCostCalculated":       (float64(item.OutputTokensSum) / 1000000) * outc.Float64,
		})
	}

	request.Encode(w, r, c.logger, http.StatusOK, list)
}
