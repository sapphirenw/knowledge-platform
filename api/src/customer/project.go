package customer

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/project"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func createProject(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	body, valid := request.Decode[createProjectRequest](w, r, c.logger)
	if !valid {
		return
	}

	// create a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// create the project
	p, err := project.CreateProject(r.Context(), tx, c.logger, c.Customer, body.Title, body.Topic, nil)
	if err != nil {
		tx.Rollback(r.Context())
		c.logger.Error("failed to create the project", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// send to the user
	request.Encode(w, r, c.logger, http.StatusOK, p)
}
