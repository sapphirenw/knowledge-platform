package document

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
)

type ParserHTML struct{}

// for html, convert to markdown and return the string
func (p *ParserHTML) Parse(data string) (string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(data)
	if err != nil {
		return "", err
	}
	return markdown, nil
}
