package conversation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Conversation struct {
	*queries.Conversation
	messages     []*gollm.Message      // internal message conversation stored
	usageRecords []*tokens.UsageRecord // token records stored with this conversation. Ephemeral, not sourced from the database on re-load
	logger       *slog.Logger
	New          bool // whether the conversation was created in this request or not
}

func GoLLMMessageFromDB(cm *queries.ConversationMessage) *gollm.Message {
	msg := &gollm.Message{
		Message: cm.Message,
	}

	// okay for this to fail

	switch cm.Role {
	case gollm.RoleSystem.ToString():
		msg.Role = gollm.RoleSystem
	case gollm.RoleUser.ToString():
		msg.Role = gollm.RoleUser
	case gollm.RoleAI.ToString():
		msg.Role = gollm.RoleAI
	case gollm.RoleToolCall.ToString():
		// get the function call arguments
		msg.Role = gollm.RoleToolCall
		msg.ToolUseID = cm.ToolUseID
		msg.ToolName = cm.ToolName

		// parse the arguments
		var args map[string]any
		json.Unmarshal(cm.ToolArguments, &args)
		msg.ToolArguments = args
	case gollm.RoleToolResult.ToString():
		msg.Role = gollm.RoleToolResult
		msg.ToolUseID = cm.ToolUseID
		msg.ToolName = cm.ToolName

		// parse the results
		var args map[string]any
		json.Unmarshal(cm.ToolResults, &args)
		msg.ToolArguments = args
	}

	return msg
}

// Attemps to parse the conversationId passed to it and fetch a conversation.
// If no conversation exists, then a new one will be created and returned.
// No conversation is created on any errors
func AutoConversation(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	customerId uuid.UUID,
	conversationId string,
	systemMessage string,
	title string,
	conversationType string,
) (*Conversation, error) {
	var conv *Conversation
	var err error
	if conversationId == "" {
		// create a new conversation
		conv, err = CreateConversation(ctx, logger, db, customerId, systemMessage, title, conversationType)
		if err != nil {
			return nil, fmt.Errorf("failed to create the conversation: %w", err)
		}
	} else {
		if _, err := uuid.Parse(conversationId); err != nil {
			return nil, fmt.Errorf("failed to parse the conversationId: '%s'", conversationId)
		}

		// get the existing conversation
		conv, err = GetConversation(ctx, logger, db, uuid.MustParse(conversationId))
		if err != nil {
			return nil, fmt.Errorf("failed to get the conversation: %w", err)
		}
	}

	return conv, err
}

