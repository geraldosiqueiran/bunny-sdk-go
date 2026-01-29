// Package internal provides internal utilities for the Bunny SDK.
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NewRequest creates a new HTTP request with context.
func NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return req, nil
}

// DecodeResponse decodes a JSON response body into the given type.
func DecodeResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result, nil
}

// ReadResponseBody reads the entire response body and closes it.
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ErrorResponse represents an error response from the Bunny API.
type ErrorResponse struct {
	Message  string `json:"Message"`
	ErrorKey string `json:"ErrorKey,omitempty"`
	Field    string `json:"Field,omitempty"`
}

// ParseErrorResponse attempts to parse an error response from the body.
func ParseErrorResponse(body []byte) *ErrorResponse {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil
	}
	return &errResp
}
