package customer

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"

	_ "net/http/pprof"
)

func TestRag(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	ctx, logger, pool, c := testInit(t)

	response1, err := c.RAG(ctx, pool, &ragRequest{
		ConversationId: "",
		Input:          "What is the meaning of the world",
	})
	require.NoError(t, err)
	fmt.Println(response1.Response)

	response2, err := c.RAG(ctx, pool, &ragRequest{
		ConversationId: response1.ConverationId,
		Input:          "Are you sure",
	})
	require.NoError(t, err)
	fmt.Println(response2.Response)

	dmodel := queries.New(pool)

	// check the conversation
	convId := utils.PGXUUIDToGoogleUUID(utils.PGXUUIDFromString(response2.ConverationId))
	require.NotNil(t, convId)
	require.NoError(t, err)
	conv, err := llm.GetConversation(ctx, logger, pool, *convId)
	require.NoError(t, err)
	require.Equal(t, 5, len(conv.Messages))

	// check the usage records
	usage, err := dmodel.GetTokenUsage(ctx, c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, usage)

	conv.PrintConversation()
}
