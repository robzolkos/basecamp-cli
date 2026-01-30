package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/rzolkos/basecamp-cli/internal/config"
)

const (
	UserAgent = "Basecamp CLI (https://github.com/rzolkos/basecamp-cli)"
	Timeout   = 30 * time.Second
)

type Client struct {
	token   string
	baseURL string
	http    *http.Client
}

func New() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	token, err := config.LoadToken()
	if err != nil {
		return nil, err
	}

	return &Client{
		token:   token,
		baseURL: cfg.APIBaseURL(),
		http:    &http.Client{Timeout: Timeout},
	}, nil
}

func (c *Client) Get(path string) (json.RawMessage, error) {
	url := c.resolveURL(path)
	return c.request(context.Background(), http.MethodGet, url, nil)
}

func (c *Client) Post(path string, data any) (json.RawMessage, error) {
	url := c.resolveURL(path)
	return c.request(context.Background(), http.MethodPost, url, data)
}

func (c *Client) Put(path string, data any) (json.RawMessage, error) {
	url := c.resolveURL(path)
	return c.request(context.Background(), http.MethodPut, url, data)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	url := c.resolveURL(path)
	return c.request(context.Background(), http.MethodDelete, url, nil)
}

// UploadFile uploads raw binary data to the given path
func (c *Client) UploadFile(path string, data []byte, contentType string, size int64) (json.RawMessage, error) {
	url := c.resolveURL(path)
	return c.uploadRequest(context.Background(), url, data, contentType, size)
}

func (c *Client) uploadRequest(ctx context.Context, url string, data []byte, contentType string, size int64) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp)
}

// GetAll fetches all pages of a paginated endpoint and returns combined results
func (c *Client) GetAll(path string) ([]json.RawMessage, error) {
	var results []json.RawMessage
	url := c.resolveURL(path)

	for url != "" {
		resp, nextURL, err := c.requestWithPagination(context.Background(), url)
		if err != nil {
			return nil, err
		}

		var page []json.RawMessage
		if err := json.Unmarshal(resp, &page); err != nil {
			return nil, err
		}
		results = append(results, page...)
		url = nextURL
	}

	return results, nil
}

func (c *Client) resolveURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return c.baseURL + path
}

func (c *Client) request(ctx context.Context, method, url string, data any) (json.RawMessage, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, data != nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp)
}

func (c *Client) requestWithPagination(ctx context.Context, url string) (json.RawMessage, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	c.setHeaders(req, false)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := c.handleResponse(resp)
	if err != nil {
		return nil, "", err
	}

	nextURL := parseNextLink(resp.Header.Get("Link"))
	return data, nextURL, nil
}

func (c *Client) setHeaders(req *http.Request, hasBody bool) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", UserAgent)
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}
}

func (c *Client) handleResponse(resp *http.Response) (json.RawMessage, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	if len(body) == 0 {
		return nil, nil
	}

	return body, nil
}

var linkRegex = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

func parseNextLink(header string) string {
	if header == "" {
		return ""
	}

	matches := linkRegex.FindStringSubmatch(header)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
