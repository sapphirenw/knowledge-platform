package customer

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/datastore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

func getWebsites(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	model := queries.New(pool)
	sites, err := model.GetWebsitesByCustomerWithCount(r.Context(), c.ID)
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

func searchWebsite(
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

	site, err := c.SearchWebsite(r.Context(), &body)
	if err != nil {
		// rollback
		c.logger.Error("failed to parse the website", "error", err)
		if strings.Contains(err.Error(), "REGEX") {
			http.Error(w, "There was an issue with your regex", http.StatusBadRequest)
		} else {
			http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		}
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, site)
}

func insertWebsite(
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

	site, err := c.InsertWebsite(r.Context(), tx, &body)
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

func insertSinglePage(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the body
	body, valid := request.Decode[insertSingleWebsitePageRequest](w, r, c.logger)
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

	if err := c.InsertSinglePage(r.Context(), tx, body.Domain); err != nil {
		// rollback
		tx.Rollback(r.Context())
		slogger.ServerError(w, c.logger, 500, "There was an internal issue", err)
		return
	}

	// send the success
	w.WriteHeader(http.StatusNoContent)
}

func getWebsitePageContent(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
	site *queries.Website,
) {
	logger := c.logger.With("handler", getWebsitePageContent)

	// parse the page id
	pageId, err := utils.GoogleUUIDFromString(chi.URLParam(r, "pageId"))
	if err != nil {
		slogger.ServerError(w, logger, 400, "failed to parse the pageId", err)
		return
	}

	dmodel := queries.New(pool)
	p, err := dmodel.GetWebsitePage(r.Context(), pageId)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get the page", err)
		return
	}

	// create the datastore object
	page, err := datastore.NewWebsitePageFromWebsitePage(r.Context(), logger, p)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to create the internal page datatype", err)
		return
	}

	// parse the options
	getCleaned := r.URL.Query().Get("getCleaned") == "true"
	getChunked := r.URL.Query().Get("getChunked") == "true"

	// get the content
	response := map[string]any{
		"page": p,
	}

	if getCleaned {
		cleaned, err := page.GetCleaned(r.Context())
		if err != nil {
			slogger.ServerError(w, logger, 500, "failed to get the cleaned page data", err)
			return
		}
		response["cleaned"] = cleaned.String()
	}

	if getChunked {
		chunks, err := page.GetChunks(r.Context())
		if err != nil {
			slogger.ServerError(w, logger, 500, "failed to get the chunked page data", err)
			return
		}
		response["chunks"] = chunks
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

func deleteWebsite(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
	site *queries.Website,
) {
	logger := c.logger.With("handler", deleteWebsite)

	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start a transaction", err)
		return
	}
	defer tx.Commit(r.Context())

	dmodel := queries.New(tx)
	if err := dmodel.DeleteWebsite(r.Context(), site.ID); err != nil {
		tx.Rollback(r.Context())
		slogger.ServerError(w, logger, 500, "failed to delete the site", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// func vectorizeWebsite(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	pool *pgxpool.Pool,
// 	c *Customer,
// 	site *queries.Website,
// ) {
// 	// create a transaction
// 	tx, err := pool.Begin(r.Context())
// 	if err != nil {
// 		c.logger.Error("failed to start transaction", "error", err)
// 		http.Error(w, "There was a database issue", http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Commit(r.Context())

// 	if err := c.VectorizeWebsite(r.Context(), tx, site); err != nil {
// 		// rollback
// 		tx.Rollback(r.Context())
// 		c.logger.Error("failed to vecorize website", "error", err)
// 		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// func vectorizeAllWebsites(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	pool *pgxpool.Pool,
// 	c *Customer,
// ) {
// 	// create a transaction
// 	tx, err := pool.Begin(r.Context())
// 	if err != nil {
// 		c.logger.Error("failed to start transaction", "error", err)
// 		http.Error(w, "There was a database issue", http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Commit(r.Context())

// 	if err := c.VectorizeAllWebsites(r.Context(), tx); err != nil {
// 		// rollback
// 		tx.Rollback(r.Context())
// 		c.logger.Error("failed to vecorize websites", "error", err)
// 		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
