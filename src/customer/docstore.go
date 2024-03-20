package customer

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func listCustomerFolder(
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

func generatePresignedUrl(
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

func notifyOfSuccessfulUpload(
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

func getDocument(
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
