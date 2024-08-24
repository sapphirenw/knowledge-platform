package resume

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

type ChecklistItem struct {
	Completed bool   `json:"completed"`
	Message   string `json:"message"`
}

func (c *Client) GetChecklist(
	ctx context.Context,
	logger *slog.Logger,
	db queries.DBTX,
) ([]*ChecklistItem, error) {
	response := make([]*ChecklistItem, 0)

	// get the attached files
	docs, err := c.GetDocuments(ctx, logger, db)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the attached documents", err)
	}

	// get the attached pages
	pages, err := c.GetWebsitePages(ctx, logger, db)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the pages", err)
	}

	completed := false
	message := ""

	// check if the doc section is completed
	if len(docs) == 0 && len(pages) == 0 {
		message = "Please attach a document or website with some information about you."
	} else {
		completed = true
	}
	response = append(response, &ChecklistItem{Completed: completed, Message: message})

	// get the about
	about, err := c.GetAbout(ctx, logger, db)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get resume about", err)
	}

	message = ""
	completed = false

	// parse this page
	if about.Name == "" {
		message = "Please fill out your name."
	} else if about.Email == "" {
		message = "Please fill out your email."
	} else if about.Title == "" {
		message = "Please fill out your title."
	} else if about.Location == "" {
		message = "Please fill out your location."
	} else {
		completed = true
	}

	response = append(response, &ChecklistItem{Completed: completed, Message: message})

	// get all the root objects and check if the information fields are filled out
	experiences, err := c.GetWorkExperiences(ctx, logger, db)
	if err != nil {
		return nil, slogger.Error(ctx, logger, "failed to get the experiences", err)
	}

	completed = true
	message = ""

	// check if empty
	if len(experiences) == 0 {
		completed = false
		message = "Please add some work experience to your resume."
	} else {
		for _, item := range experiences {
			if strings.Trim(item.Information, "") == "" {
				completed = false
				message = fmt.Sprintf("The experience '%s' has no information!", item.Position)
				break
			}
		}
	}

	response = append(response, &ChecklistItem{Completed: completed, Message: message})

	return response, nil
}
