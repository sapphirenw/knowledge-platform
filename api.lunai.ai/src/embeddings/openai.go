package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jake-landersweb/gollm/src/ltypes"
	"github.com/jake-landersweb/gollm/src/tokens"
	"github.com/pgvector/pgvector-go"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type OpenAIEmbeddings struct {
	userId int64
	model  string
	logger *slog.Logger

	// track token usage
	tokenRecords []*tokens.TokenRecord
}

type OpenAIEmbeddingsOpts struct {
	Model  string
	Logger *slog.Logger
}

func NewOpenAIEmbeddings(userId int64, opts *OpenAIEmbeddingsOpts) *OpenAIEmbeddings {
	if opts == nil {
		opts = &OpenAIEmbeddingsOpts{}
	}
	if opts.Logger == nil {
		opts.Logger = utils.DefaultLogger()
	}
	if opts.Model == "" {
		opts.Model = OPENAI_EMBEDDINGS_MODEL
	}

	opts.Logger = opts.Logger.With("userId", userId, "model", opts.Model)

	return &OpenAIEmbeddings{
		userId: userId,
		model:  opts.Model,
		logger: opts.Logger,
	}
}

func (e *OpenAIEmbeddings) UserIdString() string {
	return fmt.Sprintf("%d", e.userId)
}

func (e *OpenAIEmbeddings) Embed(ctx context.Context, input string) ([]*EmbeddingsData, error) {
	// chunk the input
	chunks := utils.ChunkStringEqualUntilN(input, OPENAI_EMBEDDINGS_INPUT_MAX)
	response, err := e.openAIEmbed(ctx, e.logger, e.model, chunks)
	if err != nil {
		return nil, err
	}

	// track token usage
	e.tokenRecords = append(e.tokenRecords, tokens.NewTokenRecordFromGPTUsage(e.model, &response.Usage))

	// convert openai response into pgvector data types
	list := make([]*EmbeddingsData, 0)
	for idx := range chunks {
		list = append(list, &EmbeddingsData{
			Raw:       chunks[idx],
			Embedding: pgvector.NewVector(utils.ConvertSlice(response.Data[idx].Embedding, func(i float64) float32 { return float32(i) })),
		})
	}

	return list, nil
}

func (e *OpenAIEmbeddings) ReportUsage(ctx context.Context, db queries.DBTX) error {
	e.logger.InfoContext(ctx, "Reporting usage", "length", len(e.tokenRecords))
	model := queries.New(db)

	// insert all internal token records
	for idx, item := range e.tokenRecords {
		e.logger.InfoContext(ctx, "Posting to database ...", "index", idx)
		_, err := model.CreateTokenUsage(ctx, &queries.CreateTokenUsageParams{
			ID:           utils.GoogleUUIDToPGXUUID(item.ID),
			CustomerID:   e.userId,
			Model:        e.model,
			InputTokens:  int32(item.InputTokens),
			OutputTokens: int32(item.OutputTokens),
			TotalTokens:  int32(item.TotalTokens),
		})
		if err != nil {
			return err
		}
		e.logger.InfoContext(ctx, "Done.", "index", idx)
	}

	e.logger.InfoContext(ctx, "Successfully reported usage")

	return nil
}

func (e *OpenAIEmbeddings) openAIEmbed(ctx context.Context, logger *slog.Logger, model string, input []string) (*OpenAIEmbeddingResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "null" {
		return nil, fmt.Errorf("the env variable `OPENAI_API_KEY` is required to be set")
	}

	// create the body
	comprequest := OpenAIEmbeddingRequest{
		Input:      input,
		Model:      model,
		Dimensions: OPENAI_EMBEDDINGS_LENGTH,
		User:       e.UserIdString(),
	}

	enc, err := json.Marshal(&comprequest)
	if err != nil {
		return nil, fmt.Errorf("there was an issue encoding the body into json: %v", err)
	}

	// create the request
	req, err := http.NewRequest("POST", OPENAI_EMBEDDINGS_BASE_URL, bytes.NewBuffer(enc))
	if err != nil {
		return nil, fmt.Errorf("there was an issue creating the http request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// send the request
	client := &http.Client{}

	retries := 3
	backoff := 1 * time.Second

	for attempt := 0; attempt < retries; attempt++ {
		logger.InfoContext(ctx, "Sending embeddings request...", "chunks", len(input))
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("there was an unknown issue with the request: %v", err)
		}

		// read the body
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("there was an issue reading the body: %v", err)
		}

		logger.InfoContext(ctx, "Completed request", "statusCode", resp.StatusCode)

		// parse into the completion response object
		var response OpenAIEmbeddingResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("there was an issue unmarshalling the request body: %v", err)
		}

		// act based on the error
		if response.Error == nil {
			return &response, nil

		} else {
			// act based on the error
			switch response.Error.Type {
			case ltypes.GPT_ERROR_INVALID:
				return nil, fmt.Errorf("there was a validation error: %s", string(body))
			case ltypes.GPT_ERROR_RATE_LIMIT:
				// rate limit, so wait some extra time and continue
				logger.WarnContext(ctx, "Rate limit error hit. Waiting for an additional 2 seconds...")
				time.Sleep(time.Second * 2)
			// case ltypes.GPT_ERROR_TOKENS_LIMIT:
			case ltypes.GPT_ERROR_AUTH:
				return nil, fmt.Errorf("the user is not authenticated: %s", string(body))
			case ltypes.GPT_ERROR_NOT_FOUND:
				return nil, fmt.Errorf("the requested resource was not found: %s", string(body))
			case ltypes.GPT_ERROR_SERVER:
				// internal server error, wait and try again
				logger.WarnContext(ctx, "There was an issue on OpenAI's side. Waiting 2 seconds and trying again ...", "body", string(body))
				time.Sleep(time.Second * 2)
			case ltypes.GPT_ERROR_PERMISSION:
				return nil, fmt.Errorf("the requested resource was not found: %s", string(body))
			default:
				return nil, fmt.Errorf("there was an unknown error: %s", string(body))
			}
		}

		if attempt < retries-1 {
			sleep := backoff + time.Duration(rand.Intn(1000))*time.Millisecond // Add jitter
			time.Sleep(sleep)
			backoff *= 2 // Double the backoff interval
		} else if resp != nil && resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("there was an issue with the request and could not recover: %s", string(body))
		}
	}

	return nil, err
}
