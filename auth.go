package diyanet

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const apiURLLogin = apiURLPrefix + "Auth/Login"
const apiURLRefreshToken = apiURLPrefix + "Auth/RefreshToken/%s"
const retrieveTokenErrorPrefix = errorPrefix + "unable to retrieve access token: "
const refreshTokenErrorPrefix = errorPrefix + "unable to refresh access token: "

var earlyExpiry = 15 * time.Minute
var past time.Time

func init() {
	past = past.Add(earlyExpiry + 1)
}

// Token uses client credentials to retrieve a token.
//
// The provided context optionally controls which HTTP client is used. See the [oauth2.HTTPClient] variable.
func (c *Config) Token(ctx context.Context) (*oauth2.Token, error) {
	return c.TokenSource(ctx).Token()
}

// HTTPClient returns an HTTP client using the provided configuration.
// The token will auto-refresh as necessary.
//
// The provided context optionally controls which HTTP client
// is returned. See the [oauth2.HTTPClient] variable.
//
// The returned [http.Client] and its Transport should not be modified.
func (c *Config) HTTPClient(ctx context.Context) *http.Client {
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

	return oauth2.ReuseTokenSourceWithExpiry(nil, source, earlyExpiry)
}

type tokenSource struct {
	ctx          context.Context
	conf         *Config
	accessToken  string
	refreshToken string
}

// Token implements [oauth2.TokenSource].
func (t *tokenSource) Token() (*oauth2.Token, error) {
	client := oauth2.NewClient(t.ctx, nil)
	defer client.CloseIdleConnections()

	if t.accessToken != "" &&
		t.refreshToken != "" &&
		getExpirationTime(t.accessToken).Round(0).Add(-10*time.Second).After(time.Now()) {
		token, err := t.requestAccessToken(
			client,
			"GET",
			fmt.Sprintf(apiURLRefreshToken, t.refreshToken),
			func(req *http.Request) { req.Header.Set("Authorization", "Bearer "+t.accessToken) },
			nil,
			refreshTokenErrorPrefix)
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
		return nil, fmt.Errorf(retrieveTokenErrorPrefix+"failed to marshal request body: %w", err)
	}

	token, err := t.requestAccessToken(
		client,
		"POST",
		apiURLLogin,
		func(req *http.Request) { req.Header.Set("Content-Type", "application/json") },
		bytes.NewBuffer(reqBody),
		retrieveTokenErrorPrefix)

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *tokenSource) requestAccessToken(
	client *http.Client,
	method string,
	url string,
	requestProcessor func(*http.Request),
	body io.Reader,
	errorPrefix string) (*oauth2.Token, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if requestProcessor != nil {
		requestProcessor(req)
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("%sfailed to make refresh token request: %w", errorPrefix, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var result Result[any]
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && !result.Ok {
			return nil, fmt.Errorf("%sAPI error: %s", errorPrefix, result.Error)
		}

		return nil, fmt.Errorf("%sreceived non-2xx status code: %s (%d)", errorPrefix, resp.Status, resp.StatusCode)
	}

	var result Result[struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%sfailed to decode response: %w", errorPrefix, err)
	}
	if !result.Ok {
		return nil, fmt.Errorf("%sAPI error: %s", errorPrefix, result.Error)
	}

	t.accessToken = result.Data.AccessToken
	t.refreshToken = result.Data.RefreshToken

	return &oauth2.Token{
		AccessToken: result.Data.AccessToken,
		TokenType:   "Bearer",
		Expiry:      getExpirationTime(result.Data.AccessToken),
	}, nil
}

func getExpirationTime(accessToken string) time.Time {
	const tokenDelim = "."

	_, s, ok := strings.Cut(accessToken, tokenDelim)
	if !ok { // no period found
		log.Printf("%sinvalid access token format", errorPrefix)
		return past
	}

	payload, s, ok := strings.Cut(s, tokenDelim)
	if !ok { // only one period found
		log.Printf("%sinvalid access token format", errorPrefix)
		return past
	}

	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		log.Printf("%sfailed to decode access token payload: %v", errorPrefix, err)
		return past
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		log.Printf("%sfailed to unmarshal access token claims: %v", errorPrefix, err)
		return past
	}

	return time.Unix(claims.Exp, 0)
}
