package docstore

import (
	"context"
)

type RemoteDocstore interface {
	// Requests a pre-signed url a client can use to upload a file
	GeneratePresignedUrl(ctx context.Context, doc *Document) (string, error)

	// downloads the raw file contents from the remote docstore
	DownloadFile(ctx context.Context, uniqueId string) ([]byte, error)

	// deletes the file from the remote docstore
	DeleteFile(ctx context.Context, uniqueId string) error

	// deletes this key and all keys that are owned by this key (root folder)
	DeleteRoot(ctx context.Context, prefix string) error

	// returns the method the client should use for the pre-signed url
	GetUploadMethod() string
}

type UploadUrlInput struct {
	ParentId  *int64 `json:"parentId"`
	Filename  string `json:"filename"`
	Mime      string `json:"mime"`
	Signature string `json:"signature"`
	Size      int64  `json:"size"`
}
