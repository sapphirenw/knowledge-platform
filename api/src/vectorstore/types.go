package vectorstore

import (
	"github.com/sapphirenw/ai-content-creation-api/src/docstore"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type DocumentResponse struct {
	*docstore.Document

	Content string `json:"content"`
}

type WebsitePageResponse struct {
	*queries.WebsitePage

	Content string `json:"content"`
}
