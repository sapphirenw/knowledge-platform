package conversation

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
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func Handler(mux chi.Router) {
	mux.Get("/", rootHandler(getConversations))

	mux.Route("/{conversationId}", func(r chi.Router) {
		r.Get("/", conversationHandler(getConversation))
	})
}

func rootHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		customer *queries.Customer,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// create the logger with the request context
			l := httplog.LogEntry(r.Context())

			id := chi.URLParam(r, "customerId")
			customerId, err := uuid.Parse(id)
			if err != nil {
				l.Error("Invalid customerId", "customerId", id)
				http.Error(w, fmt.Sprintf("Invalid customerId: %s", id), http.StatusBadRequest)
				return
			}

			// get a connection pool
			pool, err := db.GetPool()
			if err != nil {
				l.ErrorContext(r.Context(), "Error getting the connection pool", "error", err)
				http.Error(w, "There was an issue connecting to the database", http.StatusInternalServerError)
				return
			}

			// get the customer from the db
			fmt.Println(customerId)
			dmodel := queries.New(pool)
			customer, err := dmodel.GetCustomer(r.Context(), customerId)
			if err != nil {
				// check if no rows
				if strings.Contains(err.Error(), "no rows in result set") {
					slogger.ServerError(w, &l, 404, "failed to get the customer", err)
					return
				}

				l.Error("Error getting the customer", "error", err)
				http.Error(w, fmt.Sprintf("There was not a customer found with customerId: %s", id), http.StatusNotFound)
				return
			}

			// pass to the handler
			handler(w, r, pool, customer)
		},
	)
}

func conversationHandler(
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		pool *pgxpool.Pool,
		customer *queries.Customer,
		conv *Conversation,
	),
) http.HandlerFunc {
	return http.HandlerFunc(
		rootHandler(func(
			w http.ResponseWriter,
			r *http.Request,
			pool *pgxpool.Pool,
			customer *queries.Customer,
		) {
			logger := httplog.LogEntry(r.Context())

			// scan the docId into a uuid
			id := chi.URLParam(r, "conversationId")
			convId, err := uuid.Parse(id)
			if err != nil {
				logger.Error("Invalid conversationId", "conversationId", convId)
				http.Error(w, fmt.Sprintf("Invalid conversationId: %s", convId), http.StatusBadRequest)
				return
			}

			// parse as a docstore doc
			conv, err := GetConversation(r.Context(), &logger, pool, convId)
			if err != nil {
				logger.Error("Error parsing as a docstore doc", "error", err)
				http.Error(w, fmt.Sprintf("There was an internal issue: %s", err), http.StatusInternalServerError)
				return
			}

			// pass to the handler
			handler(w, r, pool, customer, conv)
		}),
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
	customer *queries.Customer,
	c *Conversation,
) {
	// return relevant conversation information
	request.Encode(w, r, c.logger, http.StatusOK, &GetConverstaionResponse{
		ConversationId: c.ID,
		Title:          c.Title,
		Messages:       c.GetMessages(),
	})
}

func getConversations(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	customer *queries.Customer,
) {
	logger := httplog.LogEntry(r.Context())
	dmodel := queries.New(pool)

	// fetch the conversations
	response, err := dmodel.GetConversationsWithCount(r.Context(), customer.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			request.Encode(w, r, &logger, http.StatusOK, []string{})
			return
		}
		logger.Error("Error getting the conversations", "error", err)
		http.Error(w, fmt.Sprintf("There was an internal issue: %s", err), http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, &logger, http.StatusOK, response)
}
