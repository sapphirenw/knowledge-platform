package project

import "context"

type generateIdeasRequest struct {
	ConversationId string `json:"conversationId"`
	Feedback       string `json:"feedback"`
	K              int    `json:"k"`
}

func (r generateIdeasRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)
	if r.K > 5 {
		p["k"] = "cannot be larger than 5"
	}
	return p
}
