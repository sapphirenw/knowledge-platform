package project

import (
	"context"
	"fmt"
	"testing"

	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"
)

func TestLinkedInPostCreate(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.GetDefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)
	project := getTestProject(t, ctx, logger, pool, c)

	// create the linkedin post
	logger.InfoContext(ctx, "Creating the linkedin post ...")
	post, err := project.NewLinkedInPost(ctx, pool, "", "My New LinkedIn Post")
	require.NoError(t, err)

	fmt.Printf("Post: %s\n", post.LinkedinPost.Title)
}

func TestLinkedInPostConfig(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.GetDefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)
	project := getTestProject(t, ctx, logger, pool, c)

	// check for a default
	dmodel := queries.New(pool)
	projectDefaultConfig, err := dmodel.GetLinkedInPostConfig(ctx, &queries.GetLinkedInPostConfigParams{
		ProjectID: utils.GoogleUUIDToPGXUUID(project.ID),
	})
	require.NoError(t, err)
	require.NotNil(t, projectDefaultConfig)

	// create a linkedin post
	post, err := project.NewLinkedInPost(ctx, pool, "", "My New LinkedIn Post")
	require.NoError(t, err)

	// check the linkedin post default
	postConfig, err := post.GetConfig(ctx, logger, pool)
	require.NoError(t, err)

	// ensure the configs are the same
	require.Equal(t, projectDefaultConfig.ID.String(), postConfig.ID.String())

	// create a new config for this specific post
	newConfig, err := project.CreateLinkedInPostConfig(ctx, pool, &createLinkedinPostConfigRequest{
		LinkedInPostId:            post.ID.String(),
		MinSections:               1,
		MaxSections:               3,
		NumDocuments:              1,
		NumWebsitePages:           1,
		LlmContentCenerationId:    utils.StringFromPGXUUID(projectDefaultConfig.LlmContentGenerationID),
		LlmVectorSummarizationId:  utils.StringFromPGXUUID(projectDefaultConfig.LlmVectorSummarizationID),
		LlmWebsiteSummarizationId: utils.StringFromPGXUUID(projectDefaultConfig.LlmWebsiteSummarizationID),
		LlmProofReadingId:         utils.StringFromPGXUUID(projectDefaultConfig.LlmProofReadingID),
	})
	require.NoError(t, err)

	// get the post config again to ensure it has been updated
	postConfigNew, err := post.GetConfig(ctx, logger, pool)
	require.NoError(t, err)

	// ensure they are equal
	require.Equal(t, newConfig.ID.String(), postConfigNew.ID.String())
}

func TestLinkedInPostGenerate(t *testing.T) {
	ctx := context.Background()
	logger := testingutils.GetDefaultLogger()
	pool := testingutils.GetDatabase(t, ctx)
	c := testingutils.GetTestCustomer(t, ctx, pool)
	project := getTestProject(t, ctx, logger, pool, c)

	// create the post
	post, err := project.NewLinkedInPost(ctx, pool, "", "Congratulations My Coworker!")
	require.NoError(t, err)

	// send the first generation
	response1, err := project.GenerateLinkedInPost(ctx, pool, post, &generateLinkedInPostRequest{
		ConversationId: "",
		Input:          "A post congratulating my co-worker for their promotion.",
	})
	require.NoError(t, err)

	// send corrections
	response2, err := project.GenerateLinkedInPost(ctx, pool, post, &generateLinkedInPostRequest{
		ConversationId: response1.ConversationId.String(),
		Input:          "Make it more aggressive! We need some ZEST in the workplace.",
	})
	require.NoError(t, err)

	// print the conversation
	conv, err := llm.GetConversation(ctx, logger, pool, response2.ConversationId)
	require.NoError(t, err)
	conv.PrintConversation()
}
