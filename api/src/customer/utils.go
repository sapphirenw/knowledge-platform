package customer

import (
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/webparse"
)

type vectorizeWebsiteResult struct {
	page    *queries.WebsitePage
	headers *webparse.ScrapeHeader
	sha256  string
	vectors []*ltypes.EmbeddingsData
}
