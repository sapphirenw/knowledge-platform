package document

import (
	"fmt"
	"strings"
)

type Filetype string

const (
	FT_none = ""

	// raw text formats
	FT_txt = "text/plain"
	FT_md  = "text/markdown"
	FT_csv = "text/csv"
	FT_tsv = "text/tab-separated-values"

	// image
	FT_pdf = "application/pdf"
	FT_png = "image/png"
	FT_jpg = "image/jpeg"

	// parsed datatypes
	FT_html = "text/html"
	FT_xml  = "application/xml"
	FT_doc  = "application/msword"
	FT_docx = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

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
	switch strings.ToLower(ext) {

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

	case "jpeg":
		fallthrough
	case "jpg":
		ft = FT_jpg

	case "doc":
		ft = FT_doc
	case "docx":
		ft = FT_docx

	case "png":
		ft = FT_png
	default:
		return FT_none, fmt.Errorf("invalid extension: %s", ext)
	}

	return ft, nil
}