func CreateConversation(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	customerId uuid.UUID,
	systemMessage string,
	title string,
	conversationType string,
) (*Conversation, error) {
	if title == "" {
		title = "(Untitled Conversation)"
	}

	dmodel := queries.New(db)
	conversation, err := dmodel.CreateConversation(ctx, &queries.CreateConversationParams{
		CustomerID:       customerId,
		Title:            title,
		ConversationType: conversationType,
		SystemMessage:    systemMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the conversation: %w", err)
	}

	// create the conversation record
	conv := &Conversation{
		Conversation: conversation,
		messages:     make([]*gollm.Message, 0),
		logger:       logger.With("conversationId", conversation.ID.String()),
		New:          true,
	}

	// Add the system message
	if err := conv.SaveMessage(ctx, db, nil, gollm.NewSystemMessage(systemMessage)); err != nil {
		return nil, fmt.Errorf("failed to sync the messages: %w", err)
	}

	return conv, nil
}

// Fetches a conversation and messages from a given conversationID
func GetConversation(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	conversationId uuid.UUID,
) (*Conversation, error) {
	dmodel := queries.New(db)

	// get the conversation
	conv, err := dmodel.GetConversation(ctx, conversationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the conversation: %w", err)
	}

	// get the messages
	msgs, err := dmodel.GetConversationMessages(ctx, conversationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %w", err)
	}

	logger.DebugContext(ctx, "Fetched messages from database", "length", len(msgs))

	// make the needed internal lists
	messages := make([]*gollm.Message, 0)
	for _, item := range msgs {
		messages = append(messages, GoLLMMessageFromDB(item))
	}

	return &Conversation{
		Conversation: conv,
		messages:     messages,
		logger:       logger,
		New:          false,
	}, nil
}

// Contains a JSON argument that should not be exposed if not necessary for normal conversations.
// The `Completion` function calls this with an empty schema.
func (c *Conversation) internalCompletion(
	ctx context.Context,
	db queries.DBTX,
	model *llm.LLM,
	message *gollm.Message,
	tools []*gollm.Tool,
	requiredTool *gollm.Tool,
	schema string,
) (*gollm.CompletionResponse, error) {
	logger := c.logger.With("model", model.Llm.ID.String())

	// create a copy of the messages array
	messages := make([]*gollm.Message, len(c.messages))
	copy(messages, c.messages)
	if message != nil {
		// add the passed message
		messages = append(messages, message)
	}

	// check the conversation state for mismatched state
	if messages[len(messages)-1].Role != gollm.RoleToolResult && messages[len(messages)-1].Role != gollm.RoleUser {
		return nil, fmt.Errorf("lastest message role: %s", messages[len(messages)-1].Role.ToString())
	}

	// run the completion
	logger.InfoContext(ctx, "Beginning conversation completion ...")
	response, err := model.Completion(ctx, c.logger, &llm.CompletionArgs{
		CustomerID:   c.CustomerID.String(),
		Messages:     messages,
		Tools:        tools,
		RequiredTool: requiredTool,
		Json:         schema != "",
		JsonSchema:   schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed conversation completion: %w", err)
	}

	// report the usage
	c.usageRecords = append(c.usageRecords, response.UsageRecord)
	logger.InfoContext(ctx, "Reporting the usage ...")
	if err := utils.ReportUsage(ctx, c.logger, db, c.CustomerID, c.usageRecords, c.Conversation); err != nil {
		c.messages = c.messages[:len(c.messages)-1]
		return nil, fmt.Errorf("failed to save the token usage")
	}

	// add messages to the conversation
	if message != nil {
		logger.InfoContext(ctx, "Saving the input message ...")
		if err := c.SaveMessage(ctx, db, model, message); err != nil {
			c.messages = c.messages[:len(c.messages)-1]
			return nil, fmt.Errorf("failed to save the input message to the conversation: %w", err)
		}
	}
	logger.InfoContext(ctx, "Saving the output message ...")
	if err := c.SaveMessage(ctx, db, model, response.Message); err != nil {
		c.messages = c.messages[:len(c.messages)-1]
		return nil, fmt.Errorf("failed to save the output message to the conversation: %w", err)
	}

	logger.InfoContext(ctx, "Successfully saved conversation")
	return response, nil
}

// Runs a completion against the model, and automatically saves the response message into the
// database
func (c *Conversation) Completion(
	ctx context.Context,
	db queries.DBTX,
	model *llm.LLM,
	message *gollm.Message,
	tools []*gollm.Tool,
	requiredTool *gollm.Tool,
) (*gollm.CompletionResponse, error) {
	response, err := c.internalCompletion(ctx, db, model, message, tools, requiredTool, "")
	if err != nil {
		if err := c.ReportError(ctx, db, err); err != nil {
			return nil, slogger.Error(ctx, c.logger, "failed to report the internal error for the convertation", err)
		}
		return nil, slogger.Error(ctx, c.logger, "failed the internal completion on the converstaion", err)
	}
	return response, nil
}

// Send a JSON completion against the model where the response is automatically serialized
// from the response message. This function calls `Conversation.Completion` under the hood.
// Note: This will not response with the entire response object as seen in Completion. Ensure
// there is no information in this object that you need.
func JsonCompletion[T any](
	conv *Conversation,
	ctx context.Context,
	db queries.DBTX,
	model *llm.LLM,
	message *gollm.Message,
	tools []*gollm.Tool,
	schema string,
) (*T, error) {
	// create a completion
	response, err := jsonCompletion[T](conv, ctx, db, model, message, tools, schema)
	if err != nil {
		// report the error on the conversation
		if err := conv.ReportError(ctx, db, err); err != nil {
			return nil, slogger.Error(ctx, conv.logger, "failed to report the internal error for the convertation", err)
		}
		return nil, err
	}
	return response, nil
}

func jsonCompletion[T any](
	conv *Conversation,
	ctx context.Context,
	db queries.DBTX,
	model *llm.LLM,
	message *gollm.Message,
	tools []*gollm.Tool,
	schema string,
) (*T, error) {
	// check the emssage
	if message.Role == gollm.RoleToolCall || message.Role == gollm.RoleToolResult {
		return nil, fmt.Errorf("the role cannot be a tool result or response to use JSON mode")
	}

	// create a completion
	response, err := conv.internalCompletion(ctx, db, model, message, tools, nil, schema)
	if err != nil {
		return nil, err
	}

	// serialize the response
	var resp T
	if err := json.Unmarshal([]byte(response.Message.Message), &resp); err != nil {
		return nil, fmt.Errorf("failed to serialize the JSON: %w", err)
	}

	return &resp, nil
}

// Adds a message to the internal messages array and saves the messages to the database
func (c *Conversation) SaveMessage(
	ctx context.Context,
	db queries.DBTX,
	model *llm.LLM,
	message *gollm.Message,
) error {
	input := &queries.CreateConversationMessageParams{
		ConversationID: c.ID,
		Role:           message.Role.ToString(),
		Message:        message.Message,
		Index:          int32(len(c.messages)),
		ToolUseID:      message.ToolUseID,
		ToolName:       message.ToolName,
	}

	if model != nil {
		input.LlmID = utils.GoogleUUIDToPGXUUID(model.Llm.ID)
		input.Model = model.Llm.Model
		input.Temperature = model.Llm.Temperature
		input.Instructions = model.Llm.Instructions
	}

	// add the arguments if valid
	if message.Role == gollm.RoleToolCall {
		enc, err := json.Marshal(message.ToolArguments)
		if err != nil {
			return fmt.Errorf("failed encode the tool arguments: %w", err)
		}
		input.ToolArguments = enc
	}
	if message.Role == gollm.RoleToolResult {
		enc, err := json.Marshal(message.ToolArguments)
		if err != nil {
			return slogger.Error(ctx, nil, "failed to envode the tool results", err)
		}
		input.ToolResults = enc
	}

	// post to the database
	dmodel := queries.New(db)
	_, err := dmodel.CreateConversationMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to save the message: %w", err)
	}

	// add the message to the internal array
	c.messages = append(c.messages, message)
	return nil
}

func (c *Conversation) ReportError(
	ctx context.Context,
	db queries.DBTX,
	err error,
) error {
	// update the conversation to have an error
	var msg string
	if err == nil {
		msg = "Unknown Error"
	} else {
		msg = err.Error()
	}
	dmodel := queries.New(db)
	_, err = dmodel.SetConversationError(ctx, &queries.SetConversationErrorParams{
		ID:           c.ID,
		ErrorMessage: &msg,
	})
	if err != nil {
		return fmt.Errorf("failed to report the internal error to the database: %w", err)
	}
	return nil
}

// return copies
func (c Conversation) GetMessages() []*gollm.Message {
	return c.messages
}

func (c *Conversation) PrintConversation() {
	gollm.PrintConversation(c.messages)
}
