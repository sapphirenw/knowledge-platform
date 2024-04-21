package customer

import "github.com/sapphirenw/ai-content-creation-api/src/queries"

type generatePresignedUrlResponse struct {
	UploadUrl  string `json:"uploadUrl"`
	Method     string `json:"method"`
	DocumentId int64  `json:"documentId"`
}

type listFolderContentsResponse struct {
	Self      *queries.Folder     `json:"self"`
	Folders   []*queries.Folder   `json:"folders"`
	Documents []*queries.Document `json:"documents"`
}

type handleWebsiteResponse struct {
	Site  *queries.Website       `json:"site"`
	Pages []*queries.WebsitePage `json:"pages"`
}
