package project

import (
	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/llm"
)

type generateIdeasResponse struct {
	Ideas          []*ProjectIdea `json:"ideas"`
	ConversationId uuid.UUID      `json:"conversationId"`
}

type generateLinkedInPostResponse struct {
	ConversationId uuid.UUID                  `json:"conversationId"`
	Messages       []*llm.ConversationMessage `json:"messages"`
	LatestMessage  string                     `json:"latestMessage"`
}
