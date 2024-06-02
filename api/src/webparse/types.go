package webparse

type ScrapeResponse struct {
	Header  *ScrapeHeader `json:"header"`
	Content string        `json:"content"`
}

type ScrapeHeader struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type SearchResponse struct {
	Query           string     `json:"query"`             // The query passed to the search engine
	NumberOfResults int        `json:"number_of_results"` // number of results parsed
	Results         []*Result  `json:"results"`           // list of result objects
	InfoBoxes       []*InfoBox `json:"infoboxes"`         // list of info boxes, usually from Wikipedia. Usually empty
	Suggestions     []string   `json:"suggestions"`       // suggested future search terms
}

type Result struct {
	Title     string   `json:"title"`      // title of the result
	Content   string   `json:"content"`    // string content description of the result | only available in web search
	Url       string   `json:"url"`        // url of the search result
	Engine    string   `json:"engine"`     // which engine this result comes from
	ParsedUrl []string `json:"parsed_url"` // list of url components
	Engines   []string `json:"engines"`    // list of engines this result appears with
	Positions []int    `json:"positions"`  // list of the search position for this result
	Score     float32  `json:"score"`      // extrapolation of search position across providers
	Category  string   `json:"category"`   // category of the result

	ImgSrc       string `json:"img_src"`              // endpoint of the image | only available in image search
	ThumbnailSrc string `json:"thumbnail_src"`        // endpoint for the thumbnail | only available in image search
	Resolution   string `json:"resolution,omitempty"` // resolution of the image. Not always available | only available in image search
	ImgFormat    string `json:"img_format,omitempty"` // image format. Not always available | only available in image search
}

type InfoBox struct {
	InfoBox string        `json:"infobox"`
	Id      string        `json:"id"`
	Content string        `json:"content"`
	Urls    []*InfoBoxUrl `json:"urls"`
	Engine  string        `json:"engine"`
	Engines []string      `json:"engines"`
}

type InfoBoxUrl struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type ParsedResult struct {
	Result     *Result
	ParsedData string
}
