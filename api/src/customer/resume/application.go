package resume

import (
	"context"
	"log/slog"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

func (c *Client) GetApplications(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) ([]*queries.ResumeApplication, error) {
	logger := l.With("desc", "get resume applications")

	dmodel := queries.New(db)
	apps, err := dmodel.GetResumeApplications(ctx, c.Resume.ID)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the applications", err)
	}

	return apps, nil
}

func (c *Client) CreateApplication(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	args *createResumeApplicationRequest,
) (*queries.ResumeApplication, error) {
	logger := l.With("desc", "create resume application")

	if args == nil {
		return nil, slogger.Error(ctx, logger, "args cannnot be empty", nil)
	}

	logger.Info("Creating a resume application", "args", *args)

	dmodel := queries.New(db)
	app, err := dmodel.CreateResumeApplication(ctx, &queries.CreateResumeApplicationParams{
		ResumeID:    c.Resume.ID,
		Title:       args.Title,
		Link:        args.Link,
		CompanySite: args.CompanySite,
		RawText:     args.RawText,
	})
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to create the resume application", err)
	}

	return app, nil
}
