package customer

import (
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func getWebsites(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	model := queries.New(pool)
	sites, err := model.GetWebsitesByCustomer(r.Context(), c.ID)
	if err != nil {
		c.logger.Error("failed to get the websites", "error", err)
		http.Error(w, "There was an internal server issue", http.StatusInternalServerError)
		return
	}
	request.Encode(w, r, c.logger, http.StatusOK, sites)
}

func getWebsite(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
	site *queries.Website,
) {
	request.Encode(w, r, c.logger, http.StatusOK, site)
}

func getWebsitePages(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
	site *queries.Website,
) {
	// get the pages
	model := queries.New(pool)
	pages, err := model.GetWebsitePagesBySite(r.Context(), site.ID)
	if err != nil {
		c.logger.Error("failed to get the website pages", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, pages)
}

func handleWesbite(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the body
	body, valid := request.Decode[handleWebsiteRequest](w, r, c.logger)
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

	site, err := c.HandleWebsite(r.Context(), tx, &body)
	if err != nil {
		// rollback
		tx.Rollback(r.Context())
		c.logger.Error("failed to insert the website", "error", err)
		if strings.Contains(err.Error(), "REGEX") {
			http.Error(w, "There was an issue with your regex", http.StatusBadRequest)
		} else {
			http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		}
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, site)
}

func vectorizeWebsite(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
	site *queries.Website,
) {
	// create a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	if err := c.VectorizeWebsite(r.Context(), tx, site); err != nil {
		// rollback
		tx.Rollback(r.Context())
		c.logger.Error("failed to vecorize website", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func vectorizeAllWebsites(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// create a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	if err := c.VectorizeAllWebsites(r.Context(), tx); err != nil {
		// rollback
		tx.Rollback(r.Context())
		c.logger.Error("failed to vecorize websites", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
