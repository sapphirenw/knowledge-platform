package customer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

const (
	toolVectorQuery = "vector_query"
)

type ragRequest struct {
	// general params
	Input        string `json:"input"`
	CheckQuality bool   `json:"checkQuality"`

	// models
	SummaryModelId string `json:"summaryModelId"`
	ChatModelId    string `json:"chatModelId"`

	// ids for scoped content
	FolderIds      []string `json:"folderIds"`
	DocumentIds    []string `json:"documentIds"`
	WebsiteIds     []string `json:"websiteIds"`
	WebsitePageIds []string `json:"websitePageIds"`
}

func (r ragRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)
	return p
}

type ragResponse struct {
	ConversationId string         `json:"conversationId"`
	Message        *gollm.Message `json:"message"`
}

func handleRAG(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	c *Customer,
) {
	// parse the request
	body, valid := request.Decode[ragRequest](w, r, c.logger)
	if !valid {
		return
	}

	// start a transaction
	tx, err := pool.Begin(r.Context())
	if err != nil {
		c.logger.Error("failed to start transaction", "error", err)
		http.Error(w, "There was a database issue", http.StatusInternalServerError)
		return
	}
	defer tx.Commit(r.Context())

	// send the response
	conversationId := r.URL.Query().Get("conversationId")
	response, err := c.RAG(r.Context(), tx, conversationId, &body)
	if err != nil {
		tx.Rollback(r.Context())
		c.logger.Error("failed to query the vectorstore", "error", err)
		http.Error(w, "There was an internal issue", http.StatusInternalServerError)
		return
	}

	request.Encode(w, r, c.logger, http.StatusOK, response)
}

func (c *Customer) RAG(
	ctx context.Context,
	db queries.DBTX,
	conversationId string,
	args *ragRequest,
) (*ragResponse, error) {
	logger := c.logger.With("function", "RAG")
	logger.InfoContext(ctx, "Beginning document retrieval pathway")

	// get the conversation

	// initial setup
	logger.DebugContext(ctx, "Getting required objects ...")
	tools := getRAGTools()

	// track all token usage across this request through a buffered channel
	// var tokenMutex sync.Mutex
	usageRecords := make([]*tokens.UsageRecord, 0)
	reportUsage := func() error {
		if err := utils.ReportUsage(ctx, logger, db, c.ID, usageRecords, nil); err != nil {
			return fmt.Errorf("failed to report the usage: %s", err)
		}
		return nil
	}

	/// get the chat llm
	logger.InfoContext(ctx, "Getting the chat llm ...")
	var chatLLMId pgtype.UUID
	chatLLMId.Scan(args.ChatModelId)
	chatLLM, err := llm.GetLLM(ctx, db, c.ID, chatLLMId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the chat llm: %s", err)
	}

	// // get the conversation
	logger.InfoContext(ctx, "Getting conversation ...")
	conv, err := conversation.AutoConversation(
		ctx,
		logger,
		db,
		c.ID,
		chatLLM,
		conversationId,
		prompts.RAG_COMPLETE_SYSTEM_PROMPT,
		"Information Chat",
		"rag",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the conversation: %s", err)
	}

	// check the state of the conversation
	if conv.New {
		// send a request for a tool usage
		message := gollm.NewUserMessage(args.Input)
		response, err := conv.Completion(ctx, db, chatLLM, message, tools, tools[0])
		if err != nil {
			return nil, fmt.Errorf("failed the completion: %s", err)
		}

		// report the usage
		usageRecords = append(usageRecords, response.UsageRecord)
		if err := reportUsage(); err != nil {
			return nil, err
		}

		// return the message
		return &ragResponse{
			ConversationId: conv.ID.String(),
			Message:        response.Message,
		}, nil
	}

	// check the state of the last message
	messages := conv.GetMessages()
	lastMessage := conv.GetMessages()[len(messages)-1]
	logger.InfoContext(ctx, "Parsing role...", "role", lastMessage.Role.ToString())
	switch lastMessage.Role {
	case gollm.RoleToolCall:
		// parse the tool call
		switch lastMessage.ToolName {
		case toolVectorQuery:
			// get the summary llm
			var summaryLLMId pgtype.UUID
			summaryLLMId.Scan(args.ChatModelId)
			summaryLLM, err := llm.GetLLM(ctx, db, c.ID, summaryLLMId)
			if err != nil {
				return nil, slogger.Error(ctx, logger, "failed to get the summary llm", err)
			}

			/// RUN THE VECTOR STORE QUERY
			vecResponse, err := runToolVectorQuery(ctx, c, lastMessage, summaryLLM, db)
			if err != nil {
				return nil, fmt.Errorf("failed to run the vector query: %s", err)
			}

			// save the message
			if err := conv.SaveMessage(ctx, db, chatLLM, vecResponse.message); err != nil {
				return nil, fmt.Errorf("failed to save the message")
			}

			return &ragResponse{
				ConversationId: conv.ID.String(),
				Message:        vecResponse.message,
			}, nil

		default:
			return nil, fmt.Errorf("invalid tool call, tool name not supported: %s", lastMessage.ToolName)
		}
	case gollm.RoleAI:
		// handle a normal request from the user
		message := gollm.NewUserMessage(args.Input)
		response, err := conv.Completion(ctx, db, chatLLM, message, tools, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to send the competion")
		}

		// report the usage records
		usageRecords = append(usageRecords, response.UsageRecord)
		if err := reportUsage(); err != nil {
			return nil, err
		}

		// compose the response
		return &ragResponse{
			ConversationId: conv.ID.String(),
			Message:        response.Message,
		}, nil
	case gollm.RoleToolResult:
		// send the completion with the current state of the conversation
		response, err := conv.Completion(ctx, db, chatLLM, nil, tools, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to send the competion")
		}

		// report the usage records
		usageRecords = append(usageRecords, response.UsageRecord)
		if err := reportUsage(); err != nil {
			return nil, err
		}

		// return
		return &ragResponse{
			ConversationId: conv.ID.String(),
			Message:        response.Message,
		}, nil

	default:
		// invalid conversation state
		return nil, fmt.Errorf("invalid conversation state. Last message role: %s", lastMessage.Role.ToString())
	}
}

