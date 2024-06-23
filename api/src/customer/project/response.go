package project

import (
	"github.com/google/uuid"
	"github.com/jake-landersweb/gollm/v2/src/gollm"
)

type generateIdeasResponse struct {
	Ideas          []*ProjectIdea `json:"ideas"`
	ConversationId uuid.UUID      `json:"conversationId"`
}

type generateLinkedInPostResponse struct {
	ConversationId uuid.UUID        `json:"conversationId"`
	Messages       []*gollm.Message `json:"messages"`
	LatestMessage  string           `json:"latestMessage"`
}
