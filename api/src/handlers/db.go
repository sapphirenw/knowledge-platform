package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
)

// Wrapper that passes a database connection to an http handler
func DB(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		logger *slog.Logger,
		pool *pgxpool.Pool,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// create the logger with the request context
			logger := httplog.LogEntry(r.Context())

			// grab a connection from the pool
			pool, err := db.GetPool()
			if err != nil {
				logger.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
				http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
				return
			}

			// pass to the handler function
			handler(w, r, &logger, pool)
		},
	)
}
