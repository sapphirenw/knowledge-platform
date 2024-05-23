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
}

func CreateLLM(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
	customer *queries.Customer,
	title string,
	modelName string,
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
	if modelName == "" {
		return nil, fmt.Errorf("modelName cannot be empty")
	}
	if temperature < 0 {
		return nil, fmt.Errorf("tempurature cannot be negative")
	}
	if instructions == "" {
		return nil, fmt.Errorf("instructions cannot be empty")
	}

	dmodel := queries.New(db)
	model, err := dmodel.CreateLLM(ctx, &queries.CreateLLMParams{
		CustomerID:   utils.GoogleUUIDToPGXUUID(customer.ID),
		Title:        title,
		Model:        modelName,
		Temperature:  temperature,
		Instructions: instructions,
		IsDefault:    isDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %s", err)
	}

	return &LLM{Llm: model}, nil
}

// Fetches an llm with the passed id. If the id is not valid, then the customer's default is used
func GetLLM(ctx context.Context, db queries.DBTX, customerId uuid.UUID, id pgtype.UUID) (*LLM, error) {
	model := queries.New(db)

	var llm *queries.Llm
	var err error

	// check if valid pgxid
	gid := utils.PGXUUIDToGoogleUUID(id)
	if gid != nil {
		// get the llm with the passed value
		llm, err = model.GetLLM(ctx, *gid)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch the llm with id: %s", err)
		}
	} else {
		// get the customer's default
		tmp, err := model.GetDefaultLLM(ctx, utils.GoogleUUIDToPGXUUID(customerId))
		if err != nil {
			return nil, fmt.Errorf("error fetching the default llm: %s", err)
		}
		llm = &queries.Llm{
			ID:           tmp.ID,
			CustomerID:   tmp.CustomerID,
			Title:        tmp.Title,
			Model:        tmp.Model,
			Temperature:  tmp.Temperature,
			Instructions: tmp.Instructions,
			IsDefault:    tmp.IsDefault,
			CreatedAt:    tmp.CreatedAt,
			UpdatedAt:    tmp.UpdatedAt,
		}
	}

	return &LLM{Llm: llm}, nil
}

func (model *LLM) GenerateSystemPrompt(prompt string) string {
	return fmt.Sprintf(prompts.PROMPT_LLM_SYSTEM, prompt, model.Instructions)
}

type CompletionArgs struct {
	Input      string
	Json       bool
	JsonSchema string
}

func (model *LLM) Completion(ctx context.Context, logger *slog.Logger, lm *gollm.LanguageModel, args *CompletionArgs) (string, error) {
	if args == nil || args.Input == "" {
		return "", fmt.Errorf("the input cannot be empty")
	}
	if args.Json && args.JsonSchema == "" {
		return "", fmt.Errorf("cannot have an empty schema with json mode enabled")
	}

	l := logger.With("args", *args)
	l.InfoContext(ctx, "Sending the completion request ...")

	input := &gollm.CompletionInput{
		Model:       model.Model,
		Temperature: model.Temperature,
		Json:        args.Json,
		JsonSchema:  args.JsonSchema,
		Input:       args.Input,
	}

	response, err := lm.DynamicCompletion(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed the dynamic completion: %s", err)
	}

	l.InfoContext(ctx, "Successfully sent the request")

	return response, nil
}
