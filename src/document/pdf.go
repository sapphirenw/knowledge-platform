package document

import (
	"bytes"
	"strings"

	"github.com/dslipak/pdf"
)

type ParserPDF struct{}

func (p *ParserPDF) Parse(data string) (string, error) {
	preader, err := pdf.NewReader(strings.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	plain, err := preader.GetPlainText()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(plain)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
