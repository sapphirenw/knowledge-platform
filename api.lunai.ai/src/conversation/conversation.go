package conversation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Conversation struct {
	*queries.Conversation

	Messages []*ConversationMessage
}

type ConversationMessage struct {
	*queries.ConversationMessage
}

// Creates a conversation from a given model, title, and list of llm messages. Creates the
// objects in the database and returns the object with the messages.
func CreateConversation(
	ctx context.Context,
	db queries.DBTX,
	c *queries.Customer,
	model *queries.Llm,
	title string,
	messages []*gollm.LanguageModelMessage,
) (*Conversation, error) {
	if title == "" {
		title = "(Untitled Conversation)"
	}

	dmodel := queries.New(db)
	conv, err := dmodel.CreateConversation(ctx, &queries.CreateConversationParams{
		CustomerID: c.ID,
		Title:      title,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %s", err)
	}

	conversation := &Conversation{Conversation: conv}

	// add the messages
	if err := conversation.AddMessages(ctx, model, messages); err != nil {
		return nil, fmt.Errorf("failed to add the messages: %s", err)
	}

	// sync the messages
	if err := conversation.SyncMessages(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to sync the messages: %s", err)
	}

	return conversation, nil
}

// Fetches a conversation and messages from a given conversationID
func GetConversation(
	ctx context.Context,
	db queries.DBTX,
	conversationId uuid.UUID,
) (*Conversation, error) {
	dmodel := queries.New(db)
	conv, err := dmodel.GetConversation(ctx, conversationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the conversation: %s", err)
	}

	msgs, err := dmodel.GetConversationMessages(ctx, conversationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation messages: %s", err)
	}

	messages := make([]*ConversationMessage, 0)
	for _, item := range msgs {
		messages = append(messages, &ConversationMessage{ConversationMessage: item})
	}

	return &Conversation{
		Conversation: conv,
		Messages:     messages,
	}, nil
}

// Adds a message to the internal messages list, but does NOT post to the database
func (c *Conversation) AddMessage(
	ctx context.Context,
	model *queries.Llm,
	message *gollm.LanguageModelMessage,
) error {
	if message == nil {
		return fmt.Errorf("the message cannot be nil")
	}

	c.Messages = append(c.Messages, &ConversationMessage{
		ConversationMessage: &queries.ConversationMessage{
			ConversationID: c.ID,
			LlmID:          utils.GoogleUUIDToPGXUUID(model.ID),
			Model:          model.Model,
			Temperature:    model.Temperature,
			Instructions:   model.Instructions,
			Role:           message.Role.ToString(),
			Message:        message.Message,
			Index:          int32(len(c.Messages)),
		},
	})

	return nil
}

// Adds a list of messages, but doe NOT post to the database
func (c *Conversation) AddMessages(
	ctx context.Context,
	model *queries.Llm,
	messages []*gollm.LanguageModelMessage,
) error {
	for _, item := range messages {
		if err := c.AddMessage(ctx, model, item); err != nil {
			return fmt.Errorf("issue adding message: %s", err)
		}
	}

	return nil
}

// Replaces the internal messages array with the given list, but does NOT post to the
// database
func (c *Conversation) ReplaceMessages(
	ctx context.Context,
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
	c.Messages = make([]*ConversationMessage, 0)
	if err := c.AddMessages(ctx, model, messages); err != nil {
		c.Messages = oldMsgs
		return fmt.Errorf("failed to add the messages: %s", err)
	}
	return nil
}

// Syncs whatever is inside the internal messages array with the given list.
// TODO -- Deletes the messages associated with the conversation and re-creates them. (may not need)
func (c *Conversation) SyncMessages(
	ctx context.Context,
	db queries.DBTX,
) error {
	dmodel := queries.New(db)

	for index, item := range c.Messages {
		if _, err := dmodel.CreateConversationMessage(ctx, &queries.CreateConversationMessageParams{
			ConversationID: c.ID,
			LlmID:          item.LlmID,
			Model:          item.Model,
			Temperature:    item.Temperature,
			Instructions:   item.Instructions,
			Role:           item.Role,
			Message:        item.Message,
			Index:          int32(index),
		}); err != nil {
			return fmt.Errorf("failed to refresh message: %s", err)
		}
	}
	return nil
}
