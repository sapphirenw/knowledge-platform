package webparse

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type ScrapeHrefsArgs struct {
	MaxDepth          int  // max depth for own domain
	MaxDepthOther     int  // max depth for other domains
	AllowOtherDomains bool // whether to allow other domains
	Limit             int  // max number of urls that will be returned
}

// Scrapes a webpage looking for ahrefs to crawl instead of the sitemap.
// returns a list of domain names to further process
func ScrapeHrefs(
	ctx context.Context,
	logger *slog.Logger,
	site *queries.Website,
	args *ScrapeHrefsArgs,
) ([]string, error) {
	if args == nil {
		args = &ScrapeHrefsArgs{}
	}
	if args.MaxDepthOther < 2 {
		args.MaxDepthOther = 2
	}
	if args.Limit == 0 {
		args.Limit = 100
	}

	// create the lists
	whitelist, blacklist, err := createLists(site)
	if err != nil {
		return nil, fmt.Errorf("REGEX: there was an issue parsing the regex: %v", err)
	}

	// create the results
	result := make([]string, 0)

	// converter and scraper
	c := colly.NewCollector(
		colly.Async(false),
		colly.IgnoreRobotsTxt(),
	)

	// configurations
	c.DisableCookies()
	c.CheckHead = false
	extensions.RandomUserAgent(c)

	// limit to the passed domain
	if !args.AllowOtherDomains {
		c.AllowedDomains = []string{
			site.Domain,
			fmt.Sprintf("%s://%s", site.Protocol, site.Domain),
		}
	}

	// if a max depth was passed, set it
	if args.MaxDepth > 2 {
		c.MaxDepth = args.MaxDepth
	}

	// setup a limit for the max number of urls visited
	c.OnRequest(func(r *colly.Request) {
		if len(result) >= args.Limit {
			r.Abort()
		}
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		u := e.Request.AbsoluteURL(e.Attr("href"))

		// ignore mail links
		if strings.Contains(u, "mailto") {
			return
		}

		// do not allow much recursion on non-main domains
		if !strings.Contains(u, site.Domain) && e.Request.Depth > args.MaxDepthOther {
			return
		}

		// check white/black list
		if isURLAllowed(u, whitelist, blacklist) {
			e.Request.Visit(u)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if len(result) < args.Limit {
			result = append(result, r.Request.URL.String())
		}
	})

	// error handler
	c.OnError(func(r *colly.Response, err error) {
		logger.ErrorContext(ctx, "There was an issue scraping the url", "url", r.Request.URL, "statusCode", r.StatusCode, "error", err)
	})

	// run the href scraper
	c.Visit(fmt.Sprintf("%s://%s%s", site.Protocol, site.Domain, site.Path))
	c.Wait()

	// remove duplicates
	result = utils.RemoveDuplicates(result, func(val string) any {
		return val
	})

	return result, nil
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
