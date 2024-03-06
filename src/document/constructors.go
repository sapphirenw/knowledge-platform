package document

import (
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

func NewDocFromBytes(filename string, data []byte) (*Doc, error) {
	return NewDocFromString(filename, string(data))
}

func NewDocFromString(filename, data string) (*Doc, error) {
	// parse the filetype
	filetype, err := parseFileType(filename)
	if err != nil {
		return nil, err
	}

	var parser Parser
	// parse the contents based on the filetype
	switch filetype {
	case FT_pdf:
		parser = &ParserPDF{}
	case FT_html:
		parser = &ParserHTML{}
	default:
		parser = &ParserTxt{}
	}

	content, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("there was an issue parsing the document: %v", err)
	}

	// clean the string
	cleaned := utils.CleanInput(content)

	return &Doc{
		Filename: filename,
		Filetype: filetype,
		Data:     cleaned,
	}, nil
}
