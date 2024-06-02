package webparse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func WebSearch(query string) (*SearchResponse, error) {
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

func ImageSearch(query string) (*SearchResponse, error) {
	return WebSearch(fmt.Sprintf("!images %s", query))
}
