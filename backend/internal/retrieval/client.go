package retrieval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"tutorai/backend/config"
)

// ServiceError is returned when the data service is unreachable or returns an
// unexpected status code.
type ServiceError struct {
	StatusCode int
	Message    string
}

func (e *ServiceError) Error() string {
	if e.StatusCode == 0 {
		return fmt.Sprintf("data service unreachable: %s", e.Message)
	}
	return fmt.Sprintf("data service error %d: %s", e.StatusCode, e.Message)
}

// Result is a single item returned by the data service.
type Result struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

// CardFilters holds the structured filters for a card retrieval request.
type CardFilters struct {
	ColorIdentity []string
	Format        string
	MaxPriceUSD   float64
	Query         string
}

// Client wraps the rag-data-service HTTP API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new data service Client from the provided config.
func NewClient(cfg config.Config, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    cfg.DataServiceURL,
		apiKey:     cfg.DataServiceAPIKey,
	}
}

// RetrieveCards queries the data service for cards matching the given filters.
func (c *Client) RetrieveCards(ctx context.Context, filters CardFilters) ([]Result, error) {
	f := map[string]any{}
	if len(filters.ColorIdentity) > 0 {
		f["color_identity"] = filters.ColorIdentity
	}
	if filters.Format != "" {
		f["format"] = filters.Format
	}
	if filters.MaxPriceUSD > 0 {
		f["max_price_usd"] = filters.MaxPriceUSD
	}

	body := map[string]any{
		"corpus": "cards",
		"query":  filters.Query,
		"top_k":  10,
	}
	if len(f) > 0 {
		body["filters"] = f
	}

	return c.post(ctx, body)
}

// RetrieveRules queries the data service for rules chunks relevant to the query.
func (c *Client) RetrieveRules(ctx context.Context, query string) ([]Result, error) {
	return c.post(ctx, map[string]any{
		"corpus": "rules",
		"query":  query,
		"top_k":  5,
	})
}

// RetrieveSlang queries the data service for slang and terminology matching the query.
func (c *Client) RetrieveSlang(ctx context.Context, query string) ([]Result, error) {
	return c.post(ctx, map[string]any{
		"corpus": "slang",
		"query":  query,
		"top_k":  5,
	})
}

func (c *Client) post(ctx context.Context, body map[string]any) ([]Result, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal retrieve request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/retrieve/tutorai", bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("create retrieve request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ServiceError{Message: err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, &ServiceError{StatusCode: resp.StatusCode, Message: resp.Status}
	}

	var result struct {
		Results []Result `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode retrieve response: %w", err)
	}

	return result.Results, nil
}
