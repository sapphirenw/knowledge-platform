package webparse

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"go.dpb.io/sitemap/data"
	"go.dpb.io/sitemap/httputil"
)

// TODO -- improve with: https://github.com/gocolly/colly/blob/master/_examples/shopify_sitemap/shopify_sitemap.go

// Recursively parses a sitemap and returns a list of all the valid domains
// relating to the site's blacklist and whitelist rules
func ParseSitemap(
	ctx context.Context,
	l *slog.Logger,
	site *queries.Website,
	max int,
) ([]string, error) {
	logger := l.With("site", site.Domain)
	logger.InfoContext(ctx, "Parsing sitemap ...")

	// compose the regex lists for the comparison
	whitelist := make([]*regexp.Regexp, len(site.Whitelist))
	for i, item := range site.Whitelist {
		r, err := regexp.Compile(item)
		if err != nil {
			return nil, fmt.Errorf("REGEX: there was an issue parsing the regex: %v", err)
		}
		whitelist[i] = r
	}
	blacklist := make([]*regexp.Regexp, len(site.Blacklist))
	for i, item := range site.Blacklist {
		r, err := regexp.Compile(item)
		if err != nil {
			return nil, fmt.Errorf("REGEX: there was an issue parsing the regex: %v", err)
		}
		blacklist[i] = r
	}

	// create data structures to handle canceling when the stop confidition is met
	pages := make([]string, 0)

	// parse the sitemap
	err := httputil.DefaultFetcher.Fetch(
		fmt.Sprintf("%s://%s/sitemap.xml", site.Protocol, site.Domain),
		data.EntryCallbackFunc(func(entry data.Entry) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			switch eT := entry.(type) {
			case *data.Sitemap:
			case *data.URL:
				// add if url passes the checks
				if allowed := isURLAllowed(eT.Location, whitelist, blacklist); allowed {
					pages = append(pages, eT.Location)

					// set break statement
					if max != 0 && len(pages) >= max {
						return fmt.Errorf("max page limit hit")

					}
				}
			}
			return nil
		}),
	)

	// ignore the control flow error
	if err != nil && !strings.Contains(err.Error(), "max page limit hit") {
		return pages, fmt.Errorf("there was an error parsing the sitemaps: %v", err)
	}

	logger.InfoContext(ctx, "Successfully parsed the sitemap", "pages", len(pages))

	return pages, nil
}

// Ensures that the url is valid or not
func isURLAllowed(url string, whitelist, blacklist []*regexp.Regexp) bool {
	// Check against blacklist first
	for _, pattern := range blacklist {
		if pattern.MatchString(url) {
			return false // URL is blacklisted
		}
	}

	// If a whitelist is provided, check if the URL is whitelisted
	if len(whitelist) > 0 {
		for _, pattern := range whitelist {
			if pattern.MatchString(url) {
				return true // URL is whitelisted
			}
		}
		// URL is not in the whitelist
		return false
	}

	// If no whitelist is provided, allow by default
	return true
}
