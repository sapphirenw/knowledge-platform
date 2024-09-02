package resume

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

type ResumeProject struct {
	*queries.ResumeProject
}

type CreateProjectArgs struct {
	Index     int    `json:"index"`
	Title     string `json:"title"`
	SubTitle  string `json:"subtitle"`
	Link      string `json:"link"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`

	parsedStart pgtype.Timestamp
	parsedEnd   pgtype.Timestamp
}

func (args *CreateProjectArgs) Validate() error {
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
