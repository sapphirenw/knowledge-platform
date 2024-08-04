package resume

import (
	"context"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

func (c *Client) AttachWebsite(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	websiteId string,
) error {
	logger := l.With("desc", "attach website to resume")

	// parse the id
	id, err := utils.GoogleUUIDFromString(websiteId)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to parse the website id", err)
	}

	// ensure the website exists
	dmodel := queries.New(db)
	if _, err := dmodel.GetWebsite(ctx, id); err != nil {
		return slogger.Error(ctx, logger, "the website does not exist", err)
	}

	// create the relationship
	if _, err := dmodel.CreateResumeWebsite(ctx, &queries.CreateResumeWebsiteParams{
		ResumeID:  c.Resume.ID,
		WebsiteID: id,
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to attach the website to the resume", err)
	}

	logger.Info("successfully added website to the resume")

	return nil
}

func (c *Client) GetWebsites(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) ([]*queries.Website, error) {
	logger := l.With("desc", "get resume websites")
	dmodel := queries.New(db)

	sites, err := dmodel.GetResumeWebsites(ctx, c.Resume.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the documents", err)
	}

	return sites, nil
}

func (c *Client) AttachWebsitePage(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	websitePageId string,
) error {
	logger := l.With("desc", "attach website page to resume")

	// parse the id
	id, err := utils.GoogleUUIDFromString(websitePageId)
	if err != nil {
		return slogger.Error(ctx, logger, "failed to parse the website page id", err)
	}

	// ensure the website exists
	dmodel := queries.New(db)
	if _, err := dmodel.GetWebsitePage(ctx, id); err != nil {
		return slogger.Error(ctx, logger, "the website page does not exist", err)
	}

	// create the relationship
	if _, err := dmodel.CreateResumeWebsitePage(ctx, &queries.CreateResumeWebsitePageParams{
		ResumeID:      c.Resume.ID,
		WebsitePageID: id,
	}); err != nil {
		return slogger.Error(ctx, logger, "failed to attach the website page to the resume", err)
	}

	logger.Info("successfully added website page to the resume")

	return nil
}

func (c *Client) GetWebsitePages(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) ([]*queries.WebsitePage, error) {
	logger := l.With("desc", "get resume website pages")
	dmodel := queries.New(db)

	pages, err := dmodel.GetResumeWebsitePages(ctx, c.Resume.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the documents", err)
	}

	return pages, nil
}
