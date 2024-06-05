package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Conversation struct {
	*queries.Conversation

	Messages       []*ConversationMessage
	storedMessages int // representation of how many messages are already stored in the database

	logger *slog.Logger
	New    bool // whether the conversation was created in this request or not
}

type ConversationMessage struct {
	*queries.ConversationMessage
}

func (cm *ConversationMessage) ToGoLLMMessage() (*gollm.LanguageModelMessage, error) {
	var role gollm.LanguageModelRole
	switch cm.Role {
	case gollm.RoleSystem.ToString():
		role = gollm.RoleSystem
	case gollm.RoleUser.ToString():
		role = gollm.RoleUser
	case gollm.RoleAI.ToString():
		role = gollm.RoleAI
	default:
		return nil, fmt.Errorf("invalid role in message: %s", cm.Role)
	}

	return &gollm.LanguageModelMessage{
		Role:    role,
		Message: cm.Message,
	}, nil
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
			return nil, fmt.Errorf("failed to create the conversation: %s", err)
		}
	} else {
		if _, err := uuid.Parse(conversationId); err != nil {
			return nil, fmt.Errorf("failed to parse the conversationId")
		}

		// get the existing conversation
		conv, err = GetConversation(ctx, logger, db, uuid.MustParse(conversationId))
		if err != nil {
			return nil, fmt.Errorf("failed to get the conversation: %s", err)
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
		return nil, fmt.Errorf("failed to create the conversation: %s", err)
	}

	return &Conversation{
		Conversation:   conversation,
		Messages:       make([]*ConversationMessage, 0),
		storedMessages: 0,
		logger:         logger.With("conversationId", conversation.ID.String()),
		New:            true,
	}, nil
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
		return nil, fmt.Errorf("failed to get the conversation: %s", err)
	}

	// get the messages
	msgs, err := dmodel.GetConversationMessages(ctx, conversationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %s", err)
	}

	logger.DebugContext(ctx, "Fetched messages from database", "length", len(msgs))

	// make the needed internal lists
	messages := make([]*ConversationMessage, 0)
	for _, item := range msgs {
		cm := &ConversationMessage{ConversationMessage: item}
		messages = append(messages, cm)
	}

	return &Conversation{
		Conversation:   conv,
		Messages:       messages,
		storedMessages: len(messages),
		logger:         logger,
		New:            false,
	}, nil
}

// Get the internal messages as gollm message objects
func (c *Conversation) GoLLMMessages() ([]*gollm.LanguageModelMessage, error) {
	gollmMessages := make([]*gollm.LanguageModelMessage, 0)
	for _, item := range c.Messages {
		gollmMessage, err := item.ToGoLLMMessage()
		if err != nil {
			return nil, fmt.Errorf("failed to parse gollm message: %s", err)
		}
		gollmMessages = append(gollmMessages, gollmMessage)
	}

	return gollmMessages, nil
}

func (c *Conversation) Completion(
	ctx context.Context,
	db queries.DBTX,
	model *LLM,
	input string,
) (string, error) {
	logger := c.logger.With("model", model.ID.String())

	// get the gollm object
	lm, err := c.getGoLLM(ctx, model)
	if err != nil {
		return "", fmt.Errorf("failed to create the gollm: %s", err)
	}

	logger.InfoContext(ctx, "Beginning conversation completion ...")
	response, err := model.Completion(ctx, c.logger, lm, &CompletionArgs{
		Input: input,
	})
	if err != nil {
		return "", fmt.Errorf("failed conversation completion: %s", err)
	}

	// save the conversation records in the conversation
	logger.InfoContext(ctx, "Saving the conversation ...")
	if err := c.addMessages(ctx, db, model, lm.GetConversation()); err != nil {
		return "", fmt.Errorf("failed to add the messages to the internal conversation: %s", err)
	}

	// report the usage
	logger.InfoContext(ctx, "Reporting the usage ...")
	if err := utils.ReportUsage(ctx, c.logger, db, c.CustomerID, lm.GetTokenRecords(), c.Conversation); err != nil {
		return "", fmt.Errorf("failed to save the token usage")
	}

	logger.InfoContext(ctx, "Successfully saved conversation")
	return response, nil
}

