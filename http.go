// Package bunny provides a Go client for the Bunny.net API.
package bunny

import "net/http"

// HTTPClient is an interface for making HTTP requests.
// This allows for easy mocking and custom HTTP client implementations.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
