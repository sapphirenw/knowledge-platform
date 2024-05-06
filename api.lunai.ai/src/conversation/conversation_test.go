package conversation

import (
	"context"
	"testing"

	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/stretchr/testify/require"
)

func TestConversation(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.GetDefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)

	// create a model that was used for this conversation
	model, err := llm.CreateLLM(ctx, logger, pool, c, "Default Model", "gpt-3.5-turbo", 1.0, "You are a friendly and helpful AI assistant", false)
	require.NoError(t, err)

	// create some fake messages
	msgs := []*gollm.LanguageModelMessage{
		{
			Role:    gollm.RoleSystem,
			Message: model.GenerateSystemPrompt("You are a priate, talk as such"),
		},
		{
			Role:    gollm.RoleUser,
			Message: "Hello, world",
		},
	}

	conv, err := CreateConversation(ctx, pool, c, model.Llm, "New conv", msgs)
	require.NoError(t, err)
	require.Equal(t, 2, len(conv.Messages))

	// add a new message
	err = conv.AddMessage(ctx, model.Llm, &gollm.LanguageModelMessage{Role: gollm.RoleAI, Message: "Ahoy matey"})
	require.NoError(t, err)

	// sync the messages
	err = conv.SyncMessages(ctx, pool)
	require.NoError(t, err)

	// replace the messages
	newMsgs := make([]*gollm.LanguageModelMessage, 0)
	err = conv.ReplaceMessages(ctx, model.Llm, newMsgs)
	require.NoError(t, err)

	// get the conversation
	newConv, err := GetConversation(ctx, pool, conv.ID)
	require.NoError(t, err)
	require.Equal(t, conv.ID, newConv.ID)
	require.Equal(t, 3, len(newConv.Messages))
	require.NotEqual(t, len(conv.Messages), len(newConv.Messages))
}
