package stream

import (
	"context"
	"io"
	"net/http"
)

const (
	// Base API URL for library management operations (list, create, delete libraries)
	defaultBaseAPIURL = "https://api.bunny.net"
	// Stream API URL for video/collection operations within a library
	defaultStreamAPIURL = "https://video.bunnycdn.com"
	defaultUserAgent    = "bunny-sdk-go/1.0"
)

// Client is a client for the Bunny.net Stream API.
type Client struct {
	apiKey       string
	httpClient   HTTPClient
	userAgent    string
	baseAPIURL   string // for library management (api.bunny.net)
	streamAPIURL string // for video/collection operations (video.bunnycdn.com)
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

// WithBaseAPIURL sets a custom base API URL for library management.
func WithBaseAPIURL(url string) Option {
	return func(c *Client) {
		c.baseAPIURL = url
	}
}

// WithStreamAPIURL sets a custom stream API URL for video/collection operations.
func WithStreamAPIURL(url string) Option {
	return func(c *Client) {
		c.streamAPIURL = url
	}
}

// NewClient creates a new Stream API client.
// Use the Global API Key (found in Account Settings > API) for library management.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:       apiKey,
		httpClient:   http.DefaultClient,
		userAgent:    defaultUserAgent,
		baseAPIURL:   defaultBaseAPIURL,
		streamAPIURL: defaultStreamAPIURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Videos returns a VideoService for managing videos in the specified library.
func (c *Client) Videos(libraryID int64) VideoService {
	return newVideoService(&streamAdapter{c}, libraryID)
}

// Libraries returns a LibraryService for managing video libraries.
func (c *Client) Libraries() LibraryService {
	return newLibraryService(&baseAdapter{c})
}

// Collections returns a CollectionService for managing collections in the specified library.
func (c *Client) Collections(libraryID int64) CollectionService {
	return newCollectionService(&streamAdapter{c}, libraryID)
}

// OEmbed returns an OEmbedService for video embedding.
func (c *Client) OEmbed() OEmbedService {
	return newOEmbedService(&streamAdapter{c})
}

// baseAdapter uses api.bunny.net for library management operations.
type baseAdapter struct {
	client *Client
}

func (a *baseAdapter) do(ctx context.Context, method, path string, body any, result any) error {
	return doRequestJSON(ctx, a.client.httpClient, a.client.baseAPIURL, a.client.apiKey, a.client.userAgent, method, path, body, result)
}

// streamAdapter uses video.bunnycdn.com for video/collection operations.
type streamAdapter struct {
	client *Client
}

func (a *streamAdapter) do(ctx context.Context, method, path string, body any, result any) error {
	return doRequestJSON(ctx, a.client.httpClient, a.client.streamAPIURL, a.client.apiKey, a.client.userAgent, method, path, body, result)
}

func (a *streamAdapter) doRaw(ctx context.Context, method, path string, body io.Reader, contentType string) error {
	return doRequestRaw(ctx, a.client.httpClient, a.client.streamAPIURL, a.client.apiKey, a.client.userAgent, method, path, body, contentType)
}
