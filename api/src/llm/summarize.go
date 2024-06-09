package llm

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
)

// Summarizes the input text using the provided llm. This function will chunk and run
// the summaries in go-routines to improve performance.
func (llm *LLM) Summarize(
	ctx context.Context,
	logger *slog.Logger,
	customerId uuid.UUID,
	tokenRecords chan *tokens.UsageRecord,
	input string,
) (string, error) {
	return "", nil
	// // chunk the string if needed
	// estimatedTokens, err := llm.GetEstimatedTokens(input)
	// if err != nil {
	// 	return "", err
	// }

	// chunks := make([]string, 0)
	// if estimatedTokens > llm.InputTokenLimit {
	// 	logger.InfoContext(ctx, "Chunking input ...", "tokens", estimatedTokens)
	// 	chunks = gollm.ChunkStringEqualUntilN(input, int(llm.InputTokenLimit))
	// } else {
	// 	chunks = append(chunks, input)
	// }
	// logger.InfoContext(ctx, "Processing chunks ...", "length", len(chunks))

	// // create an internal go-routine sync state for the chunks
	// var wg sync.WaitGroup
	// errCh := make(chan error, len(chunks))
	// responses := make(chan string, len(chunks))

	// // loop through and summarize all content as neeed and chunk at the end
	// for i, item := range chunks {
	// 	wg.Add(1)
	// 	go func(index int, chunk string) {
	// 		defer wg.Done()
	// 		l := logger.With("index", index)
	// 		l.InfoContext(ctx, "Processing chunk ...")

	// 		response, err := SingleCompletion(
	// 			ctx, llm, logger, customerId,
	// 			prompts.SUMMARY_SYSTEM_PROMPT,
	// 			tokenRecords,
	// 			&CompletionArgs{Input: chunk},
	// 		)
	// 		if err != nil {
	// 			errCh <- fmt.Errorf("failed to summarize content: %s", err)
	// 			return
	// 		}
	// 		responses <- response
	// 		l.InfoContext(ctx, "Successfully processed chunk")
	// 	}(i, item)
	// }

	// // collect the routines
	// wg.Wait()
	// close(errCh)
	// close(responses)

	// // parse for errors
	// for err := range errCh {
	// 	if err != nil {
	// 		return "", err
	// 	}
	// }

	// // compose the summary buffer
	// summary := new(bytes.Buffer)
	// for item := range responses {
	// 	summary.WriteString(item + "\n")
	// }

	// // return the entire buffer
	// return summary.String(), err
}
