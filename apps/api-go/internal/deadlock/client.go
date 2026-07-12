package deadlock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"cursed-apple-stats/apps/api-go/internal/store"
)

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) GetMatchHistory(ctx context.Context, accountID int64) ([]store.PlayerMatchHistoryEntry, error) {
	endpoint, err := url.JoinPath(c.baseURL, fmt.Sprintf("/v1/players/%d/match-history", accountID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deadlock match history request failed: %s", resp.Status)
	}

	var entries []store.PlayerMatchHistoryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}