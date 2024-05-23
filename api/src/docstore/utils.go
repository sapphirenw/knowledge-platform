package docstore

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func fileIdFromDoc(doc *queries.Document) string {
	return fileIdFromRaw(doc.CustomerID, doc.ID, doc.Filename)
}

func fileIdFromRaw(customerId uuid.UUID, documentId uuid.UUID, filename string) string {
	return fmt.Sprintf("%s/%s_%s", customerId.String(), documentId.String(), filename)
}
