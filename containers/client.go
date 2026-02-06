package containers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/geraldo/bunny-sdk-go/internal"
)

const (
	defaultBaseURL   = "https://api.bunny.net/mc"
	defaultUserAgent = "bunny-sdk-go/1.0"
)

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a client for the Bunny.net Magic Containers API.
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
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets a custom user agent string.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// NewClient creates a new Magic Containers API client.
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

// Helper functions for building query strings

func buildListPath(basePath string, opts *ListOptions) string {
	params := url.Values{}
	if opts.NextCursor != "" {
		params.Set("nextCursor", opts.NextCursor)
	}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if len(params) == 0 {
		return basePath
	}
	return basePath + "?" + params.Encode()
}

func buildStatisticsPath(basePath string, opts *StatisticsOptions) string {
	params := url.Values{}
	if opts.FromDate != "" {
		params.Set("fromDate", opts.FromDate)
	}
	if opts.ToDate != "" {
		params.Set("toDate", opts.ToDate)
	}
	if opts.Granularity != "" {
		params.Set("granularity", string(opts.Granularity))
	}
	if len(params) == 0 {
		return basePath
	}
	return basePath + "?" + params.Encode()
}

// =============================================================================
// Service Factory Methods
// =============================================================================

// Applications returns an ApplicationService for managing applications.
func (c *Client) Applications() ApplicationService {
	return newApplicationService(&clientAdapter{c})
}

// Registries returns a RegistryService for managing container registries.
func (c *Client) Registries() RegistryService {
	return newRegistryService(&clientAdapter{c})
}

// ContainerTemplates returns a ContainerTemplateService for the specified application.
func (c *Client) ContainerTemplates(appID string) ContainerTemplateService {
	return newContainerTemplateService(&clientAdapter{c}, appID)
}

// Endpoints returns an EndpointService for the specified application.
func (c *Client) Endpoints(appID string) EndpointService {
	return newEndpointService(&clientAdapter{c}, appID)
}

// Autoscaling returns an AutoscalingService for the specified application.
func (c *Client) Autoscaling(appID string) AutoscalingService {
	return newAutoscalingService(&clientAdapter{c}, appID)
}

// Regions returns a RegionService for managing regions.
func (c *Client) Regions() RegionService {
	return newRegionService(&clientAdapter{c})
}

// RegionSettings returns a RegionSettingsService for the specified application.
func (c *Client) RegionSettings(appID string) RegionSettingsService {
	return newRegionSettingsService(&clientAdapter{c}, appID)
}

// Limits returns a LimitsService for retrieving user limits.
func (c *Client) Limits() LimitsService {
	return newLimitsService(&clientAdapter{c})
}

// Nodes returns a NodeService for listing nodes.
func (c *Client) Nodes() NodeService {
	return newNodeService(&clientAdapter{c})
}

// Pods returns a PodService for the specified application.
func (c *Client) Pods(appID string) PodService {
	return newPodService(&clientAdapter{c}, appID)
}

// Volumes returns a VolumeService for the specified application.
func (c *Client) Volumes(appID string) VolumeService {
	return newVolumeService(&clientAdapter{c}, appID)
}

// LogForwarding returns a LogForwardingService for managing log forwarding.
func (c *Client) LogForwarding() LogForwardingService {
	return newLogForwardingService(&clientAdapter{c})
}
