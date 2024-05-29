package content

import (
	"context"
	"fmt"

	"github.com/sapphirenw/ai-content-creation-api/src/queries"
)

func NewLinkedInPost(
	ctx context.Context,
	db queries.DBTX,
	p *queries.Project,
	pIdea *queries.ProjectIdea,
	title string,
) (*queries.LinkedinPost, error) {
	dmodel := queries.New(db)

	// create the project library record
	_, err := dmodel.CreateProjectLibraryRecord(ctx, &queries.CreateProjectLibraryRecordParams{
		ProjectID:   p.ID,
		Title:       title,
		ContentType: string(ContentType_LinkedInPost),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the project library record: %s", err)
	}

	return nil, nil

}
