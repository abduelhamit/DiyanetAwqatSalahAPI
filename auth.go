package diyanet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const apiUrlLogin = apiUrlPrefix + "Auth/Login"
const apiUrlRefreshToken = apiUrlPrefix + "Auth/RefreshToken/%s"

// Config encapsulates the credentials (email and password) used to authenticate with Diyanet services.
type Config struct {
	// Email is the user's email address used for authentication.
	Email string

	// Password is the user's password used for authentication.
	Password string

	// AccessTokenExpiryDelta is the duration before token expiry when a token refresh should be attempted.
	//
	// If zero, a default of 45 minutes is used.
	AccessTokenExpiryDelta time.Duration
}

func (c *Config) accessTokenExpiryDelta() time.Duration {
	if c.AccessTokenExpiryDelta == 0 {
		return 45 * time.Minute
	}
	return c.AccessTokenExpiryDelta
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
	ctx          context.Context
	conf         *Config
	refreshToken string
}

// Token implements [oauth2.TokenSource].
func (t *tokenSource) Token() (*oauth2.Token, error) {
	const tokenErrorPrefix = errorPrefix + "unable to retrieve access token: "

	client := oauth2.NewClient(t.ctx, nil)
	defer client.CloseIdleConnections()

	if t.refreshToken != "" {
		token, err := t.refreshAccessToken(client)
		if err == nil {
			return token, nil
		}
		log.Println(err)
	}

	jsonData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    t.conf.Email,
		Password: t.conf.Password,
	}

	reqBody, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf(tokenErrorPrefix+"failed to marshal request body: %w", err)
	}

	token, err := t.requestAccessToken(
		client,
		"POST",
		apiUrlLogin,
		"application/json",
		bytes.NewBuffer(reqBody),
		tokenErrorPrefix)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *tokenSource) refreshAccessToken(client *http.Client) (*oauth2.Token, error) {
	const refreshAccessTokenErrorPrefix = errorPrefix + "unable to refresh access token: "
	return t.requestAccessToken(
		client,
		"GET",
		fmt.Sprintf(apiUrlRefreshToken, t.refreshToken),
		"",
		nil,
		refreshAccessTokenErrorPrefix)
}

func (t *tokenSource) requestAccessToken(
	client *http.Client,
	method string,
	url string,
	contentType string,
	body io.Reader,
	errorPrefix string) (*oauth2.Token, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%sfailed to make refresh token request: %w", errorPrefix, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%sreceived non-2xx status code: %s (%d)", errorPrefix, resp.Status, resp.StatusCode)
	}

	var result Result[struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%sfailed to decode response: %w", errorPrefix, err)
	}
	if !result.IsSuccess {
		return nil, fmt.Errorf("%sAPI error: %s", errorPrefix, result.Error)
	}

	t.refreshToken = result.Data.RefreshToken
	return &oauth2.Token{
		AccessToken: result.Data.AccessToken,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(t.conf.accessTokenExpiryDelta()),
	}, nil
}
