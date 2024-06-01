package websearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type SearchResponse struct {
	Query           string     `json:"query"`
	NumberOfResults int        `json:"number_of_results"`
	Results         []*Result  `json:"results"`
	InfoBoxes       []*InfoBox `json:"infoboxes"`
	Suggestions     []string   `json:"suggestions"`
}

type Result struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Url       string   `json:"url"`
	Engine    string   `json:"engine"`
	ParsedUrl []string `json:"parsed_url"`
	Engines   []string `json:"engines"`
	Positions []int    `json:"positions"`
	Score     float32  `json:"score"`
	Category  string   `json:"category"`
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

func Web(query string) (*SearchResponse, error) {
	endpoint, exists := os.LookupEnv("WEBSEARCH_ENDPOINT")
	if !exists {
		return nil, fmt.Errorf("the env variable `WEBSEARCH_ENDPOINT` is required")
	}

	// send the request
	u := fmt.Sprintf("%s/search?q=%s&format=json", endpoint, url.QueryEscape(query))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to send the search request: %s", err)
	}
	defer resp.Body.Close()

	// parse the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the body: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the response was not 200: %d - %s", resp.StatusCode, string(body))
	}

	// parse the json
	var response SearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse the json: %s", err)
	}

	return &response, err
}

func Image() {}
