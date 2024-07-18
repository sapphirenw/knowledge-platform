package llm

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/textsplitter"
)

type SummarizeResponse struct {
	Summary      string
	UsageRecords []*tokens.UsageRecord
}

// Summarizes the input text using the provided llm. This function will chunk and run
// the summaries in go-routines to improve performance.
func (llm *LLM) Summarize(
	ctx context.Context,
	logger *slog.Logger,
	customerId uuid.UUID,
	input string,
) (*SummarizeResponse, error) {
	// chunk the string if needed
	estimatedTokens, err := llm.GetEstimatedTokens(input)
	if err != nil {
		return nil, err
	}

	chunks := make([]string, 0)
	if estimatedTokens > llm.AvailableModel.InputTokenLimit {
		splitter := textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(int(llm.AvailableModel.InputTokenLimit)),
			textsplitter.WithChunkOverlap(100),
		)

		logger.InfoContext(ctx, "Chunking input ...", "tokens", estimatedTokens)
		chunks, err = splitter.SplitText(input)
		if err != nil {
			return nil, fmt.Errorf("failed to split the text")
		}
	} else {
		chunks = append(chunks, input)
	}
	logger.InfoContext(ctx, "Processing chunks ...", "length", len(chunks))

	// create an internal go-routine sync state for the chunks
	var wg sync.WaitGroup
	errCh := make(chan error, len(chunks))
	responses := make(chan string, len(chunks))
	usageRecordsChan := make(chan *tokens.UsageRecord, len(chunks))

	// loop through and summarize all content as neeed and chunk at the end
	for i, item := range chunks {
		wg.Add(1)
		go func(index int, chunk string) {
			defer wg.Done()
			l := logger.With("index", index)
			l.InfoContext(ctx, "Processing chunk ...")

			response, err := llm.SingleCompletion(ctx, logger, customerId, prompts.SUMMARY_SYSTEM_PROMPT, chunk)
			if err != nil {
				errCh <- fmt.Errorf("failed to summarize content: %s", err)
				return
			}
			responses <- response.Message.Message
			usageRecordsChan <- response.UsageRecord
			l.InfoContext(ctx, "Successfully processed chunk")
		}(i, item)
	}

	// collect the routines
	wg.Wait()
	close(errCh)
	close(responses)
	close(usageRecordsChan)

	// parse for errors
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	// compose the summary buffer
	summary := new(bytes.Buffer)
	for item := range responses {
		summary.WriteString(item + "\n")
	}

	// create a list of usage records
	usageRecords := make([]*tokens.UsageRecord, 0)
	for item := range usageRecordsChan {
		if item != nil {
			usageRecords = append(usageRecords, item)
		}
	}

	// return the entire buffer
	return &SummarizeResponse{
		Summary:      summary.String(),
		UsageRecords: usageRecords,
	}, err
}