type runVectorQueryResponse struct {
	message      *gollm.Message
	usageRecords []*tokens.UsageRecord
	docs         []*queries.Document
	pages        []*queries.WebsitePage
}

func runToolVectorQuery(
	ctx context.Context,
	customer *Customer,
	lastMessage *gollm.Message,
	summaryLLM *llm.LLM,
	db queries.DBTX,
) (*runVectorQueryResponse, error) {
	logger := customer.logger.With("func", "runVectorQuery")
	dmodel := queries.New(db)
	usageRecords := make([]*tokens.UsageRecord, 0)

	// parse the argument
	vecQuery := lastMessage.ToolArguments["vector_query"].(string)
	if vecQuery == "" {
		return nil, fmt.Errorf("failed to use the tool")
	}

	// simplify the vector query
	logger.InfoContext(ctx, "Simplifying vector query ...", "query", vecQuery)

	// ensure the arguments are present
	vectorQuery, exists := lastMessage.ToolArguments["vector_query"]
	if !exists {
		return nil, slogger.Error(ctx, logger, "the argument 'vector_query' does not exist", nil)
	}

	// get the simple query llm from the database
	logger.InfoContext(ctx, "Getting the simplify llm ...")
	tmp, err := dmodel.GetInteralLLM(ctx, "Vector Query Generator")
	if err != nil {
		return nil, fmt.Errorf("failed to get the simple query LLM: %s", err)
	}
	simpleQueryLLM := llm.FromObjects(&tmp.Llm, &tmp.AvailableModel)

	// run a single completion
	simpleQueryResponse, err := simpleQueryLLM.SingleCompletion(
		ctx, logger, customer.ID, prompts.RAG_SIMPLE_QUERY_SYSTEM_PROMPT,
		vectorQuery.(string),
	)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the simple query", err)
	}
	usageRecords = append(usageRecords, simpleQueryResponse.UsageRecord)

	// parse the response
	simpleQueries := strings.Split(simpleQueryResponse.Message.Message, ",")

	// get the embeddings
	embs := customer.GetEmbeddings(ctx)
	vectorResponses := make([]*vectorstore.QueryResponse, 0)

	for _, item := range simpleQueries {
		logger.InfoContext(ctx, "Running vector query ...", "query", item)

		// run all the simple queries against the vector store
		vectorResponse, err := vectorstore.Query(ctx, logger, db, &vectorstore.QueryInput{
			CustomerID: customer.ID,
			Embeddings: embs,
			Query:      item,
			K:          4,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query the vectorstore: %s", err)
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
		response, err := summaryLLM.Summarize(ctx, logger, customer.ID, bufStr)
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

	// create a new message record with all the metadata needed
	message := gollm.NewToolResultMessage(lastMessage.ToolUseID, lastMessage.ToolName, toolResponse)
	arguments := make(map[string]any)
	arguments["docs"] = docs
	arguments["pages"] = pages
	message.ToolArguments = arguments

	// add the usage records
	usageRecords = append(usageRecords, embs.GetUsageRecords()...)

	// compose the response
	return &runVectorQueryResponse{
		message:      message,
		usageRecords: usageRecords,
		docs:         docs,
		pages:        pages,
	}, nil
}

func getRAGTools() []*gollm.Tool {
	funcs := make([]*gollm.Tool, 0)
	funcs = append(funcs, &gollm.Tool{
		Title:       toolVectorQuery,
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
	})
	return funcs
}
