package bunny_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/geraldo/bunny-sdk-go"
	"github.com/geraldo/bunny-sdk-go/internal/testutil"
	"github.com/geraldo/bunny-sdk-go/stream"
)

// TestAPIError tests basic APIError functionality
func TestAPIError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		errorKey       string
		field          string
		expectedString string
	}{
		{
			name:           "basic error",
			statusCode:     400,
			message:        "Bad request",
			expectedString: "bunny: Bad request (status: 400)",
		},
		{
			name:           "error with key",
			statusCode:     400,
			message:        "Invalid input",
			errorKey:       "INVALID_INPUT",
			expectedString: "bunny: Invalid input (status: 400, key: INVALID_INPUT)",
		},
		{
			name:           "error with field",
			statusCode:     400,
			message:        "Field validation failed",
			field:          "email",
			expectedString: "bunny: Field validation failed (status: 400, field: email)",
		},
		{
			name:           "error with both key and field",
			statusCode:     422,
			message:        "Validation error",
			errorKey:       "VALIDATION_ERROR",
			field:          "username",
			expectedString: "bunny: Validation error (status: 422, field: username)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &bunny.APIError{
				StatusCode: tt.statusCode,
				Message:    tt.message,
				ErrorKey:   tt.errorKey,
				Field:      tt.field,
			}

			if err.Error() != tt.expectedString {
				t.Errorf("Error() = %q, want %q", err.Error(), tt.expectedString)
			}
		})
	}
}

// TestIsAuthError tests authentication error detection
func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		isAuth     bool
	}{
		{"401 Unauthorized", http.StatusUnauthorized, true},
		{"403 Forbidden", http.StatusForbidden, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"404 Not Found", http.StatusNotFound, false},
		{"500 Internal Error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &bunny.APIError{
				StatusCode: tt.statusCode,
				Message:    "test error",
			}

			if got := err.IsAuthError(); got != tt.isAuth {
				t.Errorf("IsAuthError() = %v, want %v", got, tt.isAuth)
			}
		})
	}
}

// TestIsNotFound tests not found error detection
func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		isNotFound bool
	}{
		{"404 Not Found", http.StatusNotFound, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"401 Unauthorized", http.StatusUnauthorized, false},
		{"500 Internal Error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &bunny.APIError{
				StatusCode: tt.statusCode,
				Message:    "test error",
			}

			if got := err.IsNotFound(); got != tt.isNotFound {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.isNotFound)
			}
		})
	}
}

// TestIsRateLimited tests rate limit error detection
func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		isRateLimited bool
	}{
		{"429 Too Many Requests", http.StatusTooManyRequests, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"500 Internal Error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &bunny.APIError{
				StatusCode: tt.statusCode,
				Message:    "test error",
			}

			if got := err.IsRateLimited(); got != tt.isRateLimited {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.isRateLimited)
			}
		})
	}
}

// TestIsRetryable tests retryable error detection
func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		isRetryable bool
	}{
		{"429 Too Many Requests", http.StatusTooManyRequests, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
		{"502 Bad Gateway", http.StatusBadGateway, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"401 Unauthorized", http.StatusUnauthorized, false},
		{"404 Not Found", http.StatusNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &bunny.APIError{
				StatusCode: tt.statusCode,
				Message:    "test error",
			}

			if got := err.IsRetryable(); got != tt.isRetryable {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.isRetryable)
			}
		})
	}
}

// TestNotFoundError tests NotFoundError type
func TestNotFoundError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Message":"Resource not found"}`
			return testutil.NewMockResponse(404, body), nil
		},
	}

	// Use stream.NewClient which has Videos() method
	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check if it contains 404
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected error to contain 404, got: %s", err.Error())
	}
}

// TestAuthError tests AuthError type
func TestAuthError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"401 Unauthorized", http.StatusUnauthorized},
		{"403 Forbidden", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					body := `{"Message":"Authentication failed"}`
					return testutil.NewMockResponse(tt.statusCode, body), nil
				},
			}

			client := stream.NewClient("invalid-key", stream.WithHTTPClient(mock))
			_, err := client.Videos(123).List(context.Background(), nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// Verify error message contains status code
			errMsg := err.Error()
			if !strings.Contains(errMsg, "Authentication failed") {
				t.Errorf("expected error to contain 'Authentication failed', got: %s", errMsg)
			}
		})
	}
}

// TestRateLimitError tests RateLimitError type
func TestRateLimitError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Message":"Rate limit exceeded"}`
			return testutil.NewMockResponse(http.StatusTooManyRequests, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "429") {
		t.Errorf("expected error to contain 429, got: %s", err.Error())
	}
}

