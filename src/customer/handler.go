package customer

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func Handler(mux chi.Router) {
	mux.Get("/", customerHandler(GetCustomer))
	mux.Post("/generatePresignedUrl", customerHandler(GeneratePresignedUrl))

	// documents
	mux.Get("/documents/{documentId}", customerHandler(NotifyOfSuccessfulUpload))
	mux.Put("/documents/{documentId}/validate", customerHandler(NotifyOfSuccessfulUpload))

	// folders
	mux.Get("/root", customerHandler(ListCustomerFolder))
	mux.Get("/folders/{folderId}", customerHandler(ListCustomerFolder))
}

// Custom handler that parses the customerId from the request, fetches the customer from the database
// and passes a valid database connection pool writer to the handler
func customerHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		c *Customer,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// create the logger with the request context
			l := httplog.LogEntry(r.Context())

			// parse the customerId
			idStr := chi.URLParam(r, "customerId")
			customerId, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				l.Error("Invalid customerId", "customerId", idStr)
				http.Error(w, fmt.Sprintf("Invalid customerId: %s", idStr), http.StatusBadRequest)
				return
			}

			// grab a connection from the pool
			pool, err := db.GetPool()
			if err != nil {
				l.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
				http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
				return
			}

			l.InfoContext(r.Context(), "Fetching the customer...", "customerId", customerId)

			// get the customer
			customer, err := NewCustomer(r.Context(), &l, customerId, pool)
			if err != nil {
				pool.Close() // ensure the pool gets released
				// check if no rows
				if err.Error() == "no rows in result set" {
					http.NotFound(w, r)
					return
				}

				l.ErrorContext(r.Context(), "Error getting the customer", "error", err)
				http.Error(w, "There was an issue getting the customer", http.StatusInternalServerError)
				return
			}

			// pass to the handler function
			handler(w, r, pool, customer)
		},
	)
}

func GetCustomer(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// return the customer
	request.Encode(w, r, c.logger, http.StatusOK, c.Customer)
}

func ListCustomerFolder(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	var err error
	var folder *queries.Folder

	// if the param was passed then get the folder
	if chi.URLParam(r, "folderId") != "" {
		// parse the folder from the query args
		folder, err = parseFolderFromRequest(r, pool)
		if err != nil {
			c.logger.ErrorContext(r.Context(), "Error getting folder", "error", err)
			http.Error(w, "There was an issue getting the folder", http.StatusInternalServerError)
			return
		}
	}

	// fetch the data inside the customer's folder
	response, err := c.ListFolderContents(r.Context(), pool, folder)
	if err != nil {
		c.logger.ErrorContext(r.Context(), "Error listing folder contents", "error", err)
		http.Error(w, "There was an internal server error", http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

func GeneratePresignedUrl(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the body
	body, valid := request.Decode[generatePresignedUrlRequest](w, r, c.logger)
	if !valid {
		return
	}

	// use the customer to generate the presigned url
	response, err := c.GeneratePresignedUrl(r.Context(), pool, &body)
	if err != nil {
		c.logger.ErrorContext(r.Context(), "generating the presigned url", "error", err)
		http.Error(w, "There was an issue generating the presigned url", http.StatusInternalServerError)
		return
	}

	// return the response to the user
	request.Encode(w, r, c.logger, http.StatusOK, response)
}

func NotifyOfSuccessfulUpload(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the documentId
	idStr := chi.URLParam(r, "documentId")
	documentId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Error("Invalid documentId", "documentId", idStr)
		http.Error(w, fmt.Sprintf("Invalid documentId: %s", idStr), http.StatusBadRequest)
		return
	}

	// start the transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// send the validation request against the customer
	if err = c.NotifyOfSuccessfulUpload(r.Context(), tx, documentId); err != nil {
		tx.Rollback(r.Context())
		c.logger.Error("failed to validate the document record", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}

	// let the user know the request was successful
	w.WriteHeader(http.StatusNoContent)
}

func GetDocument(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	doc, err := parseDocumentFromRequest(r, pool)
	if err != nil {
		c.logger.Error("failed to get the document", "error", err)
		http.Error(w, "There was an internal server issue", http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, doc)
}