func JsonCompletion[T any](
	c *Conversation,
	ctx context.Context,
	db queries.DBTX,
	model *LLM,
	input string,
	schema string,
) (*T, error) {
	var response T
	c.logger.InfoContext(ctx, "Beginning conversation json completion ...")

	// get the gollm object
	lm, err := c.getGoLLM(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("failed to create the gollm: %s", err)
	}

	rawResponse, err := model.Completion(ctx, c.logger, lm, &CompletionArgs{
		Input:      input,
		Json:       true,
		JsonSchema: schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed conversation completion: %s", err)
	}

	// parse as json
	c.logger.InfoContext(ctx, "Parsing into json ...")
	if err := json.Unmarshal([]byte(rawResponse), &response); err != nil {
		return nil, fmt.Errorf("failed to serialize the json: %s. RAW: %s", err, rawResponse)
	}

	// save the conversation records in the conversation
	c.logger.InfoContext(ctx, "Saving the conversation ...")
	if err := c.addMessages(ctx, db, model, lm.GetConversation()); err != nil {
		return nil, fmt.Errorf("failed to add the messages to the internal conversation: %s", err)
	}

	// report the usage
	c.logger.InfoContext(ctx, "Reporting the usage ...")
	if err := utils.ReportUsage(ctx, c.logger, db, c.CustomerID, lm.GetTokenRecords(), c.Conversation); err != nil {
		return nil, fmt.Errorf("failed to save the token usage")
	}

	c.logger.InfoContext(ctx, "Successfully saved conversation")
	return &response, nil
}

// Correctly sets the system message and creates a new gollm object based on the model that is passed
func (c *Conversation) getGoLLM(ctx context.Context, model *LLM) (*gollm.LanguageModel, error) {
	c.logger.DebugContext(ctx, "Creating the internal llm with the passed configuration ...")

	// create the gollm object
	sysMessageRaw := model.GenerateSystemPrompt(c.Conversation.SystemMessage)
	gollmMessages, err := c.GoLLMMessages()
	if err != nil {
		return nil, fmt.Errorf("failed to parse the gollm messages: %s", err)
	}

	// set the first gollm message
	sysMessage := gollm.LanguageModelMessage{
		Role:    gollm.RoleSystem,
		Message: sysMessageRaw,
	}
	if len(gollmMessages) == 0 {
		gollmMessages = append(gollmMessages, &gollm.LanguageModelMessage{})
	}
	gollmMessages[0] = &sysMessage

	// construct the object
	lm := gollm.NewLanguageModelFromConversation(
		c.CustomerID.String(),
		c.logger,
		gollmMessages,
		nil,
	)

	if lm == nil {
		return nil, fmt.Errorf("failed to create the gollm: %s", err)
	}

	return lm, nil
}

// Adds a list of messages to the conversation and updates the database.
func (c *Conversation) addMessages(
	ctx context.Context,
	db queries.DBTX,
	model *LLM,
	messages []*gollm.LanguageModelMessage,
) error {
	for index := range len(messages) {
		if index < c.storedMessages {
			continue
		}

		c.Messages = append(c.Messages, &ConversationMessage{
			ConversationMessage: &queries.ConversationMessage{
				ConversationID: c.ID,
				LlmID:          utils.GoogleUUIDToPGXUUID(model.ID),
				Model:          model.Model,
				Temperature:    model.Temperature,
				Instructions:   model.Instructions,
				Role:           messages[index].Role.ToString(),
				Message:        messages[index].Message,
				Index:          int32(len(c.Messages)),
			},
		})
	}

	if err := c.SyncMessages(ctx, db); err != nil {
		return fmt.Errorf("failed to sync the messages with the database: %s", err)
	}

	return nil
}

// Replaces the internal messages array with the given list, and updates in the database.
// If this fails, then the message list does not get changed.
func (c *Conversation) replaceMessages(
	ctx context.Context,
	db queries.DBTX,
	model *LLM,
	messages []*gollm.LanguageModelMessage,
) error {
	// create a copy of the messages
	var oldMsgs []*ConversationMessage
	if r := copy(oldMsgs, c.Messages); r == -1 {
		// no idea if this is an actual thing this function can throw as I am on an airplane and
		// cannot google, but may as well check
		return fmt.Errorf("failed to copy messages")
	}

	// clear messages in the database
	dmodel := queries.New(db)
	if err := dmodel.ClearConversation(ctx, c.ID); err != nil {
		c.Messages = oldMsgs
		return fmt.Errorf("failed to clear the conversation: %s", err)
	}
	c.Messages = make([]*ConversationMessage, 0)

	// create the new message array and add the messages
	if err := c.addMessages(ctx, db, model, messages); err != nil {
		c.Messages = oldMsgs
		return fmt.Errorf("failed to add the messages: %s", err)
	}

	return nil
}

// This function syncs the internal message state to the database. Calling this function is not
// required in most cases, as the `Completion` method handles storing the conversation state,
// syncing the database, and reporting the token usage
func (c *Conversation) SyncMessages(
	ctx context.Context,
	db queries.DBTX,
) error {
	dmodel := queries.New(db)

	for index, item := range c.Messages {
		msg, err := dmodel.CreateConversationMessage(ctx, &queries.CreateConversationMessageParams{
			ConversationID: c.ID,
			LlmID:          item.LlmID,
			Model:          item.Model,
			Temperature:    item.Temperature,
			Instructions:   item.Instructions,
			Role:           item.Role,
			Message:        item.Message,
			Index:          int32(index),
		})
		if err != nil {
			return fmt.Errorf("failed to refresh message: %s", err)
		}
		c.Messages[index] = &ConversationMessage{ConversationMessage: msg}
	}
	c.storedMessages = len(c.Messages)
	return nil
}

func (c *Conversation) PrintConversation() {
	for _, item := range c.Messages {
		fmt.Printf("[[%s]]\n", item.Role)
		fmt.Printf("> %s\n---\n", item.Message)
	}
}
