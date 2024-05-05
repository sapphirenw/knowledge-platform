package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
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

func (p *Project) GenerateIdeas(ctx context.Context, db queries.DBTX, n int) ([]*ProjectIdea, error) {
	logger := p.logger.With()
	logger.InfoContext(ctx, "Generating project ideas")

	// get the generation model
	model, err := p.GetGenerationModel(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to get the model: %s", err)
	}

	// generate a system prompt from the model
	sys_prompt := "You are a model that has been specifically designed to generate ideas for online content based on a topic. You will receive a title, a topic, and a number of ideas to generate as input from the user. You will proceed to generate engaging ideas in the JSON format specified."
	sys_prompt = model.GenerateSystemPrompt(sys_prompt)

	// create an llm object
	lm := gollm.NewLanguageModel(p.CustomerID.String(), p.logger, &gollm.NewLanguageModelArgs{
		SystemMessage: sys_prompt,
	})

	// create an input for idea generation
	prompt := fmt.Sprintf("Title: %s\nTopic: %sNumber: %d", p.Title, p.Topic, n)
	response, err := model.Completion(ctx, logger, lm, &llm.CompletionArgs{
		Input:      prompt,
		Json:       true,
		JsonSchema: `{"ideas": [{"title", string}]}`,
	})
	if err != nil {
		return nil, fmt.Errorf("failed model completion: %s", err)
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

	// report usage
	if err := utils.ReportUsage(ctx, logger, db, p.CustomerID, lm.GetTokenRecords()); err != nil {
		return nil, fmt.Errorf("failed to report model usage: %s", err)
	}

	return ideas.Ideas, nil
}
