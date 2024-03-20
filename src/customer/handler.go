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
	mux.Get("/", customerHandler(getCustomer))

	// documents
	mux.Post("/generatePresignedUrl", customerHandler(generatePresignedUrl))
	mux.Get("/documents/{documentId}", customerHandler(notifyOfSuccessfulUpload))
	mux.Put("/documents/{documentId}/validate", customerHandler(notifyOfSuccessfulUpload))

	// folders
	mux.Get("/root", customerHandler(listCustomerFolder))
	mux.Get("/folders/{folderId}", customerHandler(listCustomerFolder))

	// websites
	mux.Get("/websites", customerHandler(getWebsites))
	mux.Post("/websites", customerHandler(handleWesbite))
	mux.Get("/websites/{websiteId}", customerHandler(getWebsite))
	mux.Get("/websites/{websiteId}/pages", customerHandler(getWebsitePages))
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

func parseDocumentFromRequest(
	r *http.Request,
	db queries.DBTX,
) (*queries.Document, error) {
	id := chi.URLParam(r, "documentId")
	documentId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid parameter: %v", id)
	}

	// get the document from the db
	model := queries.New(db)
	doc, err := model.GetDocument(r.Context(), documentId)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the document: %v", err)
	}

	return doc, nil
}

func parseFolderFromRequest(
	r *http.Request,
	db queries.DBTX,
) (*queries.Folder, error) {
	id := chi.URLParam(r, "folderId")
	folderId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid parameter: %v", id)
	}

	// get the document from the db
	model := queries.New(db)
	folder, err := model.GetFolder(r.Context(), folderId)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the folder: %v", err)
	}

	return folder, nil
}

func parseSiteFromRequest(
	r *http.Request,
	db queries.DBTX,
) (*queries.Website, error) {
	id := chi.URLParam(r, "websiteId")
	websiteId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid parameter: %v", id)
	}

	// get the document from the db
	model := queries.New(db)
	website, err := model.GetWebsite(r.Context(), websiteId)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the folder: %v", err)
	}

	return website, nil
}

func getCustomer(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// return the customer
	request.Encode(w, r, c.logger, http.StatusOK, c.Customer)
}
