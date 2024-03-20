package customer

import (
	"context"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type generatePresignedUrlRequest struct {
	Filename  string `json:"filename"`
	Mime      string `json:"mime"`
	Signature string `json:"signature"`
	Size      int64  `json:"size"`
	ParentId  *int64 `json:"parentId,omitempty"`
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
	Owner *queries.Folder
	Name  string
}

func (r *createFolderRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Owner == nil {
		p["owner"] = "cannot be nil"
	}
	if r.Name == "" {
		p["name"] = "cannot be empty"
	}
	return p
}

type handleWebsiteRequest struct {
	Domain    string   `json:"domain"`
	Blacklist []string `json:"blacklist"`
	Whitelist []string `json:"whitelist"`
	Insert    bool     `json:"insert"`
}

func (r handleWebsiteRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if r.Domain == "" {
		p["domain"] = "cannot be empty"
	}
	return p
}
