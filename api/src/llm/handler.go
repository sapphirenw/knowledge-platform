package llm

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func Handler(mux chi.Router) {
	mux.Get("/publicModels", getPublicModels)
}

func getPublicModels(w http.ResponseWriter, r *http.Request) {
	l := httplog.LogEntry(r.Context())
	logger := l.With("handler", "getPublicModels")

	pool, err := db.GetPool()
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get the database", err)
		return
	}

	// get the llms
	dmodel := queries.New(pool)
	models, err := dmodel.GetPublicLLMs(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to get the llms", err)
		return
	}

	request.Encode(w, r, logger, 200, models)
}
