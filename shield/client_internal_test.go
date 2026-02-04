package shield

import (
	"context"
	"net/http"
	"testing"
)

// Test marshal error via internal test
func TestDoRequest_MarshalError(t *testing.T) {
	c := NewClient("key")
	// Channel type cannot be marshaled
	badBody := make(chan int)
	err := c.doRequest(context.Background(), http.MethodPost, "/test", badBody, nil)
	if err == nil {
		t.Fatal("expected marshal error")
	}
}

// Test request creation error with nil context
func TestDoRequest_RequestError(t *testing.T) {
	c := NewClient("key")
	//nolint:staticcheck // intentionally passing nil context
	err := c.doRequest(nil, http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected request error")
	}
}
