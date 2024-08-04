package customer

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/project"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/resume"
	"github.com/sapphirenw/ai-content-creation-api/src/datastore"
	"github.com/sapphirenw/ai-content-creation-api/src/handlers"
	"github.com/sapphirenw/ai-content-creation-api/src/middleware"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func Handler(mux chi.Router) {
	mux.Use(middleware.BetaAuthToken)
	mux.Get("/", customerHandler(getCustomer))

	// language model configurations
	mux.Put("/updateLLMConfigurations", customerHandler(updateCustomerLLMConfigurations))
	mux.Route("/llms", func(r chi.Router) {
		r.Get("/", customerHandler(getAvailableLLMs))
		r.Post("/", customerHandler(createModel))
		r.Route("/{llmId}", func(r chi.Router) {
			r.Put("/", customerHandler(updateModel))
		})
	})

	// datastore

	mux.Route("/datastore", func(r chi.Router) {
		r.Delete("/", customerHandler(deleteRemoteDatastore))
		r.Post("/purge", customerHandler(purgeDatastore))
	})

	// documents
	mux.Post("/generatePresignedUrl", customerHandler(generatePresignedUrl))
	mux.Route("/documents/{documentId}", func(r chi.Router) {
		r.Get("/", documentHandler(getDocument))
		r.Put("/validate", documentHandler(notifyOfSuccessfulUpload))
		r.Get("/raw", documentHandler(getDocumentRaw))
		r.Get("/cleaned", documentHandler(getDocumentCleaned))
		r.Get("/chunked", documentHandler(getDocumentChunked))
	})

	// folders
	mux.Route("/folders", func(r chi.Router) {
		r.Get("/", customerHandler(listCustomerFolder))
		r.Post("/", customerHandler(createFolder))
		r.Get("/{folderId}", customerHandler(listCustomerFolder))
	})

	// websites
	mux.Post("/insertSingleWebsitePage", customerHandler(insertSinglePage))
	mux.Route("/websites", func(r chi.Router) {
		r.Get("/", customerHandler(getWebsites))
		r.Put("/", customerHandler(searchWebsite))
		r.Post("/", customerHandler(insertWebsite))
		r.Route("/{websiteId}", func(r chi.Router) {
			r.Get("/", websiteHandler(getWebsite))
			r.Delete("/", websiteHandler(deleteWebsite))
			r.Get("/pages", websiteHandler(getWebsitePages))
			r.Get("/pages/{pageId}/content", websiteHandler(getWebsitePageContent))
		})
	})

	// vectorstore
	mux.Route("/vectorstore", func(r chi.Router) {
		r.Put("/query", customerHandler(queryVectorStore))
		r.Put("/queryDocs", customerHandler(queryVectorStoreDocuments))
		r.Put("/queryWebsitePages", customerHandler(queryVectorStoreWebsitePages))
		r.Put("/queryRaw", customerHandler(queryVectorStoreRaw))
		r.Get("/vectorize", customerHandler(getAllVectorizeRequests))
		r.Post("/vectorize", customerHandler(createVectorizeRequest))
		r.Get("/vectorize/{id}", customerHandler(getVectorizeRequest))
	})

	// usage
	mux.Route("/usage", func(r chi.Router) {
		r.Get("/", customerHandler(getUsage))
	})
	mux.Get("/usageGrouped", customerHandler(getUsageGrouped))

	// project
	mux.Route("/projects", project.Handler)
	// mux.Route("/projects", func(r chi.Router) {
	// 	r.Post("/", customerHandler(createProject))
	// })

	// conversations
	mux.Route("/conversations", conversation.Handler)

	// rag
	mux.Post("/rag", customerHandler(handleRAG))
	mux.Get("/rag2", customerHandler(handleRag2))

	// resume
	mux.Route("/resumes", resume.Handler)

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
		handlers.Customer(func(w http.ResponseWriter, r *http.Request, logger *slog.Logger, pool *pgxpool.Pool, c *queries.Customer) {
			handler(w, r, pool, &Customer{Customer: c, logger: logger})
		}),
	)
}

func documentHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		c *Customer,
		doc *datastore.Document,
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
			doc, err := datastore.GetDocument(r.Context(), c.logger, pool, documentId)
			if err != nil {
				c.logger.Error("Error parsing as a docstore doc", "error", err)
				http.Error(w, fmt.Sprintf("There was an internal issue: %w", err), http.StatusInternalServerError)
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
