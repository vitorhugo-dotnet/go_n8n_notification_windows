package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// State represents the alert payload returned by the n8n endpoint.
type State struct {
	Enabled  bool   `json:"enabled"`
	Priority string `json:"priority"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Sound    string `json:"sound"`
}

// Fetcher retrieves alert state from an HTTP endpoint.
type Fetcher struct {
	client      *http.Client
	endpointURL string
}

// NewFetcher creates a Fetcher with a 10-second timeout.
func NewFetcher(endpointURL string) *Fetcher {
	return &Fetcher{
		client:      &http.Client{Timeout: 10 * time.Second},
		endpointURL: endpointURL,
	}
}

// Fetch performs a GET request and decodes the alert state.
func (f *Fetcher) Fetch(ctx context.Context) (*State, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var s State
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}
	return &s, nil
}
