package commands

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type SearchCmd struct{}

type SearchResult struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	Type             string  `json:"type"`
	PlainTextContent string  `json:"plain_text_content"`
	CreatedAt        string  `json:"created_at"`
	AppURL           string  `json:"app_url"`
	Creator          Creator `json:"creator"`
	Bucket           struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"bucket"`
	Parent struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"parent"`
}

type SearchResultOutput struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Snippet string `json:"snippet,omitempty"`
	Project string `json:"project"`
	Creator string `json:"creator"`
	URL     string `json:"url"`
}

type SearchOutput struct {
	Query   string               `json:"query"`
	Results []SearchResultOutput `json:"results"`
}

func (c *SearchCmd) Run(args []string) error {
	// Parse query and flags
	var query, searchType, projectID string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--type":
			if i+1 < len(args) {
				searchType = args[i+1]
				i++
			}
		case "--project":
			if i+1 < len(args) {
				projectID = args[i+1]
				i++
			}
		default:
			if query == "" {
				query = args[i]
			}
		}
	}

	if query == "" {
		return errors.New("search query required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	// Build search URL
	params := url.Values{}
	params.Set("q", query)
	if searchType != "" {
		params.Set("type", searchType)
	}
	if projectID != "" {
		params.Set("bucket_id", projectID)
	}

	searchURL := "/search.json?" + params.Encode()
	data, err := cl.Get(searchURL)
	if err != nil {
		return err
	}

	var results []SearchResult
	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}

	output := SearchOutput{
		Query:   query,
		Results: make([]SearchResultOutput, len(results)),
	}

	for i, r := range results {
		output.Results[i] = SearchResultOutput{
			ID:      r.ID,
			Title:   r.Title,
			Type:    r.Type,
			Snippet: stripHTML(r.PlainTextContent),
			Project: r.Bucket.Name,
			Creator: r.Creator.Name,
			URL:     r.AppURL,
		}
	}

	return PrintJSON(output)
}
