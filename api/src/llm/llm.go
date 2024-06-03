package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
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
			Public:       tmp.Public,
			CreatedAt:    tmp.CreatedAt,
			UpdatedAt:    tmp.UpdatedAt,
		}
	}

	return &LLM{Llm: llm}, nil
}

func (model *LLM) GenerateSystemPrompt(prompt string) string {
	if model.Instructions == "" && prompt == "" {
		return "Follow the internal instructions you have been given."
	}
	if model.Instructions == "" {
		return prompt
	}
	if prompt == "" {
		return model.Instructions
	}
	return fmt.Sprintf(prompts.LLM_SYSTEM, model.Instructions, prompt)
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

	l := logger.With("completionType", "multi", "args", *args)
	l.InfoContext(ctx, "Sending the completion request ...")

	// check whether to add llm specific instructions
	msg := args.Input
	numMessages := len(lm.GetConversation())
	if numMessages == 0 {
		logger.DebugContext(ctx, "Adding general llm instructions")
		msg = fmt.Sprintf("General Instructions: %s\n\nSpecific Instructions: %s", model.Instructions, msg)
	} else {
		logger.DebugContext(ctx, "Continuing conversation", "numMessages", numMessages)
	}

	input := &gollm.CompletionInput{
		Model:       model.Model,
		Temperature: model.Temperature,
		Json:        args.Json,
		JsonSchema:  args.JsonSchema,
		Input:       msg,
	}

	response, err := lm.DynamicCompletion(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed the dynamic completion: %s", err)
	}

	l.InfoContext(ctx, "Successfully sent the request")

	return response, nil
}

type SingleCompletionResponse[T any] struct {
	Result       T
	UsageRecords []*tokens.TokenRecord
}

// for performing a single shot completion against an llm, and reporting the usage to the database.
// This is not to be used for conversations, only single operations against the model.
func SingleCompletion(
	ctx context.Context,
	model *LLM,
	logger *slog.Logger,
	customerId uuid.UUID,
	sysMessage string,
	usageRecords chan *tokens.TokenRecord,
	args *CompletionArgs,
) (string, error) {
	if args == nil || args.Input == "" {
		return "", fmt.Errorf("the input cannot be empty")
	}
	if args.Json && args.JsonSchema == "" {
		return "", fmt.Errorf("cannot have an empty schema with json mode enabled")
	}

	l := logger.With("completionType", "single", "args", *args)
	l.DebugContext(ctx, "Creating the gollm ...")
	lm := gollm.NewLanguageModel(customerId.String(), l, sysMessage, nil)
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
	l.DebugContext(ctx, "Gathered response")

	// parse usage records
	for _, item := range lm.GetTokenRecords() {
		usageRecords <- item
	}

	l.InfoContext(ctx, "Successfully sent the one-shot request")
	return response, nil
}

func SingleCompletionJson[T any](
	ctx context.Context,
	model *LLM,
	logger *slog.Logger,
	customerId uuid.UUID,
	sysMessage string,
	usageRecords chan *tokens.TokenRecord,
	args *CompletionArgs,
) (*T, error) {
	var result T

	response, err := SingleCompletion(ctx, model, logger, customerId, sysMessage, usageRecords, args)
	if err != nil {
		return nil, err
	}

	// parse the json
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %s", err)
	}

	return &result, nil
}
