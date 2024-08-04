package resume

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
	"github.com/sapphirenw/ai-content-creation-api/src/utils"
)

type Client struct {
	Resume   *queries.Resume
	Customer *queries.Customer
}

func NewClient(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	customer *queries.Customer,
	resumeId string,
) (*Client, error) {
	logger := l.With("resumeId", resumeId)

	// parse the id
	id, err := utils.GoogleUUIDFromString(resumeId)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to parse the resumeId", err)
	}

	dmodel := queries.New(db)
	resume, err := dmodel.GetResume(ctx, id)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		logger.Info("creating a new resume item")
		// create new item
		resume, err = dmodel.CreateCustomerResume(ctx, &queries.CreateCustomerResumeParams{
			ID:         customer.ID,
			CustomerID: customer.ID,
			Title:      "Default",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create the resume object")
		}
	} else if err != nil {
		// unknown error
		return nil, slogger.Error(ctx, logger, "failed to fetch the resume", err)
	}

	return &Client{
		Resume:   resume,
		Customer: customer,
	}, nil
}

func (c *Client) GetAbout(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) (*queries.ResumeAbout, error) {
	logger := l.With("desc", "get resume about")
	dmodel := queries.New(db)

	about, err := dmodel.GetResumeAbout(ctx, c.Resume.ID)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		logger.Info("no resume about found, creating a new one")
		// create a new one
		about, err = dmodel.CreateResumeAbout(ctx, &queries.CreateResumeAboutParams{
			ResumeID: c.Resume.ID,
		})
		if err != nil {
			return nil, slogger.Error(ctx, logger, "failed to create the resume about", err)
		}
	} else if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the resume about", err)
	}

	return about, nil
}

func (c *Client) SetTitle(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	title string,
) error {
	logger := l.With("desc", "set resume title")
	dmodel := queries.New(db)

	resume, err := dmodel.SetResumeTitle(ctx, &queries.SetResumeTitleParams{
		ID:    c.Resume.ID,
		Title: title,
	})
	if err != nil {
		return slogger.Error(ctx, logger, "failed to set the resume title", err)
	}

	// set internal field
	c.Resume = resume
	return nil
}

func (c *Client) SetAbout(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	args *queries.CreateResumeAboutParams,
) (*queries.ResumeAbout, error) {
	logger := l.With("desc", "set resume about")
	dmodel := queries.New(db)

	args.ResumeID = c.Resume.ID
	about, err := dmodel.CreateResumeAbout(ctx, args)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to set the about", err)
	}

	return about, err
}
