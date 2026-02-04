package shield

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

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a client for the Bunny.net Shield/WAF API.
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

// NewClient creates a new Shield/WAF API client.
// Use the Global API Key (found in Account Settings > API).
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

// httpClient is the internal interface for making API requests.
type httpClient interface {
	do(ctx context.Context, method, path string, body any, result any) error
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

// Zones returns a ZoneService for managing Shield zones.
func (c *Client) Zones() ZoneService {
	return newZoneService(&clientAdapter{c})
}

// WAF returns a WAFService for managing WAF rules and configuration.
func (c *Client) WAF() WAFService {
	return newWAFService(&clientAdapter{c})
}

// AccessLists returns an AccessListService for managing access lists in the specified zone.
func (c *Client) AccessLists(zoneID string) AccessListService {
	return newAccessListService(&clientAdapter{c}, zoneID)
}

// RateLimits returns a RateLimitService for managing rate limit rules.
func (c *Client) RateLimits() RateLimitService {
	return newRateLimitService(&clientAdapter{c})
}

// BotDetection returns a BotDetectionService for the specified zone.
func (c *Client) BotDetection(zoneID string) BotDetectionService {
	return newBotDetectionService(&clientAdapter{c}, zoneID)
}

// UploadScanning returns an UploadScanningService for the specified zone.
func (c *Client) UploadScanning(zoneID string) UploadScanningService {
	return newUploadScanningService(&clientAdapter{c}, zoneID)
}

// Metrics returns a MetricsService for accessing Shield metrics.
func (c *Client) Metrics() MetricsService {
	return newMetricsService(&clientAdapter{c})
}

// EventLogs returns an EventLogsService for accessing security event logs.
func (c *Client) EventLogs() EventLogsService {
	return newEventLogsService(&clientAdapter{c})
}

// DDoS returns a DDoSService for DDoS protection information.
func (c *Client) DDoS() DDoSService {
	return newDDoSService(&clientAdapter{c})
}

// Promo returns a PromoService for promotional information.
func (c *Client) Promo() PromoService {
	return newPromoService(&clientAdapter{c})
}
