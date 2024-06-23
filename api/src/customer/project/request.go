package project

import (
	"context"
)

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

type addIdeasRequest struct {
	Ideas          []*ProjectIdea `json:"ideas"`
	ConversationId string         `json:"conversationId"`
}

func (r addIdeasRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)
	if len(r.Ideas) == 0 {
		p["ideas"] = "cannot be empty"
	}
	return p
}

type createLinkedinPostConfigRequest struct {
	LinkedInPostId            string `json:"linkedInPostId"`
	MinSections               int    `json:"minSections"`
	MaxSections               int    `json:"maxSections"`
	NumDocuments              int    `json:"numDocuments"`
	NumWebsitePages           int    `json:"numWebsitePages"`
	LlmContentCenerationId    string `json:"llmContentCenerationId"`
	LlmVectorSummarizationId  string `json:"llmVectorSummarizationId"`
	LlmWebsiteSummarizationId string `json:"llmWebsiteSummarizationId"`
	LlmProofReadingId         string `json:"llmProofReadingId"`
}

func (r createLinkedinPostConfigRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)
	if r.MinSections == 0 {
		p["MinSections"] = "cannot be empty"
	}
	if r.MaxSections == 0 {
		p["MaxSections"] = "cannot be empty"
	}
	if r.NumDocuments == 0 {
		p["NumDocuments"] = "cannot be empty"
	}
	if r.NumWebsitePages == 0 {
		p["NumWebsitePages"] = "cannot be empty"
	}
	if r.LlmContentCenerationId == "" {
		p["LlmContentCenerationId"] = "cannot be empty"
	}
	if r.LlmVectorSummarizationId == "" {
		p["LlmVectorSummarizationId"] = "cannot be empty"
	}
	if r.LlmWebsiteSummarizationId == "" {
		p["LlmWebsiteSummarizationId"] = "cannot be empty"
	}
	if r.LlmProofReadingId == "" {
		p["LlmProofReadingId"] = "cannot be empty"
	}
	return p
}

type generateLinkedInPostRequest struct {
	ConversationId string `json:"conversationId"`
	Input          string `json:"input"`
}

func (r generateLinkedInPostRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)
	if r.Input == "" {
		p["input"] = "cannot be empty"
	}
	return p
}
