package customer

import (
	"github.com/sapphirenw/ai-content-creation-api/src/embeddings"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type vectorizeWebsiteResult struct {
	page    *queries.WebsitePage
	sha256  string
	vectors []*embeddings.EmbeddingsData
}
