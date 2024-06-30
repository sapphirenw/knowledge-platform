package tool

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jake-landersweb/gollm/v2/src/gollm"
	"github.com/jake-landersweb/gollm/v2/src/tokens"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type ToolType string

const (
	VectorQuery ToolType = "vector_query"
)

func GetToolType(input string) (ToolType, error) {
	switch input {
	case string(VectorQuery):
		return VectorQuery, nil
	default:
		return "", fmt.Errorf("invalid tool type: %s", input)
	}
}

type RunToolArgs struct {
	Database    queries.DBTX
	Customer    *queries.Customer
	LastMessage *gollm.Message
	ToolLLM     *llm.LLM
}

func (args *RunToolArgs) Validate() error {
	if args.Database == nil {
		return fmt.Errorf("the database cannot be nil")
	}
	if args.Customer == nil {
		return fmt.Errorf("the customer cannot be nil")
	}
	if args.LastMessage == nil {
		return fmt.Errorf("the last message cannot be nil")
	}
	if args.ToolLLM == nil {
		return fmt.Errorf("the llm cannot be nil")
	}
	return nil
}

type ToolResponse struct {
	Message      *gollm.Message
	UsageRecords []*tokens.UsageRecord
}

type Tool interface {
	GetSchema() *gollm.Tool
	Run(
		ctx context.Context,
		l *slog.Logger,
		args *RunToolArgs,
	) (*ToolResponse, error)
	GetType() ToolType
}

func NewTool(name ToolType) Tool {
	switch name {
	case VectorQuery:
		return newToolVectorQuery()
	}
	panic(fmt.Sprintf("FATAL: invalid place of program reached. name is invalid: %s", name))
}

func ToolsToGollm(tools []Tool) []*gollm.Tool {
	response := make([]*gollm.Tool, len(tools))
	for i, item := range tools {
		response[i] = item.GetSchema()
	}
	return response
}
