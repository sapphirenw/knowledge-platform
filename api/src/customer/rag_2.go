package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/tool"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type ragMessageType string

const (
	ragError             ragMessageType = "error"
	ragNewMessage        ragMessageType = "newMessage"
	ragNewConversationId ragMessageType = "newConversationId"
	ragTitleUpdate       ragMessageType = "titleUpdate"
	ragChangeChatLLM     ragMessageType = "changeChatLLM"
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

	Error          string         `json:"error,omitempty"`
	ChatMessage    *gollm.Message `json:"chatMessage,omitempty"`
	ConversationId string         `json:"conversationId"`
	NewTitle       string         `json:"newTitle,omitempty"`
	ChatLLM        *llm.LLM       `json:"chatLLM,omitempty"`
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

type ragUserMessage struct {
	MessageType string `json:"messageType"`

	Message   string `json:"message"`
	ChatLLMID string `json:"chatLLMId"`
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

	dmodel := queries.New(pool)

	var chatLLM *llm.LLM

	// check if a conversation id was passed, and attempt to fetch the conversation
	convId, err := utils.GoogleUUIDFromString(r.URL.Query().Get("id"))
	if err == nil {
		// get the conversation
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
			logger.Error("failed to get the conversation", "error", err)
		} else {
			// get the chat llm as well
			currLLM, err := dmodel.GetChatLLM(r.Context(), &queries.GetChatLLMParams{
				CustomerID: c.ID,
				ID:         convId,
			})
			if err == nil {
				// set the current chat llm
				logger.Debug("the conversation has a saved llm, using this", "llm.ID", currLLM.Llm.ID)
				chatLLM = llm.FromObjects(&currLLM.Llm, &currLLM.AvailableModel)
			} else {
				if strings.Contains(err.Error(), "no rows in result set") {
					logger.Info("No saved llm exists")
				} else {
					logger.Warn("unknow error getting the chatllm", "error", err)
				}
			}
		}
	} else {
		logger.Debug("invalid conversation id", "conv.ID", r.URL.Query().Get("id"))
	}

	// fetch a default llm to use if no conversation llm was found
	if chatLLM == nil {
		chatLLM, err = c.GetChatLLM(r.Context(), logger, pool)
		if err != nil {
			writeRagResponse(r.Context(), logger, conn, newRmError("failed to get the llm", err))
			return
		}
	}

	// send a message to the user populating the chatllm object
	writeRagResponse(r.Context(), logger, conn, &ragMessage{
		MessageType: ragChangeChatLLM,
		ChatLLM:     chatLLM,
	})

	for {
		// read the message that was passed from the user
		// the message will always be a string, and always be of type user.
		logger.Debug("Reading message")
		mt, message, err := conn.ReadMessage()
		if err != nil {
			writeRagResponse(r.Context(), logger, conn, newRmError("failed to read the message", err))
			break
		}

		logger.Debug("Recieved message from user", "message", string(message), "type", mt)

		// parse the user message
		var userMessage ragUserMessage
		if err := json.Unmarshal(message, &userMessage); err != nil {
			writeRagResponse(r.Context(), logger, conn, newRmError("failed to read the user message", err))
			break
		}

		// parse the user request
		logger.Info("Handling user message", "type", userMessage.MessageType)

		// handle chanding chat llm
		if userMessage.MessageType == "changeChatLLM" {
			// parse the id
			chatLLMId, err := utils.GoogleUUIDFromString(userMessage.ChatLLMID)
			if err != nil {
				writeRagResponse(r.Context(), logger, conn, newRmError("invalid message id", err))
				// do not exit the connection
			} else {
				// fetch the llm
				dmodel := queries.New(pool)
				fetchedLLM, err := dmodel.GetLLM(r.Context(), chatLLMId)
				if err != nil {
					writeRagResponse(r.Context(), logger, conn, newRmError("failed to get the model", err))
				} else {
					logger.Info("Setting new chat llm", "llmid", fetchedLLM.Llm.ID.String())
					chatLLM = llm.FromObjects(&fetchedLLM.Llm, &fetchedLLM.AvailableModel)

					// send the message to the client
					writeRagResponse(r.Context(), logger, conn, &ragMessage{
						MessageType: ragChangeChatLLM,
						ChatLLM:     chatLLM,
					})

					// save the chatllm to the conversation (THIS CAN TRANSIENTLY FAIL)
					if err := dmodel.SetChatLLM(r.Context(), &queries.SetChatLLMParams{
						ID:        conv.ID,
						CurrLlmID: utils.GoogleUUIDToPGXUUID(chatLLM.Llm.ID),
					}); err != nil {
						logger.Error("failed to set the chatllm", "error", err)
					}
				}
			}
		}

		// handle the rag message
		if userMessage.MessageType == "ragMessage" {
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

			// create a transaction to run this call inside of
			tx, err := pool.Begin(r.Context())
			if err != nil {
				writeRagResponse(r.Context(), logger, conn, newRmError("failed to start the transaction", err))
				break
			}

			// send the request
			if err := c.rag2MessageHandler(r.Context(), logger, tx, conn, conv, chatLLM, gollm.NewUserMessage(userMessage.Message)); err != nil {
				tx.Rollback(r.Context())
				writeRagResponse(r.Context(), logger, conn, newRmError("failed to send the message request", err))
				break
			}

			tx.Commit(r.Context())

			// send a request to create a title if the conversation does not have one
			if conv.Title == "Information Chat" {
				logger.Info("Creating a new title ... ")
				newTitle, err := c.createRagTitle(r.Context(), logger, pool, gollm.NewUserMessage(userMessage.Message))
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
	chatLLM *llm.LLM,
	message *gollm.Message,
) error {
	// get the tools
	tools := rag2Tools()

	// get the customers chatllm
	// TODO -- enable arguments to be passed over the websocket
	logger.Debug("chatllm", "chatllm", *chatLLM)

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
	summaryLLM, err := c.GetSummaryLLM(ctx, logger, tx)
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
	if err := c.rag2MessageHandler(ctx, logger, tx, conn, conv, chatLLM, toolResponse.Message); err != nil {
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
