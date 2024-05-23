package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Project struct {
	*queries.Project

	logger              *slog.Logger
	ideaGererationModel *llm.LLM
}

func CreateProject(
	ctx context.Context,
	db queries.DBTX,
	logger *slog.Logger,
	customer *queries.Customer,
	title string,
	topic string,
	ideaGenerationModel *llm.LLM,
) (*Project, error) {
	if customer == nil {
		return nil, fmt.Errorf("the customer cannot be nil")
	}
	if title == "" {
		return nil, fmt.Errorf("the title cannot be empty")
	}
	if topic == "" {
		return nil, fmt.Errorf("the topic cannot be empty")
	}

	l := logger.With("customerId", customer.ID.String())
	l.InfoContext(ctx, "Creating new project", "customerId", customer.ID.String())

	modelId := pgtype.UUID{}
	if ideaGenerationModel != nil {
		modelId = utils.GoogleUUIDToPGXUUID(ideaGenerationModel.ID)
	}

	dmodel := queries.New(db)
	project, err := dmodel.CreateProject(ctx, &queries.CreateProjectParams{
		CustomerID:            customer.ID,
		Title:                 title,
		Topic:                 topic,
		IdeaGenerationModelID: modelId,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating the project: %s", err)
	}

	l.InfoContext(ctx, "Successfully created new project", "project", *project)

	return &Project{
		Project:             project,
		logger:              l.With("projectId", project.ID.String(), "model.ID", ideaGenerationModel.ID.String(), "model.Title", ideaGenerationModel.Title),
		ideaGererationModel: ideaGenerationModel,
	}, nil
}

func GetProject(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	projectId uuid.UUID,
) (*Project, error) {
	l := logger.With("projectId", projectId.String())
	l.DebugContext(ctx, "Fetching project")

	dmodel := queries.New(db)
	project, err := dmodel.GetProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("error getting project: %s", err)
	}

	l.DebugContext(ctx, "Successfully got project")
	return &Project{
		Project: project,
		logger:  l.With("projectId", project.ID.String(), "customerId", project.CustomerID.String()),
	}, nil

}

func (p *Project) GetGenerationModel(ctx context.Context, db queries.DBTX) (*llm.LLM, error) {
	if p.ideaGererationModel != nil {
		return p.ideaGererationModel, nil
	}
	model, err := llm.GetLLM(ctx, db, p.CustomerID, p.IdeaGenerationModelID)
	if err != nil {
		return nil, fmt.Errorf("error getting llm: %s", err)
	}
	p.ideaGererationModel = model
	p.logger = p.logger.With("model.ID", p.ideaGererationModel.ID.String(), "model.Title", p.ideaGererationModel.Title)
	return p.ideaGererationModel, nil
}

type GenerateIdeasArgs struct {
	ConversationId string
	Feedback       string
	N              int
}

type GenerateIdeasResponse struct {
	Ideas          []*ProjectIdea
	ConversationId uuid.UUID
}

func (p *Project) GenerateIdeas(
	ctx context.Context,
	db queries.DBTX,
	args *GenerateIdeasArgs,
) (*GenerateIdeasResponse, error) {
	// parse arguments
	if args == nil {
		args = &GenerateIdeasArgs{}
	}
	if args.Feedback == "" {
		args.Feedback = "None."
	}
	if args.N == 0 {
		args.N = 1
	}
	// check if uuid can be parsed
	if args.ConversationId != "" {
		if _, err := uuid.Parse(args.ConversationId); err != nil {
			return nil, fmt.Errorf("failed to parse conversationId: %s", err)
		}
	}

	logger := p.logger.With("args", *args)

	// get the generation model
	model, err := p.GetGenerationModel(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get the model: %s", err)
	}

	// determine where to get the conversation from
	var conv *llm.Conversation
	var prompt string
	if args.ConversationId == "" {
		logger.InfoContext(ctx, "Generating new ideas ...")

		// create a new conversation
		conv, err = llm.CreateConversation(ctx, logger, db, p.CustomerID, prompts.PROMPT_PROJECT_IDEA_SYSTEM, "Idea Generation", "idea-generation")
		if err != nil {
			return nil, fmt.Errorf("failed to create the conversation: %s", err)
		}

		// define a new prompt
		prompt = fmt.Sprintf("Title: %s\nTopic: %s\nNumber: %d", p.Title, p.Topic, args.N)
	} else {
		logger.InfoContext(ctx, "Generating ideas from an existing conversation ...")

		// get the existing conversation
		conv, err = llm.GetConversation(ctx, logger, db, uuid.MustParse(args.ConversationId))
		if err != nil {
			return nil, fmt.Errorf("failed to get the conversation: %s", err)
		}

		// create a prompt with the feedback
		prompt = fmt.Sprintf("This was not quite what I am looking for. Please try again with this feedback: %s.\nRemember, respond in the JSON format given previously.", args.Feedback)
	}

	logger.InfoContext(ctx, "Running the completion ...")

	// run the completion against the conversation
	response, err := conv.Completion(ctx, db, model, &llm.CompletionArgs{
		Input:      prompt,
		Json:       true,
		JsonSchema: `{"ideas": [{"title", string}]}`,
	})
	if err != nil {
		return nil, fmt.Errorf("the completion failed")
	}
	if response == "" {
		return nil, fmt.Errorf("there was an unknown with the model completion, the response was empty")
	}

	logger.InfoContext(ctx, "Successfully parsed the ideas")

	// serialize from json
	var ideas projectIdeas
	if err := json.Unmarshal([]byte(response), &ideas); err != nil {
		return nil, fmt.Errorf("failed to de-serialize json: %s %s", err, response)
	}

	return &GenerateIdeasResponse{
		Ideas:          ideas.Ideas,
		ConversationId: conv.ID,
	}, nil
}
