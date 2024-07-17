package conversation

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"
)

func TestConversation(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.GetDefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)

	// get a default llm
	model, err := llm.GetLLM(ctx, pool, c.ID, pgtype.UUID{})
	require.NoError(t, err)

	conv, err := CreateConversation(ctx, logger, pool, c.ID, "You are a pirate", "Test Conversation", "Testing")
	require.NoError(t, err)

	// create a completion event
	userMsg := gollm.NewUserMessage("Ahoy Matey!!")
	response, err := conv.Completion(ctx, pool, model, userMsg, nil, nil)
	require.NoError(t, err)

	fmt.Printf("MODEL RESPONSE: %s\n", response.Message.Message)

	// check the records against the database
	dmodel := queries.New(pool)

	// check the stored conversation
	msgs, err := dmodel.GetConversationMessages(ctx, conv.ID)
	require.NoError(t, err)
	require.Equal(t, 3, len(msgs))

	// check the saved token records
	records, err := dmodel.GetCustomerTokenUsage(ctx, &queries.GetCustomerTokenUsageParams{
		CustomerID: c.ID,
		Limit:      20,
		Column4:    1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(records))
	require.Equal(t, conv.ID.String(), utils.PGXUUIDToGoogleUUID(records[0].ConversationID).String())
}
