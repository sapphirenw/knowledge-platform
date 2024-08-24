package resume

import "context"

type setResumeTitleRequest struct {
	Title string `json:"title"`
}

func (r setResumeTitleRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)

	if r.Title == "" {
		p["title"] = "cannot be empty"
	}

	return p
}

type createResumeApplicationRequest struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	CompanySite string `json:"companySite"`
	RawText     string `json:"rawText"`
}

func (r createResumeApplicationRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)

	if r.Title == "" {
		p["title"] = "cannot be empty"
	}
	if r.RawText == "" {
		p["rawText"] = "cannot be empty"
	}

	return p
}

type attachDocumentsRequest struct {
	DocumentIDs []string `json:"documentIds"`
}

func (r attachDocumentsRequest) Valid(ctx context.Context) map[string]string {
	p := make(map[string]string)

	if len(r.DocumentIDs) == 0 {
		p["documentIds"] = "cannot be empty"
	}

	return p
}
