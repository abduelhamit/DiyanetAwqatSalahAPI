package diyanet

import (
	"context"
	"net/http"
)

// Client is a Diyanet Awqat Salah API client.
type Client struct {
	// httpClient is the HTTP client used to make requests.
	httpClient *http.Client
}

// NewClient creates a new Diyanet Awqat Salah API client using the provided configuration.
func (c *Config) NewClient(ctx context.Context) *Client {
	return &Client{
		httpClient: c.HTTPClient(ctx),
	}
}
