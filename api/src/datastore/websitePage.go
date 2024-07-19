package datastore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/textsplitter"
	"github.com/sapphirenw/ai-content-creation-api/src/webparse"
)

type WebsitePage struct {
	*queries.WebsitePage

	// cached data to reduce compute if needed
	raw      *bytes.Buffer // raw data
	metadata *bytes.Buffer // for holding the headers
	cleaned  *bytes.Buffer // data but cleaned
	logger   *slog.Logger
}

func NewWebsitePageFromWebsitePage(
	ctx context.Context,
	logger *slog.Logger,
	page *queries.WebsitePage,
) (*WebsitePage, error) {
	return &WebsitePage{WebsitePage: page, logger: logger}, nil
}

func (p *WebsitePage) GetRaw(ctx context.Context) (*bytes.Buffer, error) {
	if p.raw == nil {
		// scrape the page
		response, err := webparse.ScrapeSingle(ctx, p.logger, p.WebsitePage)
		if err != nil {
			return nil, fmt.Errorf("failed to scrape the page: %s", err)
		}

		// create a buffer
		buf := new(bytes.Buffer)
		_, err = buf.WriteString(response.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to write to the buffer: %s", err)
		}

		met := new(bytes.Buffer)
		enc, _ := json.Marshal(response.Header)
		if _, err := buf.Write(enc); err != nil {
			return nil, fmt.Errorf("failed to write the header: %s", err)
		}

		p.raw = buf
		p.metadata = met
	}

	return p.raw, nil
}

func (p *WebsitePage) GetCleaned(ctx context.Context) (*bytes.Buffer, error) {
	if p.cleaned != nil {
		return p.cleaned, nil
	}

	// read the raw data
	raw, err := p.GetRaw(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read the raw data: %s", err)
	}

	// clean the data
	cleaned, err := ParseHTML(raw.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to clean the data")
	}

	// create a new buffer
	buf := new(bytes.Buffer)
	_, err = buf.WriteString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("failed to write to the buffer")
	}
	p.cleaned = buf

	return p.cleaned, nil
}

func (p *WebsitePage) GetChunks(ctx context.Context) ([]string, error) {

	content, err := p.GetCleaned(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the cleaned content")
	}

	// chunk the content as a markdown doc
	splitter := textsplitter.NewMarkdownTextSplitter(
		textsplitter.WithChunkSize(2000),
		textsplitter.WithChunkOverlap(200),
	)
	chunks, err := splitter.SplitText(content.String())
	if err != nil {
		return nil, fmt.Errorf("failed to split the text")
	}

	return chunks, nil
}

func (p *WebsitePage) GetMetadata(ctx context.Context) (*bytes.Buffer, error) {
	if p.metadata != nil {
		return p.metadata, nil
	}

	if _, err := p.GetRaw(ctx); err != nil {
		return nil, fmt.Errorf("failed to get the raw data to fetch the headers: %s", err)
	}

	if p.metadata == nil {
		return nil, fmt.Errorf("unknown error getting the metadata")
	}

	return p.metadata, nil
}

func (p *WebsitePage) GetSha256() (string, error) {
	return p.Sha256, nil
}

func (p *WebsitePage) getSummary() string {
	if p.Summary == "" || p.Sha256 != p.SummarySha256 {
		return ""
	}
	return p.Summary
}

func (p *WebsitePage) setSummary(s string) error {
	p.Summary = s
	p.SummarySha256 = p.Sha256
	return nil
}
