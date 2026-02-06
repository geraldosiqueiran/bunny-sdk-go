package containers

import (
	"context"
	"net/http"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")

	if client.apiKey != "test-api-key" {
		t.Errorf("expected apiKey to be 'test-api-key', got '%s'", client.apiKey)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL to be '%s', got '%s'", defaultBaseURL, client.baseURL)
	}
	if client.userAgent != defaultUserAgent {
		t.Errorf("expected userAgent to be '%s', got '%s'", defaultUserAgent, client.userAgent)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	mockHTTP := &testutil.MockHTTPClient{}
	client := NewClient("test-key",
		WithHTTPClient(mockHTTP),
		WithBaseURL("https://custom.api.com"),
		WithUserAgent("custom-agent/1.0"),
	)

	if client.baseURL != "https://custom.api.com" {
		t.Errorf("expected custom baseURL, got '%s'", client.baseURL)
	}
	if client.userAgent != "custom-agent/1.0" {
		t.Errorf("expected custom userAgent, got '%s'", client.userAgent)
	}
}

func TestClientSetsAuthHeader(t *testing.T) {
	var capturedReq *http.Request
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedReq = req
			return testutil.NewMockResponse(200, `{"items":[],"meta":{"totalItems":0}}`), nil
		},
	}

	client := NewClient("my-api-key", WithHTTPClient(mock))
	_, _ = client.Applications().List(context.Background(), nil)

	if capturedReq == nil {
		t.Fatal("expected request to be captured")
	}
	if capturedReq.Header.Get("AccessKey") != "my-api-key" {
		t.Errorf("expected AccessKey header, got '%s'", capturedReq.Header.Get("AccessKey"))
	}
}

func TestAPIError(t *testing.T) {
	err := newAPIError(404, "not found", "NOT_FOUND", "appId")

	if err.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", err.StatusCode)
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestAPIErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		contains []string
	}{
		{
			name:     "basic error",
			err:      newAPIError(400, "bad request", "", ""),
			contains: []string{"status 400", "bad request"},
		},
		{
			name:     "error with key",
			err:      newAPIError(404, "not found", "NOT_FOUND", ""),
			contains: []string{"status 404", "not found", "error_key: NOT_FOUND"},
		},
		{
			name:     "error with field",
			err:      newAPIError(422, "invalid", "", "name"),
			contains: []string{"status 422", "invalid", "field: name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			for _, s := range tt.contains {
				if !contains(msg, s) {
					t.Errorf("expected error message to contain '%s', got '%s'", s, msg)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
