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

func (b generatePresignedUrlRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if b.Filename == "" {
		p["filename"] = "cannot be empty"
	}
	if b.Mime == "" {
		p["mime"] = "cannot be empty"
	}
	if b.Signature == "" {
		p["signature"] = "cannot be empty"
	}
	if b.Size == 0 {
		p["size"] = "cannot be 0"
	}
	return p
}

type createFolderRequest struct {
	Owner *queries.Folder
	Name  string
}

func (b *createFolderRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string, 0)
	if b.Owner == nil {
		p["owner"] = "cannot be nil"
	}
	if b.Name == "" {
		p["name"] = "cannot be empty"
	}
	return p
}
