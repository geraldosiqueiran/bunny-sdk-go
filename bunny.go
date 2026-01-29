// Package bunny provides a Go client for the Bunny.net API.
//
// The SDK supports both Stream API (video management) and Storage API (zone management and file operations).
//
// Example usage:
//
//	// Create a client for Stream API
//	streamClient := bunny.NewStreamClient("your-stream-api-key")
//
//	// Create a client for Storage management API
//	client := bunny.NewClient("your-global-api-key")
//
//	// Create a client for Storage file operations
//	fileClient := bunny.NewStorageClient("my-zone", "zone-password", "storage")
package bunny

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

	"github.com/geraldo/bunny-sdk-go/internal"
)

const (
	defaultUserAgent      = "bunny-sdk-go/1.0"
	defaultStreamBaseURL  = "https://video.bunnycdn.com"
	defaultStorageBaseURL = "https://api.bunny.net"
)

// Client is the main Bunny.net API client for management operations.
type Client struct {
	apiKey         string
	httpClient     HTTPClient
	userAgent      string
	streamBaseURL  string
	storageBaseURL string
}

// NewClient creates a new Bunny.net management API client.
// Use the Global API Key for storage zone management operations.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:         apiKey,
		httpClient:     http.DefaultClient,
		userAgent:      defaultUserAgent,
		streamBaseURL:  defaultStreamBaseURL,
		storageBaseURL: defaultStorageBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// StreamClient is a client for the Bunny.net Stream API.
type StreamClient struct {
	apiKey     string
	httpClient HTTPClient
	userAgent  string
	baseURL    string
}

// NewStreamClient creates a new Bunny.net Stream API client.
// Use the Stream Library API Key (found in the Stream library settings).
func NewStreamClient(apiKey string, opts ...Option) *StreamClient {
	// Create a temporary Client to apply options
	c := &Client{
		httpClient:    http.DefaultClient,
		userAgent:     defaultUserAgent,
		streamBaseURL: defaultStreamBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}

	return &StreamClient{
		apiKey:     apiKey,
		httpClient: c.httpClient,
		userAgent:  c.userAgent,
		baseURL:    c.streamBaseURL,
	}
}

// StorageClient is a client for Bunny.net Edge Storage file operations.
type StorageClient struct {
	zoneName   string
	accessKey  string
	httpClient HTTPClient
	userAgent  string
	baseURL    string
}

// StorageRegion represents a storage region endpoint.
type StorageRegion string

const (
	RegionFalkenstein StorageRegion = "storage"  // Falkenstein, Germany (default)
	RegionNewYork     StorageRegion = "ny"       // New York
	RegionLosAngeles  StorageRegion = "la"       // Los Angeles
	RegionSingapore   StorageRegion = "sg"       // Singapore
	RegionSydney      StorageRegion = "syd"      // Sydney
	RegionStockholm   StorageRegion = "se"       // Stockholm
	RegionSaoPaulo    StorageRegion = "br"       // Sao Paulo
	RegionJohannesburg StorageRegion = "jh"      // Johannesburg
	RegionLondon      StorageRegion = "uk"       // London
)

// NewStorageClient creates a new Bunny.net Edge Storage client for file operations.
// Use the Storage Zone Password (found in FTP & API Access settings).
func NewStorageClient(zoneName, accessKey string, region StorageRegion, opts ...Option) *StorageClient {
	baseURL := fmt.Sprintf("https://%s.bunnycdn.com", region)

	// Create a temporary Client to apply options
	c := &Client{
		httpClient: http.DefaultClient,
		userAgent:  defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}

	return &StorageClient{
		zoneName:   zoneName,
		accessKey:  accessKey,
		httpClient: c.httpClient,
		userAgent:  c.userAgent,
		baseURL:    baseURL,
	}
}

// doRequest performs an HTTP request and handles errors.
func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	return doRequestWithClient(ctx, c.httpClient, c.storageBaseURL, c.apiKey, c.userAgent, method, path, body, result)
}

// doRequest performs an HTTP request for the Stream client.
func (c *StreamClient) doRequest(ctx context.Context, method, path string, body any, result any) error {
	return doRequestWithClient(ctx, c.httpClient, c.baseURL, c.apiKey, c.userAgent, method, path, body, result)
}

// doRequestWithClient is a shared helper for performing HTTP requests.
func doRequestWithClient(ctx context.Context, httpClient HTTPClient, baseURL, apiKey, userAgent, method, path string, body any, result any) error {
	fullURL := baseURL + path

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

	setAuthHeader(req, apiKey)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
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

// handleErrorResponse parses and returns an appropriate error from an HTTP response.
func handleErrorResponse(resp *http.Response) error {
	body, err := internal.ReadResponseBody(resp)
	if err != nil {
		return newError(resp.StatusCode, "failed to read error response", "", "")
	}

	errResp := internal.ParseErrorResponse(body)
	if errResp != nil && errResp.Message != "" {
		return newError(resp.StatusCode, errResp.Message, errResp.ErrorKey, errResp.Field)
	}

	return newError(resp.StatusCode, string(body), "", "")
}

// buildListURL builds a URL with pagination and search query parameters.
func buildListURL(basePath string, opts *ListOptions) string {
	if opts == nil {
		return basePath
	}

	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.ItemsPerPage > 0 {
		params.Set("itemsPerPage", strconv.Itoa(opts.ItemsPerPage))
	}
	if opts.Search != "" {
		params.Set("search", opts.Search)
	}
	if opts.OrderBy != "" {
		params.Set("orderBy", opts.OrderBy)
	}

	if len(params) == 0 {
		return basePath
	}

	if strings.Contains(basePath, "?") {
		return basePath + "&" + params.Encode()
	}
	return basePath + "?" + params.Encode()
}
