package customer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/ltypes"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

const (
	toolVectorQuery = "query_user_information"
)

type ragRequest struct {
	// general params
	Input          string `json:"input"`
	ConversationId string `json:"conversationId"`
	CheckQuality   bool   `json:"checkQuality"`

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
	ConversationId string                 `json:"conversationId"`
	Documents      []*queries.Document    `json:"documents"`
	WebsitePages   []*queries.WebsitePage `json:"websitePages"`
	Message        *gollm.Message         `json:"message"`
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

	response, err := c.RAG(r.Context(), tx, &body)
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
	args *ragRequest,
) (*ragResponse, error) {
	logger := c.logger.With("function", "RAG")
	logger.InfoContext(ctx, "Beginning document retrieval pathway")
	// dmodel := queries.New(db)

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
	conv, err := llm.AutoConversation(
		ctx,
		logger,
		db,
		c.ID,
		chatLLM,
		args.ConversationId,
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
			/// RUN THE VECTOR STORE QUERY
			vecResponse, err := runVectorQuery(ctx, c, conv, lastMessage, db, chatLLM)
			if err != nil {
				return nil, fmt.Errorf("failed to run the vector query: %s", err)
			}

			return &ragResponse{
				ConversationId: conv.ID.String(),
				Documents:      vecResponse.vectorResponse.Documents,
				WebsitePages:   vecResponse.vectorResponse.WebsitePages,
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
	message        *gollm.Message
	usageRecords   []*tokens.UsageRecord
	vectorResponse *queries.QueryVectorStoreResponse
}

func runVectorQuery(
	ctx context.Context,
	customer *Customer,
	conv *llm.Conversation,
	lastMessage *gollm.Message,
	db queries.DBTX,
	lm *llm.LLM,
) (*runVectorQueryResponse, error) {
	logger := customer.logger.With("func", "runVectorQuery")
	// parse the argument
	vecQuery := lastMessage.ToolArguments["vector_query"].(string)
	if vecQuery == "" {
		return nil, fmt.Errorf("failed to use the tool")
	}

	logger.InfoContext(ctx, "Running vector query ...", "query", vecQuery)

	// get the embeddings
	embs := customer.GetEmbeddings(ctx)

	// run all the simple queries against the vector store
	vectorResponse, err := vectorstore.Query(ctx, &vectorstore.QueryAllInput{
		QueryInput: &vectorstore.QueryInput{
			CustomerId: customer.ID,
			Embeddings: embs,
			DB:         db,
			Query:      vecQuery,
			K:          4,
			Logger:     logger,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query the vectorstore: %s", err)
	}

	// create a string from the results
	responseBuffer := new(bytes.Buffer)
	for _, item := range vectorResponse.Vectors {
		if _, err := responseBuffer.WriteString(item.Raw); err != nil {
			logger.ErrorContext(ctx, "There was an issue writing to the buffer", "error", err)
		}
	}

	responseString := responseBuffer.String()
	if responseString == "" {
		responseString = "[No information was found]"
	}
	toolResponse := fmt.Sprintf("Query Response:\n%s", responseString)

	// save message to the conversation
	message := gollm.NewToolResultMessage(lastMessage.ToolUseID, lastMessage.ToolName, toolResponse)
	if err := conv.SaveMessage(ctx, db, lm, message); err != nil {
		return nil, fmt.Errorf("failed to save the message")
	}

	// compose the response
	return &runVectorQueryResponse{
		message:        message,
		usageRecords:   embs.GetUsageRecords(),
		vectorResponse: vectorResponse,
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
