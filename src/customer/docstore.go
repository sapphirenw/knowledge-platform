package customer

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
		if strings.Contains(err.Error(), "violates unique constraint") {
			// this file already exists
			request.Encode(w, r, c.logger, http.StatusConflict, "this file already exists")
			return
		}
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

func createFolder(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the body
	body, valid := request.Decode[*createFolderRequest](w, r, c.logger)
	if !valid {
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

	response, err := c.CreateFolder(r.Context(), tx, body)
	if err != nil {
		c.logger.Error("failed to create the folder", "error", err)
		http.Error(w, "There was an internal server issue", http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

func purgeDatastore(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the body
	body, valid := request.Decode[*purgeDatastoreRequest](w, r, c.logger)
	if !valid {
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

	// send the purge request
	if err = c.PurgeDatastore(r.Context(), tx, body); err != nil {
		tx.Rollback(r.Context())
		c.logger.Error("failed to purge the datastore", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	// let the user know the request was successful
	w.WriteHeader(http.StatusNoContent)
}

func vectorizeDatastore(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// start the transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	if err := c.VectorizeDatastore(r.Context(), tx); err != nil {
		tx.Rollback(r.Context())
		c.logger.Error("failed to vectorize the datastore", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteRemoteDatastore(
	w http.ResponseWriter,
	r *http.Request,
	_ *pgxpool.Pool,
	c *Customer,
) {
	if err := c.DeleteRemoteDatastore(r.Context()); err != nil {
		c.logger.Error("failed to delete the remote datastore", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
