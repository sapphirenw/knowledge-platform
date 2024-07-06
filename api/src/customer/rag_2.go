package customer

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/tool"
)

type RagMessageType string

const (
	Error             RagMessageType = "error"
	Success           RagMessageType = "success"
	Loading           RagMessageType = "loading"
	NewConversationId RagMessageType = "newConversationId"
	TitleUpdate       RagMessageType = "titleUpdate"
)

/*
A RAG message is the content that is sent to the user through a websocket.
this can be many different types of messages, such as control flow (loading),
messages to and from the AI, tool calls, error throws, and so on.
*/
type RagMessage struct {
	MessageType RagMessageType `json:"messageType"`
	Error       error          `json:"error"`
	Status      string         `json:"status"`
	ChatMessage *gollm.Message `json:"chatMessage"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRag2(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	logger := c.logger.With("handler", "rag2")

	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slogger.ServerError(w, logger, 500, "failed to upgrade the connection", err)
		return
	}
	defer conn.Close()

	logger.Info("Opened ws connection")

	// get the conversation
	// fetch or create a conversation
	conv, err := conversation.AutoConversation(
		r.Context(),
		logger,
		pool,
		c.ID,
		r.URL.Query().Get("id"),
		prompts.RAG_COMPLETE_SYSTEM_PROMPT,
		"Information Chat",
		"rag",
	)
	if err != nil {
		request.WriteWs(r.Context(), logger, conn, "failed to get the converation", err)
		return
	}
	// conn.WriteMessage(websocket.TextMessage, []byte("Successfully got conversation"))
	logger = logger.With("conversationId", conv.ID.String())

	for {
		// read the message that was passed from the user
		// the message will always be a string, and always be of type user.
		logger.Debug("Reading message")
		mt, message, err := conn.ReadMessage()
		if err != nil {
			request.WriteWs(r.Context(), logger, conn, "failed to read the message", err)
			break
		}

		logger.Debug("Recived message from user", "message", string(message), "type", mt)

		// create the user message
		userMessage := gollm.NewUserMessage(string(message))

		// create a transaction to run this call inside of
		tx, err := pool.Begin(r.Context())
		if err != nil {
			request.WriteWs(r.Context(), logger, conn, "failed to begin the transaction", err)
			break
		}

		// process the message using the handler. This will handle all operations that can occur
		// when a user sends a message, including tool calls and writing multiple messages to the user
		// when a transaction is in process.
		if err := c.rag2MessageHandler(r.Context(), logger, tx, conn, conv, userMessage); err != nil {
			request.WriteWs(r.Context(), logger, conn, "failed to handle the message", err)
			tx.Rollback(r.Context())
			break
		}
		tx.Commit(r.Context())
	}

	logger.Info("Closing ws connection")
}

// Handles the initial message recieved from the user. This will either write the AI
// response to the user, or it will perform the tool call chain.
func (c *Customer) rag2MessageHandler(
	ctx context.Context,
	logger *slog.Logger,
	tx pgx.Tx,
	conn *websocket.Conn,
	conv *conversation.Conversation,
	message *gollm.Message,
) error {
	// get the tools
	tools := rag2Tools()

	// get the llm as required from the user
	// TODO -- handle making this dynamic from the user
	logger.Debug("Getting the chat llm ...")
	var chatLLMId pgtype.UUID
	chatLLMId.Scan("")
	chatLLM, err := llm.GetLLM(ctx, tx, c.ID, chatLLMId)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to get the chatllm", err)
	}

	// create a completion
	logger.Debug("sending a completion response in the rag handler")
	completionResponse, err := conv.Completion(ctx, tx, chatLLM, message, tool.ToolsToGollm(tools), nil)
	if err != nil {
		return slogger.Error(ctx, logger, "failed the completion", err)
	}

	// write the message to the connection
	if err := request.WriteWs(ctx, logger, conn, completionResponse.Message, nil); err != nil {
		return slogger.Error(ctx, logger, "failed to write the message", err)
	}

	// parse the message response
	logger.Debug("rag completion response role", "role", completionResponse.Message.Role.ToString())
	switch completionResponse.Message.Role {
	case gollm.RoleAI:
		// response from the AI, no futher work to do in this chain.
		return nil
	case gollm.RoleToolCall:
		// perform the tool call chain
		logger.Debug("calling the rag2 tool handler")
		return c.rag2ToolCallHandler(ctx, logger, tx, conn, conv, chatLLM, completionResponse.Message)
	default:
		return slogger.Error(ctx, logger, "unexpected message role from the AI", nil, "role", completionResponse.Message.Role.ToString())
	}
}

func (c *Customer) rag2ToolCallHandler(
	ctx context.Context,
	logger *slog.Logger,
	tx pgx.Tx,
	conn *websocket.Conn,
	conv *conversation.Conversation,
	chatLLM *llm.LLM,
	message *gollm.Message,
) error {

	// parse the tool call
	toolType, err := tool.GetToolType(message.ToolName)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to parse the tool name", err)
	}
	parsedTool := tool.NewTool(toolType)

	// get the summary llm
	// TODO -- handle the summary model
	summaryLLM, err := llm.GetLLMString(ctx, tx, c.ID, "")
	if err != nil {
		return slogger.Error(ctx, logger, "failed to get the summary llm", err)
	}

	logger.Debug("running the rag2 parsed tool run call")
	toolResponse, err := parsedTool.Run(ctx, logger, &tool.RunToolArgs{
		Database:    tx,
		Customer:    c.Customer,
		LastMessage: message,
		ToolLLM:     summaryLLM,
	})
	if err != nil {
		return slogger.Error(ctx, logger, "failed to run the tool", err)
	}

	// write the message to the connection
	if err := request.WriteWs(ctx, logger, conn, toolResponse.Message, nil); err != nil {
		return slogger.Error(ctx, logger, "failed to write the message", err)
	}

	// recursively run the message handler
	logger.Debug("recursively calling the rag2 message handler")
	if err := c.rag2MessageHandler(ctx, logger, tx, conn, conv, toolResponse.Message); err != nil {
		return slogger.Error(ctx, logger, "failed to recursively call the message handler", err)
	}

	return nil
}

func rag2Tools() []tool.Tool {
	tools := make([]tool.Tool, 0)
	tools = append(tools, tool.NewTool(tool.VectorQuery))
	return tools
}
