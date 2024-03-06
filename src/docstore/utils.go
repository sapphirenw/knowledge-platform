package docstore

import (
	"fmt"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func createUniqueFileId(customer *queries.Customer, filename string) string {
	return fmt.Sprintf("%d_%s", customer.ID, filename)
}

func parseUniqueFileId(customer *queries.Customer, fileId string) string {
	return strings.ReplaceAll(fileId, fmt.Sprintf("%d_", customer.ID), "")
}
