package project

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func Handler(mux chi.Router) {
	mux.Get("/", projectHandler(getProject))
	mux.Post("/generateIdeas", projectHandler(generateIdeas))

	mux.Post("/ideas", projectHandler(addIdeas))
	mux.Get("/ideas", projectHandler(getIdeas))
	mux.Get("/ideas/{ideaId}", projectHandler(nil))
}

func projectHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		p *Project,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// create the logger with the request context
			l := httplog.LogEntry(r.Context())

			id := chi.URLParam(r, "projectId")
			projectId, err := uuid.Parse(id)
			if err != nil {
				l.Error("Invalid projectId", "projectId", id)
				http.Error(w, fmt.Sprintf("Invalid projectId: %s", id), http.StatusBadRequest)
				return
			}

			// get a connection pool
			pool, err := db.GetPool()
			if err != nil {
				l.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
				http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
				return
			}

			// get the folder from the db
			project, err := GetProject(r.Context(), &l, pool, projectId)
			if err != nil {
				defer pool.Close() // ensure the pool gets released
				// check if no rows
				if strings.Contains(err.Error(), "no rows in result set") {
					http.NotFound(w, r)
					return
				}

				l.Error("Error getting the project", "error", err)
				http.Error(w, fmt.Sprintf("There was no project found with projectId: %s", id), http.StatusNotFound)
				return
			}

			// pass to the handler
			handler(w, r, pool, project)
		},
	)
}

func getProject(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	p *Project,
) {
	request.Encode(w, r, p.logger, http.StatusOK, p)
}

func generateIdeas(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	p *Project,
) {
	// parse the request
	body, valid := request.Decode[generateIdeasRequest](w, r, p.logger)
	if !valid {
		return
	}

	// create a tx
	tx, err := pool.Begin(r.Context())
	if err != nil {
		p.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// create the ideas
	response, err := p.GenerateIdeas(r.Context(), tx, &body)
	if err != nil {
		p.logger.Error("failed to generate ideas", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// return to the user
	request.Encode(w, r, p.logger, http.StatusOK, response)
}

func addIdeas(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	p *Project,
) {
	// parse the request
	body, valid := request.Decode[addIdeasRequest](w, r, p.logger)
	if !valid {
		return
	}

	// create a tx
	tx, err := pool.Begin(r.Context())
	if err != nil {
		p.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// create the ideas
	response, err := p.AddIdeas(r.Context(), tx, &body)
	if err != nil {
		p.logger.Error("failed to add ideas", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// return to the user
	request.Encode(w, r, p.logger, http.StatusOK, response)
}

func getIdeas(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	p *Project,
) {
	// create a tx
	tx, err := pool.Begin(r.Context())
	if err != nil {
		p.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// create the ideas
	response, err := p.GetIdeas(r.Context(), tx)
	if err != nil {
		p.logger.Error("failed to get ideas", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// return to the user
	request.Encode(w, r, p.logger, http.StatusOK, response)
}
