package llm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
)

func Handler(mux chi.Router) {
	mux.Get("/", nil)

	mux.Route("/{conversationId}", func(r chi.Router) {
		r.Get("/", conversationHandler(getConversation))
	})
}

func conversationHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		conv *Conversation,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// create the logger with the request context
			l := httplog.LogEntry(r.Context())

			id := chi.URLParam(r, "conversationId")
			conversationId, err := uuid.Parse(id)
			if err != nil {
				l.Error("Invalid conversationId", "conversationId", id)
				http.Error(w, fmt.Sprintf("Invalid conversationId: %s", id), http.StatusBadRequest)
				return
			}

			// get a connection pool
			pool, err := db.GetPool()
			if err != nil {
				l.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
				http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
				return
			}

			// get the folder from the db
			conv, err := GetConversation(r.Context(), &l, pool, conversationId)
			if err != nil {
				defer pool.Close() // ensure the pool gets released
				// check if no rows
				if strings.Contains(err.Error(), "no rows in result set") {
					l.Error("Not found", "error", err)
					http.NotFound(w, r)
					return
				}

				l.Error("Error getting the conversation", "error", err)
				http.Error(w, fmt.Sprintf("There was no project found with conversationId: %s", id), http.StatusNotFound)
				return
			}

			// pass to the handler
			handler(w, r, pool, conv)
		},
	)
}

type GetConverstaionResponse struct {
	ConversationId uuid.UUID        `json:"conversationId"`
	Title          string           `json:"title"`
	Messages       []*gollm.Message `json:"messages"`
}

func getConversation(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Conversation,
) {
	// return relevant conversation information
	request.Encode(w, r, c.logger, http.StatusOK, &GetConverstaionResponse{
		ConversationId: c.ID,
		Title:          c.Title,
		Messages:       c.GetMessages(),
	})
}

// func getConversations(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {
// 	logger := httplog.LogEntry(r.Context())

// 	// get the database
// 	pool, err := db.GetPool()
// 	if err != nil {
// 		logger.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
// 		http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
// 		return
// 	}
// 	dmodel := queries.New(pool)

// 	// fetch the conversations
// 	response, err := dmodel.GetConversationsWithCount(ctx, )
// }
