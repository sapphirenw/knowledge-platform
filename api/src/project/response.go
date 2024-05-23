package project

import "github.com/google/uuid"

type generateIdeasResponse struct {
	Ideas          []*ProjectIdea `json:"ideas"`
	ConversationId uuid.UUID      `json:"conversationId"`
}
