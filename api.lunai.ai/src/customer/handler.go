package customer

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func Handler(mux chi.Router) {
	mux.Get("/", customerHandler(getCustomer))

	// datastore
	mux.Delete("/datastore", customerHandler(deleteRemoteDatastore))
	mux.Post("/datastore/purge", customerHandler(purgeDatastore))

	// documents
	mux.Put("/vectorizeDocuments", customerHandler(vectorizeDatastore))
	mux.Post("/generatePresignedUrl", customerHandler(generatePresignedUrl))
	mux.Get("/documents/{documentId}", documentHandler(getDocument))
	mux.Put("/documents/{documentId}/validate", documentHandler(notifyOfSuccessfulUpload))
	// mux.Put("/documents/{documentId}/vectorize", customerHandler(notifyOfSuccessfulUpload))

	// folders
	mux.Get("/root", customerHandler(listCustomerFolder))
	mux.Get("/folders/{folderId}", customerHandler(listCustomerFolder))
	mux.Post("/folders", customerHandler(createFolder))

	// websites
	mux.Get("/websites", customerHandler(getWebsites))
	mux.Post("/websites", customerHandler(handleWesbite))
	mux.Put("/vectorizeWebsites", customerHandler(vectorizeAllWebsites))
	mux.Get("/websites/{websiteId}", websiteHandler(getWebsite))
	mux.Get("/websites/{websiteId}/pages", websiteHandler(getWebsitePages))
	mux.Put("/websites/{websiteId}/vectorize", websiteHandler(vectorizeWebsite))

	// vectorstore
	mux.Put("/vectorstore/query", customerHandler(queryVectorStore))

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
			customerId, err := uuid.Parse(idStr)
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

func documentHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		c *Customer,
		doc *docstore.Document,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		customerHandler(func(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, c *Customer) {
			// scan the docId into a uuid
			docId := chi.URLParam(r, "documentId")
			documentId, err := uuid.Parse(docId)
			if err != nil {
				c.logger.Error("Invalid documentId", "documentId", docId)
				http.Error(w, fmt.Sprintf("Invalid documentId: %s", docId), http.StatusBadRequest)
				return
			}

			// parse as a docstore doc
			doc, err := docstore.NewDocument(r.Context(), pool, documentId)
			if err != nil {
				c.logger.Error("Error parsing as a docstore doc", "error", err)
				http.Error(w, fmt.Sprintf("There was an internal issue: %s", err), http.StatusInternalServerError)
				return
			}

			// pass to the handler
			handler(w, r, pool, c, doc)
		}),
	)
}

func websiteHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		c *Customer,
		site *queries.Website,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		customerHandler(func(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, c *Customer) {
			id := chi.URLParam(r, "websiteId")
			siteId, err := uuid.Parse(id)
			if err != nil {
				c.logger.Error("Invalid folderId", "siteId", id)
				http.Error(w, fmt.Sprintf("Invalid siteId: %s", id), http.StatusBadRequest)
				return
			}

			// get the folder from the db
			model := queries.New(pool)
			site, err := model.GetWebsite(r.Context(), siteId)
			if err != nil {
				c.logger.Error("Error getting the website", "error", err)
				http.Error(w, fmt.Sprintf("There was no site found with websiteId: %s", id), http.StatusNotFound)
				return
			}

			// pass to the handler
			handler(w, r, pool, c, site)
		}),
	)
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
