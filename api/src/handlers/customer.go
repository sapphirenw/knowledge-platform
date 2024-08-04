package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func Customer(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		logger *slog.Logger,
		pool *pgxpool.Pool,
		c *queries.Customer,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		DB(func(w http.ResponseWriter, r *http.Request, logger *slog.Logger, pool *pgxpool.Pool) {
			// parse the customerId
			idStr := chi.URLParam(r, "customerId")
			customerId, err := uuid.Parse(idStr)
			if err != nil {
				logger.Error("Invalid customerId", "customerId", idStr)
				http.Error(w, fmt.Sprintf("Invalid customerId: %s", idStr), http.StatusBadRequest)
				return
			}

			logger.InfoContext(r.Context(), "Fetching the customer...", "customerId", customerId)

			// get the customer
			dmodel := queries.New(pool)
			customer, err := dmodel.GetCustomer(r.Context(), customerId)
			if err != nil {
				// check if no rows
				if err.Error() == "no rows in result set" {
					slogger.ServerError(w, logger, 404, "There was no customers found", err)
					return
				}

				logger.ErrorContext(r.Context(), "Error getting the customer", "error", err)
				http.Error(w, "There was an issue getting the customer", http.StatusInternalServerError)
				return
			}
			logger = logger.With("customerId", customer.ID, "name", customer.Name)

			// pass to the handler function
			handler(w, r, logger, pool, customer)
		}),
	)
}
