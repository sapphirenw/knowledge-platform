package beta

import (
	"net/http"
	"strings"

	"github.com/go-chi/httplog/v2"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
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
		http.Error(w, "`name` is a required url parameter", http.StatusInternalServerError)
		return
	}

	// grab a connection from the pool
	pool, err := db.GetPool()
	if err != nil {
		logger.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
		http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
		return
	}

	// create a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// create the customer
	var customer *queries.Customer
	model := queries.New(tx)
	customer, err = model.GetCustomerByName(r.Context(), name)
	if err == nil {
		request.Encode(w, r, &logger, http.StatusOK, customer)
		return
	}

	// check if a new customer needs to be created
	if strings.Contains(err.Error(), "no rows in result set") {
		// create a new customer
		customer, err = model.CreateCustomer(r.Context(), name)
		if err == nil {
			request.Encode(w, r, &logger, http.StatusOK, customer)
			return
		}
	}

	// if here, then there was an issue and changes need to be rolled back
	tx.Rollback(r.Context())
	logger.Error("failed to create or get the test customer", "error", err)
	http.Error(w, "There was an internal issue", http.StatusInternalServerError)
}
