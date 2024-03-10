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

func CustomerHandler(mux chi.Router) {
	mux.Get("/", handle(GetCustomer))
	mux.Post("/generatePresignedUrl", handle(GeneratePresignedUrl))

	// documents
	mux.Put("/documents/{documentId}/validate", handle(NotifyOfSuccessfulUpload))

	// folders
	mux.Get("/folders/{folderId}", handle(ListCustomerFolder))
}

// Custom handler that parses the customerId from the request, fetches the customer from the database
// and passes a valid database connection pool writer to the handler
func handle(handler func(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, c *Customer)) http.HandlerFunc {
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
	if err := request.Encode(w, r, http.StatusOK, c.Customer); err != nil {
		c.logger.ErrorContext(r.Context(), "Error encoding data", "error", err)
		http.Error(w, "There was an internal server error", http.StatusInternalServerError)
		return
	}
}

func ListCustomerFolder(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the folder title from the query args
	fidStr := chi.URLParam(r, "folderId")
	folderId, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		c.logger.Error("Invalid folderId", "folderId", fidStr)
		http.Error(w, fmt.Sprintf("Invalid folderId: %s", fidStr), http.StatusBadRequest)
		return
	}

	// fetch the folder using the foldername
	model := queries.New(pool)
	folder, err := model.GetFolder(r.Context(), folderId)
	if err != nil {
		c.logger.ErrorContext(r.Context(), "Error getting folder", "error", err)
		http.Error(w, "There was an issue getting the folder", http.StatusInternalServerError)
		return
	}

	// fetch the data inside the customer's folder
	response, err := c.ListFolderContents(r.Context(), pool, folder)
	if err != nil {
		c.logger.ErrorContext(r.Context(), "Error listing folder contents", "error", err)
		http.Error(w, "There was an internal server error", http.StatusInternalServerError)
		return
	}

	if err := request.Encode(w, r, http.StatusOK, response); err != nil {
		c.logger.ErrorContext(r.Context(), "Error encoding data", "error", err)
		http.Error(w, "There was an internal server error", http.StatusInternalServerError)
		return
	}
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
	if err := request.Encode(w, r, http.StatusOK, response); err != nil {
		c.logger.ErrorContext(r.Context(), "Error encoding data", "error", err)
		http.Error(w, "There was an internal server error", http.StatusInternalServerError)
		return
	}
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
		c.logger.Error("failed to validate the document record", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		tx.Rollback(r.Context())
		return
	}

	// let the user know the request was successful
	w.WriteHeader(http.StatusNoContent)
}
