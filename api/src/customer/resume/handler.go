package resume

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/handlers"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func Handler(r chi.Router) {
	// r.Post("/", handlers.Customer(createResumeHandler))
	// r.Get("/", handlers.Customer(getAllResumesHandler))

	r.Route("/{resumeId}", func(r chi.Router) {
		r.Get("/", resumeHandler(getCustomerResume))
		r.Get("/about", resumeHandler(getResumeAboutHandler))
		r.Post("/about", resumeHandler(setResumeAboutHandler))
		r.Post("/setTitle", resumeHandler(setResumeTitleHandler))

		r.Get("/checklist", resumeHandler(getResumeChecklistHandler))

		// docs
		r.Get("/resume", resumeHandler(getResumeHandler))
		r.Post("/documents", resumeHandler(attachDocumentsHandler))
		r.Get("/documents", resumeHandler(getDocumentsHandler))

		// work experience
		r.Route("/workExperiences/{experienceId}", func(r chi.Router) {

		})

		// applications
		r.Get("/applications", resumeHandler(getResumeApplicationsHandler))
		r.Post("/applications", resumeHandler(createResumeApplicationHandler))
		r.Route("/applications/{applicationId}", func(r chi.Router) {})
	})
}

func resumeHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		logger *slog.Logger,
		pool *pgxpool.Pool,
		client *Client,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		handlers.Customer(func(w http.ResponseWriter, r *http.Request, logger *slog.Logger, pool *pgxpool.Pool, c *queries.Customer) {
			tx, err := pool.Begin(r.Context())
			if err != nil {
				slogger.ServerError(w, logger, 500, "failed to start the transaction", err)
			}
			defer tx.Commit(r.Context())

			client, err := NewClient(r.Context(), logger, tx, c, chi.URLParam(r, "resumeId"))
			if err != nil {
				slogger.ServerError(w, logger, 400, "failed to get the resume", err)
				return
			}

			handler(w, r, logger.With("resumeId", client.Resume.ID.String()), pool, client)
		}),
	)
}

func getResumeAboutHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	about, err := client.GetAbout(r.Context(), logger, pool)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get the resume about", err)
		return
	}

	request.Encode(w, r, logger, 200, about)
}

// func createResumeHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	logger *slog.Logger,
// 	pool *pgxpool.Pool,
// 	c *queries.Customer,
// ) {
// 	// parse the body
// 	body, valid := request.Decode[setResumeTitleRequest](w, r, logger)
// 	if !valid {
// 		return
// 	}

// 	tx, err := pool.Begin(r.Context())
// 	if err != nil {
// 		slogger.ServerError(w, logger, 500, "failed to start the transaction", err)
// 	}
// 	defer tx.Commit(r.Context())

// 	dmodel := queries.New(tx)
// 	resume, err := dmodel.CreateResume(r.Context(), &queries.CreateResumeParams{
// 		CustomerID: c.ID,
// 		Title:      body.Title,
// 	})
// 	if err != nil {
// 		slogger.ServerError(w, logger, 500, "failed to create the resume", err)
// 		tx.Rollback(r.Context())
// 		return
// 	}

// 	request.Encode(w, r, logger, 200, resume)
// }

func setResumeTitleHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	// parse the body
	body, valid := request.Decode[setResumeTitleRequest](w, r, logger)
	if !valid {
		return
	}

	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start the transaction", err)
		return
	}
	defer tx.Commit(r.Context())

	if err := client.SetTitle(r.Context(), logger, pool, body.Title); err != nil {
		slogger.ServerError(w, logger, 500, "failed to set the title", err)
		tx.Rollback(r.Context())
		return
	}

	w.WriteHeader(200)
}

func setResumeAboutHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	// parse the body
	var body queries.CreateResumeAboutParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slogger.ServerError(w, logger, 500, "failed to read the body", err)
		return
	}

	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start the transaction", err)
		return
	}
	defer tx.Commit(r.Context())

	if _, err := client.SetAbout(r.Context(), logger, tx, &body); err != nil {
		slogger.ServerError(w, logger, 500, "failed to set the resume about", err)
		tx.Rollback(r.Context())
		return
	}

	w.WriteHeader(200)
}

// func getAllResumesHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	logger *slog.Logger,
// 	pool *pgxpool.Pool,
// 	c *queries.Customer,
// ) {
// 	dmodel := queries.New(pool)
// 	response, err := dmodel.GetResumesCustomer(r.Context(), c.ID)
// 	if err != nil {
// 		slogger.ServerError(w, logger, 500, "failed to get the resumes", err)
// 		return
// 	}

// 	request.Encode(w, r, logger, 200, response)
// }

func getCustomerResume(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	request.Encode(w, r, logger, 200, client.Resume)
}

func getResumeApplicationsHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	apps, err := client.GetApplications(r.Context(), logger, pool)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get applications", err)
		return
	}

	request.Encode(w, r, logger, 200, apps)
}

func createResumeApplicationHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	// parse the body
	var body createResumeApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slogger.ServerError(w, logger, 500, "failed to read the body", err)
		return
	}

	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start the transaction", err)
		return
	}
	defer tx.Commit(r.Context())

	apps, err := client.CreateApplication(r.Context(), logger, tx, &body)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to create the application", err)
		return
	}

	request.Encode(w, r, logger, 200, apps)
}

func getResumeChecklistHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	response, err := client.GetChecklist(r.Context(), logger, pool)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get the checklist", err)
	}

	request.Encode(w, r, logger, 200, response)
}

func attachDocumentsHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	// parse the body
	var body attachDocumentsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slogger.ServerError(w, logger, 500, "failed to read the body", err)
		return
	}

	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start the transaction", err)
		return
	}
	defer tx.Commit(r.Context())

	for _, id := range body.DocumentIDs {
		if err := client.AttachDocument(r.Context(), logger, tx, id, false); err != nil {
			tx.Rollback(r.Context())
			slogger.ServerError(w, logger, 500, "failed to attach the document", err, "docId", id)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getDocumentsHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	docs, err := client.GetDocuments(r.Context(), logger, pool)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get the documents", err)
		return
	}

	request.Encode(w, r, logger, 200, docs)
}

func getResumeHandler(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	client *Client,
) {
	doc, err := client.GetResume(r.Context(), logger, pool)
	if err != nil {
		if errors.Is(errors.Unwrap(err), pgx.ErrNoRows) {
			w.WriteHeader(404)
			return
		}
		slogger.ServerError(w, logger, 500, "failed to get the resume", err)
		return
	}

	request.Encode(w, r, logger, 200, doc)
}
