package webscrape

import (
	"context"
	"testing"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/require"
)

func TestSitemap(t *testing.T) {
	logger := utils.DefaultLogger()

	site := queries.Website{
		CustomerID: 1,
		Protocol:   "https",
		Domain:     "crosschecksports.com",
	}

	// parse based on all pages
	pages, err := ParseSitemap(context.TODO(), logger, &site)
	if err != nil {
		t.Error(err)
	}
	total := len(pages)

	// add a blacklist
	site.Blacklist = append(site.Blacklist, "docs/.*$")
	pages, err = ParseSitemap(context.TODO(), logger, &site)
	if err != nil {
		t.Error(err)
	}
	black := len(pages)

	// add a whitelist
	site.Blacklist = []string{}
	site.Whitelist = append(site.Whitelist, "docs/.*$")
	pages, err = ParseSitemap(context.TODO(), logger, &site)
	if err != nil {
		t.Error(err)
	}
	white := len(pages)

	// asserts
	require.Equal(t, 31, total)
	require.Equal(t, total, black+white)

	// add both
	site.Blacklist = []string{"create.*$"}
	site.Whitelist = []string{"docs/.*$"}
	pages, err = ParseSitemap(context.TODO(), logger, &site)
	if err != nil {
		t.Error(err)
	}
	require.Less(t, len(pages), white)
	require.Less(t, len(pages), black)
}

func TestWebscrape(t *testing.T) {
	logger := utils.DefaultLogger()
	site := queries.Website{
		CustomerID: 1,
		Protocol:   "https",
		Domain:     "crosschecksports.com",
	}

	res, err := Scrape(context.TODO(), logger, &site)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("res is nil")
		return
	}
	require.NotEmpty(t, res)
	require.Equal(t, 11, len(*res))
}

func TestScrapeSingle(t *testing.T) {
	logger := utils.DefaultLogger()
	page := &queries.WebsitePage{
		ID:         1,
		CustomerID: 1,
		WebsiteID:  1,
		Url:        "https://crosschecksports.com/docs/create-team",
		Sha256:     "",
		IsValid:    true,
	}

	content, err := ScrapeSingle(context.TODO(), logger, page)
	require.NoError(t, err)
	require.NotEmpty(t, content)
}
