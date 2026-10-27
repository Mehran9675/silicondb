package data

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// maxFileSize bounds a single download, so a misconfigured source can't
// exhaust memory.
const maxFileSize = 32 << 20

// Client downloads files from the remote database.
type Client struct {
	base string
	http *http.Client
}

// NewClient returns a client for the database rooted at base.
func NewClient(base string) *Client {
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return &Client{base: base, http: &http.Client{Timeout: 30 * time.Second}}
}

// Fetch downloads path relative to the database root. version is appended
// as a query parameter so CDNs don't serve a stale copy of a changed file.
func (c *Client) Fetch(ctx context.Context, path, version string) ([]byte, error) {
	u := c.base + path
	if version != "" {
		u += "?v=" + url.QueryEscape(version)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", path, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFileSize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxFileSize {
		return nil, fmt.Errorf("%s: larger than %d bytes", path, maxFileSize)
	}
	return body, nil
}
