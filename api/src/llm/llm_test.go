package llm

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"
)

func TestLLMConsistency(t *testing.T) {
	ctx := context.Background()
	pool := testingutils.GetDatabase(t, ctx)

	c := testingutils.GetTestCustomer(t, ctx, pool)
	cid := utils.GoogleUUIDToPGXUUID(c.ID)

	model := queries.New(pool)

	// ensure there is a single default at startup time
	defaults, err := model.GetStandardLLMs(ctx)
	require.NoError(t, err)
	count := 0
	for _, item := range defaults {
		if item.Llm.IsDefault {
			count++
		}
	}
	require.Equal(t, 1, count)

	// create new llm
	ll1, err := model.CreateLLM(ctx, &queries.CreateLLMParams{
		CustomerID:   cid,
		Model:        "gpt-3.5-turbo",
		Temperature:  1.0,
		Instructions: "You are a helpful and friendly assistant",
		IsDefault:    false,
	})
	require.NoError(t, err)

	// create the llm object
	obj, err := GetLLM(ctx, pool, c.ID, utils.GoogleUUIDToPGXUUID(ll1.ID))
	require.NoError(t, err)
	fmt.Println(*obj)

	// get the customer default llm by passing invalid uuid
	obj, err = GetLLM(ctx, pool, c.ID, pgtype.UUID{})
	require.NoError(t, err)
	require.Equal(t, true, obj.Llm.IsDefault)
}
