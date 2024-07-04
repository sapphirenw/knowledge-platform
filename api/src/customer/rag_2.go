package customer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type ragRequestResponse struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"path"`
}

func handleRag2Init(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "rag2Init")

	// create a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to start a transaction", err)
		return
	}
	defer tx.Commit(r.Context())

	// fetch or create a conversation
	conv, err := conversation.AutoConversation(
		r.Context(),
		logger,
		tx,
		c.ID,
		r.URL.Query().Get("id"),
		prompts.RAG_COMPLETE_SYSTEM_PROMPT,
		"Information Chat",
		"rag",
	)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to create a conversation", err)
		return
	}

	// compose the path that the client will use to connect
	ragPath := fmt.Sprintf("v1/customers/%s/rag2?id=%s", c.ID.String(), conv.ID)

	// send the request to the client
	request.Encode(w, r, logger, 200, ragRequestResponse{
		ID:   conv.ID,
		Path: ragPath,
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func rag2Handler(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "rag2")

	// parse the id
	id, err := utils.GoogleUUIDFromString(r.URL.Query().Get("id"))
	if err != nil {
		slogger.ServerError(w, logger, 400, "invalid id", err)
		return
	}

	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to upgrade the connection", err)
		return
	}
	defer conn.Close()

	logger.Info("Opened ws connection")

	// get the conversation
	conv, err := conversation.GetConversation(r.Context(), logger, pool, id)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("failed to get the conversation"))
		slogger.ServerError(w, logger, 500, "failed to get the conversation", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Successfully got conversation"))

	logger.Info(conv.Title)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			slogger.Error(r.Context(), logger, "failed to read the message", err)
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, []byte(fmt.Sprintf("From server: %s", message)))
		if err != nil {
			slogger.Error(r.Context(), logger, "failed to write the message", err)
			break
		}
	}

	logger.Info("Closing ws connection")
}
