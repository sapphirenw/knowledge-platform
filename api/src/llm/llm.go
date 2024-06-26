package llm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type LLM struct {
	*queries.Llm
	*queries.AvailableModel
}

type CompletionArgs struct {
	CustomerID   string
	Messages     []*gollm.Message
	Tools        []*gollm.Tool
	RequiredTool *gollm.Tool

	Json       bool
	JsonSchema string
}

func CreateLLM(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	customer *queries.Customer,
	title string,
	availableModelId string,
	temperature float64,
	instructions string,
	isDefault bool,
) (*LLM, error) {
	if customer == nil {
		return nil, fmt.Errorf("customer cannot be nil")
	}
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if availableModelId == "" {
		return nil, fmt.Errorf("availableModelId cannot be empty")
	}
	if temperature < 0 {
		return nil, fmt.Errorf("tempurature cannot be negative")
	}
	if instructions == "" {
		return nil, fmt.Errorf("instructions cannot be empty")
	}

	dmodel := queries.New(db)

	// get the available model
	amodel, err := dmodel.GetAvailableModel(ctx, availableModelId)
	if err != nil {
		return nil, fmt.Errorf("there was not model found: %s", err)
	}
	if amodel.IsDepreciated {
		return nil, fmt.Errorf("this model has been depreciated: %s", err)
	}

	// create the llm object
	model, err := dmodel.CreateLLM(ctx, &queries.CreateLLMParams{
		CustomerID:   utils.GoogleUUIDToPGXUUID(customer.ID),
		Title:        title,
		Model:        amodel.ID,
		Temperature:  temperature,
		Instructions: instructions,
		IsDefault:    isDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %s", err)
	}

	// get the generated object
	obj, err := dmodel.GetLLM(ctx, model.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the newly created object: %s", err)
	}

	return &LLM{Llm: &obj.Llm, AvailableModel: &obj.AvailableModel}, nil
}

// Fetches an llm with the passed id. If the id is not valid, then the customer's default is used
func GetLLM(ctx context.Context, db queries.DBTX, customerId uuid.UUID, id pgtype.UUID) (*LLM, error) {
	model := queries.New(db)

	// check if valid pgxid
	gid := utils.PGXUUIDToGoogleUUID(id)
	if gid != nil {
		// get the llm with the passed value
		llm, err := model.GetLLM(ctx, *gid)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch the llm with id: %s", err)
		}
		return &LLM{Llm: &llm.Llm, AvailableModel: &llm.AvailableModel}, nil

	} else {
		// get the customer's default
		llm, err := model.GetDefaultLLM(ctx, utils.GoogleUUIDToPGXUUID(customerId))
		if err != nil {
			return nil, fmt.Errorf("error fetching the default llm: %s", err)
		}
		return &LLM{Llm: &llm.Llm, AvailableModel: &llm.AvailableModel}, nil
	}
}

func FromObjects(llm *queries.Llm, availableModel *queries.AvailableModel) *LLM {
	return &LLM{Llm: llm, AvailableModel: availableModel}
}

func (model *LLM) GetEstimatedTokens(input string) (int32, error) {
	tokens, err := gollm.TokenEstimate(model.Llm.Model, input)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate token usage: %s", err)
	}
	return int32(tokens), nil
}

func (model *LLM) GenerateSystemPrompt(prompt string) string {
	if model.Llm.Instructions == "" && prompt == "" {
		return "Follow the internal instructions you have been given."
	}
	if model.Llm.Instructions == "" {
		return prompt
	}
	if prompt == "" {
		return model.Llm.Instructions
	}
	return fmt.Sprintf(prompts.LLM_SYSTEM, model.Llm.Instructions, prompt)
}

func (model *LLM) Completion(
	ctx context.Context,
	logger *slog.Logger,
	args *CompletionArgs,
) (*gollm.CompletionResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("the input cannot be empty")
	}
	if args.CustomerID == "" {
		return nil, fmt.Errorf("the CustomerID cannot be empty")
	}
	if args.Json && args.JsonSchema == "" {
		return nil, fmt.Errorf("cannot have an empty schema with json mode enabled")
	}
	if args.Messages == nil || len(args.Messages) < 2 {
		return nil, fmt.Errorf("the messages array must be filled")
	}

	l := logger.With("completionType", "multi")
	l.InfoContext(ctx, "Sending the completion request ...")

	// create a copy of the list
	msgs := make([]*gollm.Message, len(args.Messages))
	copy(msgs, args.Messages)

	// add model specific instructions to the system message
	msgs[0].Message = fmt.Sprintf("General Instructions: %s\n\nSpecific Instructions: %s", model.Llm.Instructions, msgs[0].Message)

	input := &gollm.CompletionInput{
		Model:        model.Llm.Model,
		Temperature:  model.Llm.Temperature,
		Json:         args.Json,
		JsonSchema:   args.JsonSchema,
		Conversation: msgs,
		Tools:        args.Tools,
		RequiredTool: args.RequiredTool,
	}

	lm := gollm.NewLanguageModel(args.CustomerID, logger, nil)
	response, err := lm.Completion(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed the dynamic completion: %s", err)
	}

	l.InfoContext(ctx, "Successfully sent the request")

	return response, nil
}

// Convenience function wrapper around the `Completion` function for performing one-off requests
func (model *LLM) SingleCompletion(
	ctx context.Context,
	logger *slog.Logger,
	customerId uuid.UUID,
	systemMessage string,
	input string,
) (*gollm.CompletionResponse, error) {
	messages := make([]*gollm.Message, 0)
	messages = append(messages, gollm.NewSystemMessage(systemMessage))
	messages = append(messages, gollm.NewUserMessage(input))

	return model.Completion(ctx, logger, &CompletionArgs{
		CustomerID: customerId.String(),
		Messages:   messages,
	})
}
