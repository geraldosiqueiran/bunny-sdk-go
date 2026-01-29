package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/geraldo/bunny-sdk-go/internal"
)

const (
	defaultBaseURL   = "https://api.bunny.net"
	defaultUserAgent = "bunny-sdk-go/1.0"
)

// Client is a client for the Bunny.net Storage Zone Management API.
// Use this client for managing storage zones (create, list, update, delete).
// For file operations, use NewFileService instead.
type Client struct {
	apiKey     string
	httpClient HTTPClient
	userAgent  string
	baseURL    string
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithUserAgent sets a custom user agent string.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// NewClient creates a new Storage Zone Management API client.
// Use the Global API Key (found in account settings).
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		userAgent:  defaultUserAgent,
		baseURL:    defaultBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Zones returns a ZoneService for managing storage zones.
func (c *Client) Zones() ZoneService {
	return newZoneService(&clientAdapter{c})
}

// clientAdapter adapts Client to the httpClient interface.
type clientAdapter struct {
	client *Client
}

func (a *clientAdapter) do(ctx context.Context, method, path string, body any, result any) error {
	return a.client.doRequest(ctx, method, path, body, result)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := internal.NewRequest(ctx, method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("AccessKey", c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return handleErrorResponse(resp)
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	} else {
		resp.Body.Close()
	}

	return nil
}

func handleErrorResponse(resp *http.Response) error {
	body, err := internal.ReadResponseBody(resp)
	if err != nil {
		return newAPIError(resp.StatusCode, "failed to read error response", "", "")
	}

	errResp := internal.ParseErrorResponse(body)
	if errResp != nil && errResp.Message != "" {
		return newAPIError(resp.StatusCode, errResp.Message, errResp.ErrorKey, errResp.Field)
	}

	return newAPIError(resp.StatusCode, string(body), "", "")
}
