package customer

import "github.com/sapphirenw/ai-content-creation-api/src/queries"

type CreateFolderArgs struct {
	Owner *queries.Folder
	Name  string
}

type FolderContents struct {
	Self      *queries.Folder     `json:"self"`
	Folders   []*queries.Folder   `json:"folders"`
	Documents []*queries.Document `json:"documents"`
}

// Provides metadata to the user about how to upload the document that was sent for upload
type GeneratePresignedUrlResponse struct {
	UploadUrl string
	Method    string
}