// TestServerError tests 5xx server errors
func TestServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"500 Internal Server Error", http.StatusInternalServerError, "Internal server error"},
		{"502 Bad Gateway", http.StatusBadGateway, "Bad gateway"},
		{"503 Service Unavailable", http.StatusServiceUnavailable, "Service unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					body := `{"Message":"` + tt.message + `"}`
					return testutil.NewMockResponse(tt.statusCode, body), nil
				},
			}

			client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
			_, err := client.Videos(123).List(context.Background(), nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			errMsg := err.Error()
			if !strings.Contains(errMsg, tt.message) {
				t.Errorf("expected error to contain %q, got: %s", tt.message, errMsg)
			}
		})
	}
}

// TestErrorWithoutMessage tests error handling when no message is provided
func TestErrorWithoutMessage(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Return error response without structured JSON
			return testutil.NewMockResponse(500, "Internal Server Error"), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to contain status code 500, got: %s", err.Error())
	}
}

// TestErrorResponseParsing tests various error response formats
func TestErrorResponseParsing(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		statusCode   int
		expectError  bool
	}{
		{
			name:         "valid error JSON",
			responseBody: `{"Message":"Invalid input","ErrorKey":"INVALID","Field":"name"}`,
			statusCode:   400,
			expectError:  true,
		},
		{
			name:         "error without ErrorKey",
			responseBody: `{"Message":"Not found"}`,
			statusCode:   404,
			expectError:  true,
		},
		{
			name:         "plain text error",
			responseBody: "Error occurred",
			statusCode:   500,
			expectError:  true,
		},
		{
			name:         "empty error body",
			responseBody: "",
			statusCode:   500,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return testutil.NewMockResponse(tt.statusCode, tt.responseBody), nil
				},
			}

			client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
			_, err := client.Videos(123).List(context.Background(), nil)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

// TestNetworkError tests network-level errors
func TestNetworkError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network connection failed")
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "network connection failed") && !strings.Contains(err.Error(), "request failed") {
		t.Errorf("expected network error, got: %s", err.Error())
	}
}

// TestErrorTypeAssertions tests type assertions for different error types
func TestErrorTypeAssertions(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		checkType  string
	}{
		{"NotFoundError", http.StatusNotFound, "NotFoundError"},
		{"AuthError 401", http.StatusUnauthorized, "AuthError"},
		{"AuthError 403", http.StatusForbidden, "AuthError"},
		{"RateLimitError", http.StatusTooManyRequests, "RateLimitError"},
		{"Generic APIError", http.StatusBadRequest, "APIError"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					body := `{"Message":"Error occurred"}`
					return testutil.NewMockResponse(tt.statusCode, body), nil
				},
			}

			client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
			_, err := client.Videos(123).List(context.Background(), nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// Verify error contains status code
			if !strings.Contains(err.Error(), "Error occurred") {
				t.Errorf("expected error message, got: %s", err.Error())
			}
		})
	}
}

// TestNewErrorFunction tests the newError function creates correct error types
func TestNewErrorFunction(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		errorKey   string
		field      string
		wantType   string
	}{
		{
			name:       "NotFoundError",
			statusCode: 404,
			message:    "resource not found",
			wantType:   "*bunny.NotFoundError",
		},
		{
			name:       "AuthError 401",
			statusCode: 401,
			message:    "unauthorized",
			wantType:   "*bunny.AuthError",
		},
		{
			name:       "AuthError 403",
			statusCode: 403,
			message:    "forbidden",
			wantType:   "*bunny.AuthError",
		},
		{
			name:       "RateLimitError",
			statusCode: 429,
			message:    "too many requests",
			wantType:   "*bunny.RateLimitError",
		},
		{
			name:       "Generic APIError",
			statusCode: 400,
			message:    "bad request",
			errorKey:   "BAD_REQUEST",
			field:      "input",
			wantType:   "*bunny.APIError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					body := `{"Message":"` + tt.message + `"`
					if tt.errorKey != "" {
						body += `,"ErrorKey":"` + tt.errorKey + `"`
					}
					if tt.field != "" {
						body += `,"Field":"` + tt.field + `"`
					}
					body += `}`
					return testutil.NewMockResponse(tt.statusCode, body), nil
				},
			}

			client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
			_, err := client.Videos(123).List(context.Background(), nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// Verify error message
			if !strings.Contains(err.Error(), tt.message) {
				t.Errorf("expected error to contain %q, got: %s", tt.message, err.Error())
			}
		})
	}
}
