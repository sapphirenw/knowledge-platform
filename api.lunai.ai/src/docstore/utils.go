package docstore

import (
	"fmt"
	"strings"
)

func createUniqueFileId(customerId int64, filename string, parentId *int64) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d", customerId))
	if parentId != nil && *parentId != 0 {
		builder.WriteString(fmt.Sprintf("/%d", *parentId))
	} else {
		builder.WriteString("/nil")
	}
	builder.WriteString("/" + filename)
	return builder.String()
}
