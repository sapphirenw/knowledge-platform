package middleware

import (
	"net/http"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

// Beta middleware to handle blanket auth to Beta testers. This should NOT
// be used long term. In the future, if the 'x-api-key' header is included,
// the request should be rejected.
func BetaAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get a connection from the pool
		pool, err := db.GetPool()
		if err != nil {
			slogger.ServerError(w, r, nil, 500, "failed to connect to the database")
			return
		}

		// get a specific connection to use separately
		conn, err := pool.Acquire(r.Context())
		if err != nil {
			slogger.ServerError(w, r, nil, 500, "failed to grab a connection", "error", err)
			return
		}
		defer conn.Release()

		// parse as a uuid
		apiKey, err := utils.GoogleUUIDFromString(r.Header.Get("x-api-key"))
		if err != nil {
			slogger.ServerError(w, r, nil, 403, "Not Allowed.", "error", err)
			return
		}

		// get the api key
		dmodel := queries.New(conn)
		key, err := dmodel.GetBetaApiKey(r.Context(), apiKey)
		if err != nil {
			slogger.ServerError(w, r, nil, 403, "Not Allowed.", "error", err)
			return
		}

		// essure key valid
		if key.Expired {
			slogger.ServerError(w, r, nil, 403, "Not Allowed.", "message", "api key is expired")
			return
		}

		// pass to next handler
		next.ServeHTTP(w, r)
	})
}
