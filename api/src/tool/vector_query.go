package tool

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

type ToolVectorQuery struct{}

func newToolVectorQuery() *ToolVectorQuery {
	return &ToolVectorQuery{}
}

func (t *ToolVectorQuery) GetType() ToolType {
	return VectorQuery
}

func (t *ToolVectorQuery) GetSchema() *gollm.Tool {
	return &gollm.Tool{
		Title:       string(t.GetType()),
		Description: "Send a request against the user's private stored information. Be liberal with the use of this tool, as the tool repsonse will contain valuable information to help you create more personalized answers.",
		Schema: &ltypes.ToolSchema{
			Type: "object",
			Properties: map[string]*ltypes.ToolSchema{
				"vector_query": {
					Type:        "string",
					Description: "A simple query that can be used to query the user's private information stored in a vector database. Make sure your query contains all the semantic information you are seeking, as the vector store potentially contains lots of similar information.",
				},
			},
		},
	}
}

func (t *ToolVectorQuery) Run(
	ctx context.Context,
	l *slog.Logger,
	args *RunToolArgs,
) (*ToolResponse, error) {
	logger := l.With("tool", t.GetType())
	if err := args.Validate(); err != nil {
		return nil, slogger.Error(ctx, logger, "ARGUMENT ERROR", err)
	}

	// internal arguments
	dmodel := queries.New(args.Database)
	usageRecords := make([]*tokens.UsageRecord, 0)

	// ensure the arguments are present
	vectorQuery, exists := args.LastMessage.ToolArguments["vector_query"]
	if !exists {
		return nil, slogger.Error(ctx, logger, "the argument 'vector_query' does not exist", nil)
	}

	// simplify the vector query
	logger.InfoContext(ctx, "Simplifying vector query ...", "query", vectorQuery)

	// get the simple query llm from the database
	logger.InfoContext(ctx, "Getting the simplify llm ...")
	tmp, err := dmodel.GetInteralLLM(ctx, "Vector Query Generator")
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the simple query LLM", err)
	}
	simpleQueryLLM := llm.FromObjects(&tmp.Llm, &tmp.AvailableModel)

	// run a single completion
	simpleQueryResponse, err := simpleQueryLLM.SingleCompletion(
		ctx, logger, args.Customer.ID, prompts.RAG_SIMPLE_QUERY_SYSTEM_PROMPT,
		vectorQuery.(string),
	)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the simple query", err)
	}
	usageRecords = append(usageRecords, simpleQueryResponse.UsageRecord)

	// parse the response
	simpleQueries := strings.Split(simpleQueryResponse.Message.Message, ",")

	// get the embeddings
	embs := llm.GetEmbeddings(logger, args.Customer)
	vectorResponses := make([]*vectorstore.QueryResponse, 0)

	for _, item := range simpleQueries {
		logger.InfoContext(ctx, "Running vector query ...", "query", item)

		// run all the simple queries against the vector store
		vectorResponse, err := vectorstore.Query(ctx, logger, args.Database, &vectorstore.QueryInput{
			CustomerID: args.Customer.ID,
			Embeddings: embs,
			Query:      item,
			K:          4,
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to query the vectorstore", err)
		}
		vectorResponses = append(vectorResponses, vectorResponse)
	}

	// create separate lists
	vectors := make([]*queries.VectorStore, 0)
	docs := make([]*queries.Document, 0)
	pages := make([]*queries.WebsitePage, 0)

	for _, item := range vectorResponses {
		vectors = append(vectors, item.Vectors...)
		docs = append(docs, item.Documents...)
		pages = append(pages, item.WebsitePages...)
	}

	// remove the duplicates
	vectors = utils.RemoveDuplicates(vectors, func(val *queries.VectorStore) any {
		return val.ID
	})
	docs = utils.RemoveDuplicates(docs, func(val *queries.Document) any {
		return val.ID
	})
	pages = utils.RemoveDuplicates(pages, func(val *queries.WebsitePage) any {
		return val.ID
	})

	// craft a response for the valler
	var toolResponse string

	if len(vectors) != 0 {
		// summarize the vectors
		buf := new(bytes.Buffer)
		for _, item := range vectors {
			if _, err := buf.WriteString(item.Raw); err != nil {
				logger.Warn("failed to write to the buffer", "error", err)
			}
		}

		// send the summary
		bufStr := buf.String()
		response, err := args.ToolLLM.Summarize(ctx, logger, args.Customer.ID, bufStr)
		if err != nil {
			logger.Error("failed to summarize the content", err)
			// still pass on information to response even with an error
			toolResponse = fmt.Sprintf("[Query Response]: %s", bufStr)
		} else {
			// write the summary
			toolResponse = fmt.Sprintf("[Query Response]: %s", response.Summary)
			if response.UsageRecords != nil {
				usageRecords = append(usageRecords, response.UsageRecords...)
			}
		}
	} else {
		toolResponse = "[Query Response]: No valid information found"
	}

	message := gollm.NewToolResultMessage(args.LastMessage.ToolUseID, args.LastMessage.ToolName, toolResponse)
	arguments := make(map[string]any)
	arguments["docs"] = docs
	arguments["pages"] = pages
	message.ToolArguments = arguments

	// add the usage records
	usageRecords = append(usageRecords, embs.GetUsageRecords()...)

	// compose the response
	return &ToolResponse{
		Message:      message,
		UsageRecords: usageRecords,
	}, nil
}
