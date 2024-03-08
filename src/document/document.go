package document

import (
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

/*
Convenience powerful struct that parses raw data to provide consistency in how the
input data is parsed, and what some metadata about the file is.
*/
type Doc struct {
	Filename string
	Filetype Filetype
	Data     []byte
}

// Get the size this document takes up on disk in bytes
func (d *Doc) GetSizeInBytes() int {
	return len(d.Data)
}

/*
Runs the raw data stored in the document through a document parser to extract the raw
text data, then runs the extracted text through an additional parser to clean out whitespace
characters and weird formatting. The result is a string ready to be vectorized.
*/
func (d *Doc) GetCleanedData() (string, error) {
	var content string
	var err error
	// parse the contents based on the filetype
	switch d.Filetype {
	case FT_html:
		content, err = ParseHTML(d.Data)
	default:
		// use an auto-content detection parser
		content, err = ParseDynamic(d.Data, string(d.Filetype))
		// TODO -- handle errors. May potentially need to just use the raw content here
	}

	if err != nil {
		return "", fmt.Errorf("there was an issue parsing the document: %v", err)
	}

	// clean the string
	cleaned := utils.CleanInput(content)

	return cleaned, nil
}
