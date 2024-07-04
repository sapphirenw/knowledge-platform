package customer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/customer/conversation"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/prompts"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/request"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/tool"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
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
	ragTools := getRAGTools()

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
		response, err := conv.Completion(ctx, db, chatLLM, message, tool.ToolsToGollm(ragTools), ragTools[0].GetSchema())
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
		toolType, err := tool.GetToolType(lastMessage.ToolName)
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to parse the tool name", err)
		}
		parsedTool := tool.NewTool(toolType)

		// get the summary llm
		summaryLLM, err := llm.GetLLMString(ctx, db, c.ID, args.SummaryModelId)
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to get the summary llm", err)
		}

		vecResponse, err := parsedTool.Run(ctx, logger, &tool.RunToolArgs{
			Database:    db,
			Customer:    c.Customer,
			LastMessage: lastMessage,
			ToolLLM:     summaryLLM,
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to run the tool", err)
		}

		// save the message
		if err := conv.SaveMessage(ctx, db, chatLLM, vecResponse.Message); err != nil {
			return nil, slogger.Error(ctx, logger, "failed to save the message", err)
		}

		return &ragResponse{
			ConversationId: conv.ID.String(),
			Message:        vecResponse.Message,
		}, nil
	case gollm.RoleAI:
		// handle a normal request from the user
		message := gollm.NewUserMessage(args.Input)
		response, err := conv.Completion(ctx, db, chatLLM, message, tool.ToolsToGollm(ragTools), nil)
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
		response, err := conv.Completion(ctx, db, chatLLM, nil, tool.ToolsToGollm(ragTools), nil)
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

func getRAGTools() []tool.Tool {
	tools := make([]tool.Tool, 0)
	tools = append(tools, tool.NewTool(tool.VectorQuery))
	return tools
}
