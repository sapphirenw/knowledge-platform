package customer

import (
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type vectorizeWebsiteResult struct {
	page    *queries.WebsitePage
	sha256  string
	vectors []*ltypes.EmbeddingsData
}
