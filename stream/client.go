package stream

import (
	"context"
	"io"
	"net/http"
)

const (
	defaultBaseURL   = "https://video.bunnycdn.com"
	defaultUserAgent = "bunny-sdk-go/1.0"
)

// Client is a client for the Bunny.net Stream API.
type Client struct {
	apiKey     string
	httpClient HTTPClient
	userAgent  string
	baseURL    string
}

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
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

// NewClient creates a new Stream API client.
// Use the Stream Library API Key (found in the Stream library's API settings).
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

// Videos returns a VideoService for managing videos in the specified library.
func (c *Client) Videos(libraryID int64) VideoService {
	return newVideoService(&clientAdapter{c}, libraryID)
}

// Libraries returns a LibraryService for managing video libraries.
func (c *Client) Libraries() LibraryService {
	return newLibraryService(&clientAdapter{c})
}

// Collections returns a CollectionService for managing collections in the specified library.
func (c *Client) Collections(libraryID int64) CollectionService {
	return newCollectionService(&clientAdapter{c}, libraryID)
}

// clientAdapter adapts Client to the httpClient interface used by services.
type clientAdapter struct {
	client *Client
}

func (a *clientAdapter) do(ctx context.Context, method, path string, body any, result any) error {
	return doRequestJSON(ctx, a.client.httpClient, a.client.baseURL, a.client.apiKey, a.client.userAgent, method, path, body, result)
}

func (a *clientAdapter) doRaw(ctx context.Context, method, path string, body io.Reader, contentType string) error {
	return doRequestRaw(ctx, a.client.httpClient, a.client.baseURL, a.client.apiKey, a.client.userAgent, method, path, body, contentType)
}
