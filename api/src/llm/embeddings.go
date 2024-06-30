package llm

import (
	"log/slog"

	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func GetEmbeddings(logger *slog.Logger, c *queries.Customer) gollm.Embeddings {
	emb := gollm.NewOpenAIEmbeddings(c.ID.String(), logger, nil)
	return emb
}
