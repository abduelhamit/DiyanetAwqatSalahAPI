package diyanet

import (
	"context"
	"net/http"
)

// Client is a Diyanet Awqat Salah API client.
type Client struct {
	// ctx is the context used for making requests.
	ctx context.Context
	// httpClient is the HTTP client used to make requests.
	httpClient *http.Client
}

// NewClient creates a new Diyanet Awqat Salah API client using the provided configuration.
func (c Config) NewClient(ctx context.Context) Client {
	return Client{
		ctx:        ctx,
		httpClient: c.HTTPClient(ctx),
	}
}

func (c Client) get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}
