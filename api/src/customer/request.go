package customer

import (
	"context"
	"time"
)

type generatePresignedUrlRequest struct {
	ParentId  *string `json:"parentId,omitempty"`
	Filename  string  `json:"filename"`
	Mime      string  `json:"mime"`
	Signature string  `json:"signature"`
	Size      int64   `json:"size"`
}

func (r generatePresignedUrlRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Filename == "" {
		p["filename"] = "cannot be empty"
	}
	if r.Mime == "" {
		p["mime"] = "cannot be empty"
	}
	if r.Signature == "" {
		p["signature"] = "cannot be empty"
	}
	if r.Size == 0 {
		p["size"] = "cannot be 0"
	}
	return p
}

type createFolderRequest struct {
	Owner *string `json:"owner,omitempty"`
	Name  string  `json:"name"`
}

func (r *createFolderRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Name == "" {
		p["name"] = "cannot be empty"
	}
	return p
}

type handleWebsiteRequest struct {
	Domain            string   `json:"domain"`
	Blacklist         []string `json:"blacklist"`
	Whitelist         []string `json:"whitelist"`
	UseSitemap        bool     `json:"useSitemap"`
	AllowOtherDomains bool     `json:"allowOtherDomains"`
	Pages             []string `json:"pages,omitempty"` // only included on the insert request
}

func (r handleWebsiteRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Domain == "" {
		p["domain"] = "cannot be empty"
	}
	return p
}

type insertSingleWebsitePageRequest struct {
	Domain string `json:"domain"`
}

func (r insertSingleWebsitePageRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Domain == "" {
		p["domain"] = "cannot be empty"
	}
	return p
}

type purgeDatastoreRequest struct {
	Timestamp *string `json:"timestamp"`
}

func (r purgeDatastoreRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Timestamp != nil {
		// ensure the correct format was passed
		_, err := time.Parse("2006-01-02 15:04:05", *r.Timestamp)
		if err != nil {
			p["timestamp"] = "The timestamp:" + *r.Timestamp + "is not valid"
		}
	}
	return p
}

type createCustomerRequest struct {
	Name string `json:"name"`
}

func (r createCustomerRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Name == "" {
		p["name"] = "cannot be empty"
	}

	return p
}

type queryVectorStoreRequest struct {
	Query          string `json:"query"`
	K              int    `json:"k"`
	IncludeContent bool   `json:"includeContent"`
}

func (r queryVectorStoreRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Query == "" {
		p["query"] = "cannot be empty"
	}
	if r.K == 0 || r.K > 100 {
		p["k"] = "has to be between 1 and 5"
	}

	return p
}

type createVectorRequest struct {
	Documents bool `json:"documents"`
	Websites  bool `json:"websites"`
}

func (r createVectorRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	return p
}

type updateCustomerLLMConfigsRequest struct {
	SummaryLLMID string `json:"summaryLLMId"`
	ChatLLMID    string `json:"chatLLMId"`
}

func (r updateCustomerLLMConfigsRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	return p
}

type createModelRequest struct {
	AvailableModelName string  `json:"availableModelName"`
	Title              string  `json:"title"`
	Temperature        float64 `json:"temperature"`
	Instructions       string  `json:"instructions"`
}

func (r createModelRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.AvailableModelName == "" {
		p["AvailableModelName"] = "Cannot be empty"
	}
	if r.Title == "" {
		p["Title"] = "Cannot be empty"
	}

	return p
}
