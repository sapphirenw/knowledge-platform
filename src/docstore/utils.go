package docstore

import (
	"fmt"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func createUniqueFileId(customer *queries.Customer, filename string, parentId *int64) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d", customer.ID))
	if parentId != nil {
		builder.WriteString(fmt.Sprintf("/%d", *parentId))
	} else {
		builder.WriteString("/nil")
	}
	builder.WriteString("/" + filename)
	return builder.String()
}

func parseUniqueFileId(fileId string) string {
	fmt.Println(fileId)
	parts := strings.Split(fileId, "/")
	fmt.Println(parts)
	parts = parts[2:]
	return strings.Join(parts, "/")
}
