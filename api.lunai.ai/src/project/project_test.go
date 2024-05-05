package project

import (
	"context"
	"fmt"
	"testing"

	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/stretchr/testify/require"
)

func TestCreateProjectIdea(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.DefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)

	// create a custom model to use for generating ideas on the project
	model, err := llm.CreateLLM(
		ctx, logger, pool, c, "Idea Generator", "gpt-3.5-turbo", 1.2,
		"You are inquisitive, yet bold with your outputs. You have a keen eye for what people will want to listen too, without overloading with buzzwords or unnecessarily complex language. You are also pay extreme attention to the instructions you are given.",
		false,
	)
	require.NoError(t, err)

	// create a test project
	project, err := CreateProject(
		ctx,
		pool,
		logger,
		c,
		"Crosscheck Sports",
		"Creating engaging and informative sports-related content.",
		model,
	)
	require.NoError(t, err)

	m, err := project.GetGenerationModel(ctx, pool)
	require.NoError(t, err)
	require.Equal(t, "Idea Generator", m.Title)

	ideas, err := project.GenerateIdeas(ctx, pool, 2)
	require.NoError(t, err)

	// ensure that the model usage was reported properly
	dmodel := queries.New(pool)
	records, err := dmodel.GetTokenUsage(ctx, c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, records)

	fmt.Println("\nIdeas:")
	for _, item := range ideas {
		fmt.Println("- " + item.Title)
	}

	fmt.Println("\nUsage:")
	for _, item := range records {
		fmt.Printf("- %s: %d\n", item.Model, item.TotalTokens)
	}
}
