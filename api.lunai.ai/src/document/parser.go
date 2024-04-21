package document

import (
	"bytes"
	"fmt"

	"code.sajari.com/docconv/v2"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

// custom html parser that converts to markdown.
func ParseHTML(data []byte) (string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertBytes(data)
	if err != nil {
		return "", err
	}
	return string(markdown), nil
}

// Will attempt to parse the data from the supplied mimetype
func ParseDynamic(data []byte, mime string) (string, error) {
	resp, err := docconv.Convert(bytes.NewBuffer(data), mime, true)
	if err != nil {
		return "", fmt.Errorf("there was an issue parsing the data: %v", err)
	}
	return resp.Body, nil
}
