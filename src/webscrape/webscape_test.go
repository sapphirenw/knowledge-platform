package webscrape

import (
	"context"
	"testing"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/testingutils"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestSitemap(t *testing.T) {
	logger := utils.DefaultLogger()
	site := queries.Website{
		CustomerID: testingutils.TEST_CUSTOMER_ID,
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
	assert.Equal(t, 31, total)
	assert.Equal(t, total, black+white)

	// add both
	site.Blacklist = []string{"create.*$"}
	site.Whitelist = []string{"docs/.*$"}
	pages, err = ParseSitemap(context.TODO(), logger, &site)
	if err != nil {
		t.Error(err)
	}
	assert.Less(t, len(pages), white)
	assert.Less(t, len(pages), black)
}

func TestWebscrape(t *testing.T) {
	logger := utils.DefaultLogger()
	site := queries.Website{
		CustomerID: testingutils.TEST_CUSTOMER_ID,
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
	assert.NotEmpty(t, res)
	assert.Equal(t, 11, len(*res))
}
