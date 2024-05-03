package embeddings

import (
	"context"
	"testing"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/stretchr/testify/require"
)

func TestOpenAIEmbeddings(t *testing.T) {
	ctx := context.TODO()
	input := "Hello world, this is a string that I am going to convert into an embedding!"

	// get a testing database
	pool := testingutils.GetDatabase(t, ctx)

	// create a test customer
	customer := testingutils.CreateTestCustomer(t, ctx, pool)
	require.NotNil(t, customer)

	// send the embeddings request
	embeddings := NewOpenAIEmbeddings(customer.ID, nil)
	response, err := embeddings.Embed(ctx, input)
	require.Nil(t, err)
	if err != nil {
		return
	}

	require.Equal(t, 1, len(response))

	// insert the embeddings and ensure that no dupliates get created
	for range 3 {
		err = embeddings.ReportUsage(ctx, pool)
		require.Nil(t, err)

		// test that the correct data exists in the database
		model := queries.New(pool)
		usage, err := model.GetTokenUsage(ctx, customer.ID)
		require.Nil(t, err)
		if err != nil {
			return
		}
		require.Equal(t, 1, len(usage))
	}
}
