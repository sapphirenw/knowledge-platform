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

type UploadDocumentsResponse struct {
	Doc   *queries.Document
	Url   string
	Error error
}
