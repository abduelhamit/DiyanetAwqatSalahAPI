package diyanet

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// Config encapsulates the credentials (email and password) used to authenticate with Diyanet services.
type Config struct {
	// Email is the user's email address used for authentication.
	Email string

	// Password is the user's password used for authentication.
	Password string
}

// Token uses client credentials to retrieve a token.
//
// The provided context optionally controls which HTTP client is used. See the [oauth2.HTTPClient] variable.
func (c *Config) Token(ctx context.Context) (*oauth2.Token, error) {
	return c.TokenSource(ctx).Token()
}

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary.
//
// The provided context optionally controls which HTTP client
// is returned. See the [oauth2.HTTPClient] variable.
//
// The returned [http.Client] and its Transport should not be modified.
func (c *Config) Client(ctx context.Context) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx))
}

// TokenSource returns a [oauth2.TokenSource] that returns t until t expires,
// automatically refreshing it as necessary using the provided context and the
// client ID and client secret.
//
// Most users will use [Config.Client] instead.
func (c *Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	source := &tokenSource{
		ctx:  ctx,
		conf: c,
	}
	return oauth2.ReuseTokenSource(nil, source)
}

type tokenSource struct {
	ctx  context.Context
	conf *Config
}

// Token implements [oauth2.TokenSource].
func (t *tokenSource) Token() (*oauth2.Token, error) {
	return nil, errors.New("not implemented")
}
