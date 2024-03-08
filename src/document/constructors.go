package document

import (
	"fmt"
)

func NewDoc(filename string, data []byte) (*Doc, error) {
	filetype, err := parseFileType(filename)
	if err != nil {
		return nil, fmt.Errorf("there was an issue parsing the filetype: %v", err)
	}

	return &Doc{
		Filename: filename,
		Filetype: filetype,
		Data:     data,
	}, nil
}
