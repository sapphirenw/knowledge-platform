package document

import (
	"fmt"
	"os"
	"testing"
)

func TestParseTXT(t *testing.T) {
	err := parseFile("file1.txt")
	if err != nil {
		t.Error(err)
	}
}

func TestParsePDF(t *testing.T) {
	err := parseFile("pdf.pdf")
	if err != nil {
		t.Error(err)
	}
}

func TestParseHTML(t *testing.T) {
	err := parseFile("html.html")
	if err != nil {
		t.Error(err)
	}
}

func TestParseDOCX(t *testing.T) {
	err := parseFile("docx.docx")
	if err != nil {
		t.Error(err)
	}
}

func parseFile(filename string) error {
	data, err := os.ReadFile(fmt.Sprintf("../../resources/%s", filename))
	if err != nil {
		return fmt.Errorf("there was an issue opening the file: %v", err)
	}

	doc, err := NewDoc(filename, data)
	if err != nil {
		return fmt.Errorf("there was an issue creating the document: %v", err)
	}

	cleaned, err := doc.GetCleanedData()
	if err != nil {
		return fmt.Errorf("there was an issue cleaning the document: %v", err)
	}
	if len(cleaned) == 0 {
		return fmt.Errorf("there was no data returned from the cleaning step")
	}

	fmt.Println("PRE: ", len(doc.Data))
	fmt.Println("POST:", len([]byte(cleaned)))

	return nil
}
