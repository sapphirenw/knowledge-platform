package project

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/stretchr/testify/require"
)

func TestCreateProjectIdea(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.GetDefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)

	// create a test project
	project := getTestProject(t, ctx, logger, pool, c)

	m, err := project.GetGenerationModel(ctx, pool)
	require.NoError(t, err)
	require.Equal(t, "Idea Generator", m.Title)

	// test the ability to generate ideas
	response, err := project.GenerateIdeas(ctx, pool, &generateIdeasRequest{K: 2})
	require.NoError(t, err)

	// ensure that the model usage was reported properly
	dmodel := queries.New(pool)
	records, err := dmodel.GetTokenUsage(ctx, c.ID)
	require.NoError(t, err)
	require.NotEmpty(t, records)

	fmt.Println("\nIdeas:")
	for _, item := range response.Ideas {
		fmt.Println("- " + item.Title)
	}

	fmt.Println("\nUsage:")
	for _, item := range records {
		fmt.Printf("- %s: %d\n", item.Model, item.TotalTokens)
	}

	// give feedback
	response, err = project.GenerateIdeas(ctx, pool, &generateIdeasRequest{
		ConversationId: response.ConversationId.String(),
		Feedback:       "I would like this content more tailored to NHL hockey",
		K:              2,
	})
	require.NoError(t, err)

	// ensure that tokens were tracked properly
	records, err = dmodel.GetTokenUsage(ctx, c.ID)
	require.NoError(t, err)
	require.Equal(t, 2, len(records))

	fmt.Println("\nIdeas:")
	for _, item := range response.Ideas {
		fmt.Println("- " + item.Title)
	}

	fmt.Println("\nUsage:")
	for _, item := range records {
		fmt.Printf("- %s: %d\n", item.Model, item.TotalTokens)
	}

	// print out the conversation
	conv, err := llm.GetConversation(ctx, logger, pool, response.ConversationId)
	require.NoError(t, err)
	conv.PrintConversation()
	require.Equal(t, 5, len(conv.Messages))
}

func getTestProject(t *testing.T, ctx context.Context, logger *slog.Logger, pool *pgxpool.Pool, c *queries.Customer) *Project {
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

	return project
}
