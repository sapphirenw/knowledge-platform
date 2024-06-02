package webparse

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly/v2"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func Scrape(
	ctx context.Context,
	logger *slog.Logger,
	site *queries.Website,
) (*map[string][]byte, error) {

	// create a map for the results to be piped to along with the mutex
	res := make(map[string][]byte)
	var mu sync.Mutex

	// converter and scraper
	converter := md.NewConverter("", true, nil)
	scraper := colly.NewCollector(
		colly.AllowedDomains(site.Domain),
		colly.Async(true),
	)
	scraper.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	// Find and visit all links
	scraper.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// parse the bodies from the webpage
	scraper.OnHTML("html", func(e *colly.HTMLElement) {
		// normalize url
		url, err := normalizeURL(e.Request.URL)
		if err != nil {
			logger.ErrorContext(ctx, "Error normalizing the url", "error", err)
			return
		}

		// check if this url was already processed
		if _, exists := res[url]; exists {
			return
		}

		// parse the markdown from the html elements
		markdown, err := converter.ConvertBytes(e.Response.Body)
		if err != nil {
			logger.ErrorContext(ctx, "Error parsing the markdown from the html", "error", err)
			return
		}

		// add to the map
		mu.Lock()
		res[url] = markdown
		mu.Unlock()
	})

	// visit all entries in a parsed sitemap
	scraper.OnXML("loc", func(x *colly.XMLElement) {
		x.Request.Visit(x.Attr("loc"))
	})

	// error handler
	scraper.OnError(func(r *colly.Response, err error) {
		logger.ErrorContext(ctx, "There was an issue scraping the url", "url", r.Request.URL, "statusCode", r.StatusCode)
	})

	scraper.OnRequest(func(r *colly.Request) {
		logger.DebugContext(ctx, "Visiting url", "url", r.URL)
	})

	scraper.Visit(fmt.Sprintf("%s://%s", site.Protocol, site.Domain))
	scraper.Wait() // with async = true

	return &res, nil
}

func ScrapeSingle(
	ctx context.Context,
	logger *slog.Logger,
	page *queries.WebsitePage,
) (*ScrapeResponse, error) {
	var res string
	header := ScrapeHeader{}
	var err error

	// converter and scraper
	converter := md.NewConverter("", true, nil)
	scraper := colly.NewCollector()

	// parse the bodies from the webpage
	scraper.OnHTML("body", func(e *colly.HTMLElement) {
		// parse the response
		raw, err := e.DOM.Html()
		if err != nil {
			logger.ErrorContext(ctx, "Error parsing the html", "error", err)
			return
		}

		// parse markdown from the passed html element
		markdown, err := converter.ConvertString(raw)
		if err != nil {
			logger.ErrorContext(ctx, "Error parsing the markdown from the html", "error", err)
			return
		}

		// return
		res = markdown
	})

	// parse the header
	scraper.OnHTML("head", func(e *colly.HTMLElement) {
		header.Title = e.ChildText("title")
		header.Description = e.ChildAttr(`meta[name="description"]`, "content")

		// get the keywords
		keywordsRaw := e.ChildAttr(`meta[name="keywords"]`, "content")
		keywords := make([]string, 0)
		keywords = append(keywords, strings.Split(keywordsRaw, ",")...)
		header.Tags = keywords
	})

	// error handler
	scraper.OnError(func(r *colly.Response, err error) {
		logger.ErrorContext(ctx, "There was an issue scraping the url", "url", r.Request.URL, "statusCode", r.StatusCode)
	})

	scraper.OnRequest(func(r *colly.Request) {
		logger.DebugContext(ctx, "Visiting url", "url", r.URL)
	})

	scraper.Visit(page.Url)

	return &ScrapeResponse{
		Header:  &header,
		Content: res,
	}, err
}

func normalizeURL(u *url.URL) (string, error) {
	// Ensure the path ends with a "/" if it's not empty and doesn't already have one.
	if u.Path != "" && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	// Normalize the URL by ensuring it has a trailing slash if it has no path.
	if u.Path == "" {
		u.Path = "/"
	}

	// Remove the default port for http and https schemes.
	if (u.Scheme == "http" && u.Port() == "80") || (u.Scheme == "https" && u.Port() == "443") {
		u.Host = strings.Split(u.Host, ":")[0]
	}

	return u.String(), nil
}
