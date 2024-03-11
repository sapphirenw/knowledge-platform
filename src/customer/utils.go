package customer

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

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
