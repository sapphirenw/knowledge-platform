package customer

import (
	"context"
	"time"
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
	Owner int64
	Name  string
}

func (r *createFolderRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
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
