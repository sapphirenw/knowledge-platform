package resume

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
	"github.com/sapphirenw/ai-content-creation-api/src/slogger"
)

type ResumeWorkExperience struct {
	*queries.ResumeWorkExperience
}

type CreateExperienceArgs struct {
	Company   string `json:"company"`
	Position  string `json:"position"`
	Location  string `json:"location"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	IsCurrent bool   `json:"isCurrent"`
	Index     int    `json:"index"`

	parsedStart pgtype.Timestamp
	parsedEnd   pgtype.Timestamp
}

func (args *CreateExperienceArgs) Validate() error {
	s, err := time.Parse(time.RFC3339, args.StartDate)
	if err != nil {
		return fmt.Errorf("failed to parse the start time")
	}
	var start pgtype.Timestamp
	if err := start.Scan(s); err != nil {
		return fmt.Errorf("failed to parse start time into postgres")
	}
	args.parsedStart = start

	if args.EndDate != "" {
		e, err := time.Parse(time.RFC3339, args.EndDate)
		if err != nil {
			return fmt.Errorf("failed to parse the end time")
		}
		var end pgtype.Timestamp
		if err := end.Scan(e); err != nil {
			return fmt.Errorf("failed to parse the end time into postgres")
		}
		args.parsedEnd = end
	}
	return nil
}

func (c *Client) CreateWorkExperience(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
	args *CreateExperienceArgs,
) (*queries.ResumeWorkExperience, *slogger.Err) {
	logger := l.With("desc", "create resume work experience")
	dmodel := queries.New(db)

	experience, err := dmodel.CreateResumeWorkExperience(ctx, &queries.CreateResumeWorkExperienceParams{
		ResumeID:  c.Resume.ID,
		Company:   args.Company,
		Position:  args.Position,
		Location:  args.Location,
		StartDate: args.parsedStart,
		EndDate:   args.parsedEnd,
		IsCurrent: args.IsCurrent,
		Index:     int32(args.Index),
	})
	if err != nil {
		return nil, slogger.NewErr(logger, "failed to create the work experience", err, "parse")
	}

	return experience, nil
}

func (c *Client) GetWorkExperiences(
	ctx context.Context,
	l *slog.Logger,
	db queries.DBTX,
) ([]*ResumeWorkExperience, *slogger.Err) {
	logger := l.With("desc", "get resume work experiences")
	dmodel := queries.New(db)

	experiences, err := dmodel.GetResumeWorkExperiences(ctx, c.Resume.ID)
	if err != nil {
		return nil, slogger.NewErr(logger, "failed to get the work experiences", err, "parse")
	}

	e := make([]*ResumeWorkExperience, 0)
	for _, item := range experiences {
		e = append(e, &ResumeWorkExperience{item})
	}

	return e, nil
}
