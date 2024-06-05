package datastore

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/webparse"
)

type WebsitePage struct {
	*queries.WebsitePage

	// cached data to reduce compute if needed
	raw     *bytes.Buffer // raw data
	cleaned *bytes.Buffer // data but cleaned
	logger  *slog.Logger
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

		p.raw = buf
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

func (p *WebsitePage) GetSha256() (string, error) {
	return p.Sha256, nil
}
