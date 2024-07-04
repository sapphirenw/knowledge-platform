package beta

import (
	"net/http"
	"strings"

	"github.com/go-chi/httplog/v2"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

func getCustomer(
	w http.ResponseWriter,
	r *http.Request,
) {
	logger := httplog.LogEntry(r.Context())
	var err error

	// parse the name
	name := r.URL.Query().Get("name")
	if name == "" {
		logger.ErrorContext(r.Context(), "There was no name passed in the url")
		http.Error(w, "'name' is a required url parameter", http.StatusInternalServerError)
		return
	}

	// parse the authToken for this user
	authToken, err := utils.GoogleUUIDFromString(r.URL.Query().Get("authToken"))
	if err != nil {
		slogger.ServerError(w, &logger, 400, "'authToken', is a required url parameter", err)
		return
	}

	// grab a connection from the pool
	pool, err := db.GetPool()
	if err != nil {
		logger.ErrorContext(r.Context(), "Error getting the connection pool", err)
		http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
		return
	}

	// create a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		logger.Error("failed to start transaction", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	dmodel := queries.New(tx)

	// get the api key
	key, err := dmodel.GetBetaApiKey(r.Context(), authToken)
	if err != nil {
		slogger.ServerError(w, &logger, 500, "failed to get the api key", err)
		return
	}

	// ensure that the name matches the key
	if name != key.Name {
		slogger.ServerError(w, &logger, 403, "Not Allowed.", nil)
		return
	}

	// create the customer
	var customer *queries.Customer
	customer, err = dmodel.GetCustomerByName(r.Context(), name)
	if err == nil {
		request.Encode(w, r, &logger, http.StatusOK, customer)
		return
	}

	// check if a new customer needs to be created
	if strings.Contains(err.Error(), "no rows in result set") {
		// create a new customer
		customer, err = dmodel.CreateCustomer(r.Context(), name)
		if err == nil {
			request.Encode(w, r, &logger, http.StatusOK, customer)
			return
		}
	}

	// if here, then there was an issue and changes need to be rolled back
	tx.Rollback(r.Context())
	logger.Error("failed to create or get the test customer", err)
	http.Error(w, "There was an internal issue", http.StatusInternalServerError)
}
