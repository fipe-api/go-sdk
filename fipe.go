// Package fipe is a client for the FIPE API (https://fipe.api.br),
// which provides average vehicle prices from Brazil's Tabela FIPE.
package fipe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DefaultBaseURL is the production endpoint of the FIPE API.
const DefaultBaseURL = "https://fipe.api.br/api/v2"

// Client is a FIPE API client. Create one with New; the zero value is not usable.
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the *http.Client used for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithSubscriptionToken sets the X-Subscription-Token header on every request.
// The free tier works without a token but has stricter rate limits.
func WithSubscriptionToken(token string) Option {
	return func(c *Client) { c.token = token }
}

// WithBaseURL overrides the API base URL (e.g. for testing).
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") }
}

// New returns a Client configured with the given options.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// RequestOption configures a single API request.
type RequestOption func(url.Values)

// WithReference queries a specific FIPE monthly reference table instead of
// the current one. Codes come from Client.References.
func WithReference(code int) RequestOption {
	return func(q url.Values) { q.Set("reference", strconv.Itoa(code)) }
}

func (c *Client) get(ctx context.Context, path string, out any, opts ...RequestOption) error {
	q := url.Values{}
	for _, opt := range opts {
		opt(q)
	}
	u := c.baseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("fipe: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("X-Subscription-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fipe: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return &APIError{StatusCode: resp.StatusCode, Body: string(bytes.TrimSpace(body))}
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("fipe: decode response: %w", err)
	}
	return nil
}
