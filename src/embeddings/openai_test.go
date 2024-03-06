package embeddings

import (
	"context"
	"testing"

	db "github.com/sapphirenw/ai-content-creation-api/src/database"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/stretchr/testify/assert"
)

func TestOpenAIEmbeddings(t *testing.T) {
	ctx := context.TODO()
	input := "Hello world, this is a string that I am going to convert into an embedding!"

	// create the database
	pool, err := db.GetPool()
	assert.Nil(t, err)
	if err != nil {
		return
	}

	// start a transaction
	txn, err := pool.Begin(ctx)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	defer txn.Rollback(ctx) // do not commit the changes

	// create a test customer
	customer, err := testingutils.CreateTestCustomer(txn)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	// send the embeddings request
	embeddings := NewOpenAIEmbeddings(TEST_USER_ID, nil)
	response, err := embeddings.Embed(ctx, input)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, 1, len(response))

	// insert the embeddings and ensure that no dupliates get created
	for range 3 {
		err = embeddings.ReportUsage(ctx, txn, customer)
		assert.Nil(t, err)

		// test that the correct data exists in the database
		model := queries.New(txn)
		usage, err := model.GetTokenUsage(ctx, customer.ID)
		assert.Nil(t, err)
		if err != nil {
			return
		}
		assert.Equal(t, 1, len(usage))
	}
}
