package bunny

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client for the client.
func WithHTTPClient(hc HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithUserAgent sets a custom user agent string for API requests.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithStreamBaseURL sets a custom base URL for the Stream API.
func WithStreamBaseURL(url string) Option {
	return func(c *Client) {
		c.streamBaseURL = url
	}
}

// WithStorageBaseURL sets a custom base URL for the Storage management API.
func WithStorageBaseURL(url string) Option {
	return func(c *Client) {
		c.storageBaseURL = url
	}
}
