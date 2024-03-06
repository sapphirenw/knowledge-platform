package document

import (
	"fmt"
	"strings"
)

type Filetype string

const (
	FT_none = ""

	FT_txt  = "txt"
	FT_md   = "md"
	FT_html = "html"
	FT_xml  = "xml"
	FT_csv  = "csv"
	FT_tsv  = "tsv"
	FT_pdf  = "pdf"

	FT_unknown = "unknown"
)

func parseFileType(filename string) (Filetype, error) {
	items := strings.Split(filename, ".")
	if len(items) < 2 {
		return FT_none, fmt.Errorf("there is no extension on this file")
	}

	// validate the extension
	ext := items[len(items)-1]
	var ft Filetype
	switch ext {
	case "text":
		fallthrough
	case "txt":
		ft = FT_txt
	case "markdown":
		fallthrough
	case "md":
		ft = FT_md
	case "htm":
		fallthrough
	case "html":
		ft = FT_html
	case "xml":
		ft = FT_xml
	case "csv":
		ft = FT_csv
	case "tsv":
		ft = FT_tsv
	case "pdf":
		ft = FT_pdf
	default:
		return FT_none, fmt.Errorf("invalid extension: %s", ext)
	}

	return ft, nil
}
