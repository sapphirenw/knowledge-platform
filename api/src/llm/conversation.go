package llm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Conversation struct {
	*queries.Conversation

	LanguageModel  *gollm.LanguageModel
	Messages       []*ConversationMessage
	storedMessages int // representation of how many messages are already stored in the database

	logger *slog.Logger
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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the conversation: %s", err)
	}

	l := logger.With("conversationId", conversation.ID.String())

	// create the gollm object
	lm := gollm.NewLanguageModel(customerId.String(), l, systemMessage, nil)

	return &Conversation{
		Conversation:   conversation,
		LanguageModel:  lm,
		Messages:       make([]*ConversationMessage, 0),
		storedMessages: 0,
		logger:         logger.With("conversationId", conversation.ID.String()),
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
	gollmMessages := make([]*gollm.LanguageModelMessage, 0)
	for _, item := range msgs {
		cm := &ConversationMessage{ConversationMessage: item}
		messages = append(messages, cm)

		// parse the gollm message
		gollmMessage, err := cm.ToGoLLMMessage()
		if err != nil {
			return nil, fmt.Errorf("failed to parse gollm message: %s", err)
		}
		gollmMessages = append(gollmMessages, gollmMessage)
	}

	l := logger.With("conversationId", conv.ID.String())
	lm := gollm.NewLanguageModelFromConversation(conv.CustomerID.String(), l, gollmMessages, nil)

	return &Conversation{
		Conversation:   conv,
		Messages:       messages,
		LanguageModel:  lm,
		storedMessages: len(messages),
		logger:         l,
	}, nil
}

func (c *Conversation) Completion(
	ctx context.Context,
	db queries.DBTX,
	model *LLM,
	args *CompletionArgs,
) (string, error) {
	c.logger.InfoContext(ctx, "Beginning conversation completion ...")
	response, err := model.Completion(ctx, c.logger, c.LanguageModel, args)
	if err != nil {
		return "", fmt.Errorf("failed conversation completion: %s", err)
	}

	// save the conversation records in the conversation
	c.logger.InfoContext(ctx, "Saving the conversation ...")
	if err := c.addMessages(ctx, db, model.Llm, c.LanguageModel.GetConversation()); err != nil {
		return "", fmt.Errorf("failed to add the messages to the internal conversation: %s", err)
	}

	// report the usage
	c.logger.InfoContext(ctx, "Reporting the usage ...")
	if err := utils.ReportUsage(ctx, c.logger, db, c.CustomerID, c.LanguageModel.GetTokenRecords(), c.Conversation); err != nil {
		return "", fmt.Errorf("failed to save the token usage")
	}

	c.logger.InfoContext(ctx, "Successfully saved conversation")
	return response, nil
}

// Adds a list of messages to the conversation and updates the database.
func (c *Conversation) addMessages(
	ctx context.Context,
	db queries.DBTX,
	model *queries.Llm,
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
	model *queries.Llm,
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
