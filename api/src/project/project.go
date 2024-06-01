package project

import (
	"context"
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
		logger:              l.With("projectId", project.ID.String()),
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

func (p *Project) GenerateIdeas(
	ctx context.Context,
	db queries.DBTX,
	args *generateIdeasRequest,
) (*generateIdeasResponse, error) {
	// parse arguments
	if args == nil {
		args = &generateIdeasRequest{}
	}
	if args.Feedback == "" {
		args.Feedback = "None."
	}
	if args.K == 0 {
		args.K = 1
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
	conv, err := llm.AutoConversation(
		ctx,
		logger,
		db,
		p.CustomerID,
		args.ConversationId,
		prompts.PROJECT_IDEA_SYSTEM,
		fmt.Sprintf("Idea Generation for project: %s", p.Title),
		"idea-generation",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get the conversation: %s", err)
	}

	var prompt string
	if args.ConversationId == "" {
		logger.InfoContext(ctx, "Generating new ideas ...")
		prompt = fmt.Sprintf("Title: %s\nTopic: %s\nNumber: %d", p.Title, p.Topic, args.K)
	} else {
		logger.InfoContext(ctx, "Generating ideas from an existing conversation ...")
		prompt = fmt.Sprintf("This was not quite what I am looking for. Please try again with this feedback: %s.\nRemember, respond in the JSON format given previously.", args.Feedback)
	}

	logger.InfoContext(ctx, "Running the completion ...")

	// run the completion against the conversation
	response, err := llm.JsonCompletion[projectIdeas](
		conv, ctx, db, model, prompt,
		`{"ideas": [{"title", string}]}`,
	)
	if err != nil {
		return nil, fmt.Errorf("the completion failed: %s", err)
	}

	logger.InfoContext(ctx, "Successfully parsed the ideas")

	return &generateIdeasResponse{
		Ideas:          response.Ideas,
		ConversationId: conv.ID,
	}, nil
}

func (p *Project) AddIdeas(
	ctx context.Context,
	db queries.DBTX,
	args *addIdeasRequest,
) ([]*queries.ProjectIdea, error) {
	logger := p.logger.With("func", "p.AddIdeas")

	// parse the id
	var pid pgtype.UUID
	err := pid.Scan(args.ConversationId)
	if err != nil && args.ConversationId != "" {
		return nil, fmt.Errorf("failed to parse the conversationId: %s", err)
	}

	// create all records
	logger.InfoContext(ctx, "Creating project ideas ...", "length", len(args.Ideas))

	dmodel := queries.New(db)
	ideas := make([]*queries.ProjectIdea, 0)
	for _, item := range args.Ideas {
		idea, err := dmodel.CreateProjectIdea(ctx, &queries.CreateProjectIdeaParams{
			ProjectID:      p.ID,
			ConversationID: pid,
			Title:          item.Title,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create the project idea: %s", err)
		}
		ideas = append(ideas, idea)
	}

	logger.InfoContext(ctx, "Successfully created project ideas")

	return ideas, nil
}

func (p *Project) GetIdeas(ctx context.Context, db queries.DBTX) ([]*queries.ProjectIdea, error) {
	dmodel := queries.New(db)
	ideas, err := dmodel.GetProjectIdeas(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the project ideas: %s", err)
	}
	return ideas, nil
}
