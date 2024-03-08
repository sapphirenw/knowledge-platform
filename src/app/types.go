package app

import "github.com/sapphirenw/ai-content-creation-api/src/queries"

type CreateFolderArgs struct {
	Owner *queries.Folder
	Name  string
}

type FolderContents struct {
	Self      *queries.Folder
	Folders   []queries.Folder
	Documents []queries.Document
}

// Provides metadata to the user about how to upload the document that was sent for upload
type GeneratePresignedUrlResponse struct {
	UploadUrl string
	Method    string
}
