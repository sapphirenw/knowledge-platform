package beta

import (
	"net/http"
	"os"

	"github.com/go-chi/httplog/v2"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func createBetaApiKey(
	w http.ResponseWriter,
	r *http.Request,
) {
	logger := httplog.LogEntry(r.Context())

	// read the master auth token
	serverAuthToken, exists := os.LookupEnv("API_MASTER_AUTH_TOKEN")
	if !exists {
		slogger.ServerError(w, &logger, 500, "unknown failure", nil)
		return
	}

	authToken := r.Header.Get("x-master-auth-token")
	if authToken == "" {
		slogger.ServerError(w, &logger, 403, "Not Allowed.", nil, "authStatus", "EMPTY")
		return
	}

	if authToken != serverAuthToken {
		slogger.ServerError(w, &logger, 403, "Not Allowed.", nil, "authStatus", "NOT_EQUAL")
		return
	}

	// check for the parameter
	name := r.URL.Query().Get("name")
	if name == "" {
		slogger.ServerError(w, &logger, 400, "'name' is a required url parameter", nil)
		return
	}
	isAdmin := r.URL.Query().Get("isAdmin") == "true"

	// get a connection
	pool, err := db.GetPool()
	if err != nil {
		slogger.ServerError(w, &logger, 500, "failed to connect to the database", nil)
		return
	}

	// send the create request
	dmodel := queries.New(pool)
	key, err := dmodel.CreateBetaApiKey(r.Context(), &queries.CreateBetaApiKeyParams{
		Name:    name,
		IsAdmin: isAdmin,
	})
	if err != nil {
		slogger.ServerError(w, &logger, 500, "failed to create the api key", err)
		return
	}

	request.Encode(w, r, &logger, http.StatusOK, key)
}
