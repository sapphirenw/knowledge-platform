package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/tool"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type ragMessageType string

const (
	ragLoading ragMessageType = "loading"

	ragNewMessage        ragMessageType = "newMessage"
	ragNewConversationId ragMessageType = "newConversationId"
	ragTitleUpdate       ragMessageType = "titleUpdate"
	ragError             ragMessageType = "error"
)

/*
A RAG message is the content that is sent to the user through a websocket.
this can be many different types of messages, such as control flow (loading),
messages to and from the AI, tool calls, error throws, and so on.
*/
type ragMessage struct {
	// consistent for all messages

	MessageType ragMessageType `json:"messageType"`

	// dependent on the message type

	ChatMessage    *gollm.Message `json:"chatMessage,omitempty"`
	ConversationId string         `json:"conversationId"`
	NewTitle       string         `json:"newTitle,omitempty"`
	Error          string         `json:"error,omitempty"`
}

func newRmChatMessage(msg *gollm.Message) *ragMessage {
	return &ragMessage{
		MessageType: ragNewMessage,
		ChatMessage: msg,
	}
}

func newRmError(msg string, err error) *ragMessage {

	return &ragMessage{
		MessageType: ragError,
		Error:       fmt.Sprintf("%s: %s", msg, err.Error()),
	}
}

func writeRagResponse(
	ctx context.Context,
	logger *slog.Logger,
	conn *websocket.Conn,
	message *ragMessage,
) error {
	// encode the message
	enc, err := json.Marshal(message)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to write the message", err)
	}

	// write the message on the connection
	if err := conn.WriteJSON(enc); err != nil {
		return slogger.Error(ctx, logger, "failed to write on the websocket", err)
	}

	return nil
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

	// hold the conversation in memory for the request
	var conv *conversation.Conversation

	for {
		// read the message that was passed from the user
		// the message will always be a string, and always be of type user.
		logger.Debug("Reading message")
		mt, message, err := conn.ReadMessage()
		if err != nil {
			request.WriteWs(r.Context(), logger, conn, "failed to read the message", err)
			break
		}

		// check if the conversation exists in memory
		if conv == nil {
			// fetch the conversation
			conv, err = conversation.AutoConversation(
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
				writeRagResponse(r.Context(), logger, conn, newRmError("failed to get the conversation", err))
				break
			}

			// write the conversation id update
			if err := writeRagResponse(r.Context(), logger, conn, &ragMessage{
				MessageType:    ragNewConversationId,
				ConversationId: conv.ID.String(),
			}); err != nil {
				slogger.Error(r.Context(), logger, "failed to write the rag message", err)
				break
			}
			logger = logger.With("conversationId", conv.ID.String())
		}

		logger.Debug("Recieved message from user", "message", string(message), "type", mt)

		// create the user message
		userMessage := gollm.NewUserMessage(string(message))

		// create a transaction to run this call inside of
		tx, err := pool.Begin(r.Context())
		if err != nil {
			writeRagResponse(r.Context(), logger, conn, newRmError("failed to start the transaction", err))
			break
		}

		// process the message using the handler. This will handle all operations that can occur
		// when a user sends a message, including tool calls and writing multiple messages to the user
		// when a transaction is in process.
		if err := c.rag2MessageHandler(r.Context(), logger, tx, conn, conv, userMessage); err != nil {
			writeRagResponse(r.Context(), logger, conn, newRmError("failed to handle the message", err))
			tx.Rollback(r.Context())
			break
		}
		tx.Commit(r.Context())

		// send a request to create a title if the conversation does not have one
		if conv.Title == "Information Chat" {
			logger.Info("Creating a new title ... ")
			newTitle, err := c.createRagTitle(r.Context(), logger, pool, userMessage)
			if err != nil {
				slogger.Error(r.Context(), logger, "failed to create the title, but not closing the ws", err)
			}

			// update the conversation
			dmodel := queries.New(pool)
			_, err = dmodel.UpdateConversationTitle(r.Context(), &queries.UpdateConversationTitleParams{
				ID:    conv.ID,
				Title: newTitle,
			})
			if err != nil {
				slogger.Error(r.Context(), logger, "failed tp update the conversation title, but not closing ws", err)
			}
			conv.Title = newTitle

			// send the new title to the user
			writeRagResponse(r.Context(), logger, conn, &ragMessage{
				MessageType: ragTitleUpdate,
				NewTitle:    newTitle,
			})
		}
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
	if err := writeRagResponse(ctx, logger, conn, newRmChatMessage(completionResponse.Message)); err != nil {
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
		return c.rag2ToolCallHandler(ctx, logger, tx, conn, conv, completionResponse.Message)
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
	if err := writeRagResponse(ctx, logger, conn, newRmChatMessage(toolResponse.Message)); err != nil {
		return slogger.Error(ctx, logger, "failed to write the message", err)
	}

	// recursively run the message handler
	logger.Debug("recursively calling the rag2 message handler")
	if err := c.rag2MessageHandler(ctx, logger, tx, conn, conv, toolResponse.Message); err != nil {
		return slogger.Error(ctx, logger, "failed to recursively call the message handler", err)
	}

	return nil
}

// creates a chat title based on the passed message to this function.
// it is recommended to use the first customer message as the input to this function
func (c *Customer) createRagTitle(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	message *gollm.Message,
) (string, error) {
	// get the title creation llm
	dmodel := queries.New(db)
	response, err := dmodel.GetInteralLLM(ctx, "Title creator")
	if err != nil {
		return "", slogger.Error(ctx, logger, "failed to get the internal llm for title creation", err)
	}

	lm := llm.FromObjects(&response.Llm, &response.AvailableModel)

	// create a single completion
	completion, err := lm.SingleCompletion(ctx, logger, c.ID, prompts.RAG_TITLE_GENERATION_SYSTEM_PROMPT, message.Message)
	if err != nil {
		return "", slogger.Error(ctx, logger, "failed to send the single completion for a new title", err)
	}

	// report the usage
	if err := utils.ReportUsage(ctx, logger, db, c.ID, []*tokens.UsageRecord{completion.UsageRecord}, nil); err != nil {
		return "", slogger.Error(ctx, logger, "failed to report the usage", err)
	}

	return completion.Message.Message, nil

}

func rag2Tools() []tool.Tool {
	tools := make([]tool.Tool, 0)
	tools = append(tools, tool.NewTool(tool.VectorQuery))
	return tools
}
