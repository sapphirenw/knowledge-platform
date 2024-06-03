package customer

import (
	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/vectorstore"
)

type generatePresignedUrlResponse struct {
	UploadUrl  string    `json:"uploadUrl"`
	Method     string    `json:"method"`
	DocumentId uuid.UUID `json:"documentId"`
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

type queryVectorStoreResponse struct {
	Documents    []*vectorstore.DocumentResponse    `json:"documents"`
	WebsitePages []*vectorstore.WebsitePageResponse `json:"websitePages"`
}
