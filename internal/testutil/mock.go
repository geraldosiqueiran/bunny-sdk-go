// Package testutil provides testing utilities for the Bunny SDK.
package testutil

import (
	"bytes"
	"io"
	"net/http"
)

// MockHTTPClient is a mock HTTP client for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do executes the mock function.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// NewMockResponse creates a mock HTTP response with the given status code and body.
func NewMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

// NewMockResponseWithHeaders creates a mock HTTP response with custom headers.
func NewMockResponseWithHeaders(statusCode int, body string, headers map[string]string) *http.Response {
	resp := NewMockResponse(statusCode, body)
	for k, v := range headers {
		resp.Header.Set(k, v)
	}
	return resp
}
