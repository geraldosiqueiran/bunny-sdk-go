package internal_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal"
)

// TestNewRequest tests HTTP request creation with context
func TestNewRequest(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		method  string
		url     string
		body    io.Reader
		wantErr bool
	}{
		{
			name:    "GET request without body",
			method:  http.MethodGet,
			url:     "https://api.example.com/test",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "POST request with body",
			method:  http.MethodPost,
			url:     "https://api.example.com/test",
			body:    strings.NewReader(`{"key":"value"}`),
			wantErr: false,
		},
		{
			name:    "PUT request with body",
			method:  http.MethodPut,
			url:     "https://api.example.com/test/123",
			body:    bytes.NewBufferString("test data"),
			wantErr: false,
		},
		{
			name:    "DELETE request",
			method:  http.MethodDelete,
			url:     "https://api.example.com/test/123",
			body:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := internal.NewRequest(ctx, tt.method, tt.url, tt.body)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if req == nil {
					t.Fatal("NewRequest() returned nil request")
				}

				if req.Method != tt.method {
					t.Errorf("expected method %s, got %s", tt.method, req.Method)
				}

				if req.URL.String() != tt.url {
					t.Errorf("expected URL %s, got %s", tt.url, req.URL.String())
				}

				if req.Context() != ctx {
					t.Error("request context doesn't match provided context")
				}

				// Verify body is set correctly
				if tt.body != nil && req.Body == nil {
					t.Error("expected request body to be set, got nil")
				}
				if tt.body == nil && req.Body != nil {
					t.Error("expected request body to be nil, got non-nil")
				}
			}
		})
	}
}

// TestNewRequestWithCancelledContext tests request creation with cancelled context
func TestNewRequestWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req, err := internal.NewRequest(ctx, http.MethodGet, "https://api.example.com/test", nil)

	// Request creation should succeed even with cancelled context
	// The cancellation is checked when the request is executed
	if err != nil {
		t.Errorf("NewRequest() with cancelled context failed: %v", err)
	}
	if req == nil {
		t.Fatal("NewRequest() returned nil request")
	}
}

// TestNewRequestWithInvalidURL tests request creation with invalid URL
func TestNewRequestWithInvalidURL(t *testing.T) {
	ctx := context.Background()

	// Invalid URLs that should cause errors
	invalidURLs := []string{
		"://invalid",
		"ht!tp://invalid",
		string([]byte{0x7f}), // Invalid characters
	}

	for _, url := range invalidURLs {
		t.Run("invalid_url_"+url, func(t *testing.T) {
			_, err := internal.NewRequest(ctx, http.MethodGet, url, nil)
			if err == nil {
				t.Error("expected error for invalid URL, got nil")
			}
		})
	}
}

// TestDecodeResponse tests JSON response decoding
func TestDecodeResponse(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name       string
		body       string
		wantResult TestStruct
		wantErr    bool
	}{
		{
			name:       "valid JSON",
			body:       `{"name":"test","value":42}`,
			wantResult: TestStruct{Name: "test", Value: 42},
			wantErr:    false,
		},
		{
			name:       "empty object",
			body:       `{}`,
			wantResult: TestStruct{},
			wantErr:    false,
		},
		{
			name:    "invalid JSON",
			body:    `{"name":"test"`,
			wantErr: true,
		},
		{
			name:    "not JSON",
			body:    `plain text`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
				Header:     make(http.Header),
			}

			result, err := internal.DecodeResponse[TestStruct](resp)

			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Fatal("DecodeResponse() returned nil result")
				}
				if result.Name != tt.wantResult.Name {
					t.Errorf("expected name %s, got %s", tt.wantResult.Name, result.Name)
				}
				if result.Value != tt.wantResult.Value {
					t.Errorf("expected value %d, got %d", tt.wantResult.Value, result.Value)
				}
			}
		})
	}
}

// TestDecodeResponseWithDifferentTypes tests decoding different types
func TestDecodeResponseWithDifferentTypes(t *testing.T) {
	t.Run("slice type", func(t *testing.T) {
		body := `["item1","item2","item3"]`
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
		}

		result, err := internal.DecodeResponse[[]string](resp)
		if err != nil {
			t.Fatalf("DecodeResponse() error = %v", err)
		}
		if len(*result) != 3 {
			t.Errorf("expected 3 items, got %d", len(*result))
		}
	})

	t.Run("map type", func(t *testing.T) {
		body := `{"key1":"value1","key2":"value2"}`
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
		}

		result, err := internal.DecodeResponse[map[string]string](resp)
		if err != nil {
			t.Fatalf("DecodeResponse() error = %v", err)
		}
		if len(*result) != 2 {
			t.Errorf("expected 2 items, got %d", len(*result))
		}
	})

	t.Run("primitive type", func(t *testing.T) {
		body := `42`
		resp := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
		}

		result, err := internal.DecodeResponse[int](resp)
		if err != nil {
			t.Fatalf("DecodeResponse() error = %v", err)
		}
		if *result != 42 {
			t.Errorf("expected 42, got %d", *result)
		}
	})
}

// TestReadResponseBody tests reading response body
func TestReadResponseBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "normal body",
			body:    "response body content",
			wantErr: false,
		},
		{
			name:    "empty body",
			body:    "",
			wantErr: false,
		},
		{
			name:    "large body",
			body:    strings.Repeat("x", 10000),
			wantErr: false,
		},
		{
			name:    "JSON body",
			body:    `{"key":"value","number":123}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			body, err := internal.ReadResponseBody(resp)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadResponseBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(body) != tt.body {
					t.Errorf("expected body %q, got %q", tt.body, string(body))
				}
			}
		})
	}
}

// TestReadResponseBodyClosesBody tests that ReadResponseBody closes the body
func TestReadResponseBodyClosesBody(t *testing.T) {
	closeCalled := false
	body := &mockReadCloser{
		reader: strings.NewReader("test content"),
		onClose: func() error {
			closeCalled = true
			return nil
		},
	}

	resp := &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	_, err := internal.ReadResponseBody(resp)
	if err != nil {
		t.Fatalf("ReadResponseBody() error = %v", err)
	}

	if !closeCalled {
		t.Error("ReadResponseBody() did not close the body")
	}
}

// TestParseErrorResponse tests error response parsing
func TestParseErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		wantMsg  string
		wantKey  string
		wantField string
		wantNil  bool
	}{
		{
			name:    "complete error response",
			body:    []byte(`{"Message":"Invalid input","ErrorKey":"INVALID_INPUT","Field":"email"}`),
			wantMsg: "Invalid input",
			wantKey: "INVALID_INPUT",
			wantField: "email",
		},
		{
			name:    "error with message only",
			body:    []byte(`{"Message":"Not found"}`),
			wantMsg: "Not found",
		},
		{
			name:    "error with key only",
			body:    []byte(`{"Message":"Error occurred","ErrorKey":"ERROR"}`),
			wantMsg: "Error occurred",
			wantKey: "ERROR",
		},
		{
			name:    "invalid JSON",
			body:    []byte(`not valid json`),
			wantNil: true,
		},
		{
			name:    "empty body",
			body:    []byte(``),
			wantNil: true,
		},
		{
			name:    "empty object",
			body:    []byte(`{}`),
			wantMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ParseErrorResponse(tt.body)

			if tt.wantNil {
				if result != nil {
					t.Errorf("expected nil result, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("ParseErrorResponse() returned nil")
			}

			if result.Message != tt.wantMsg {
				t.Errorf("expected message %q, got %q", tt.wantMsg, result.Message)
			}
			if result.ErrorKey != tt.wantKey {
				t.Errorf("expected error key %q, got %q", tt.wantKey, result.ErrorKey)
			}
			if result.Field != tt.wantField {
				t.Errorf("expected field %q, got %q", tt.wantField, result.Field)
			}
		})
	}
}

// TestParseErrorResponseWithVariousFormats tests different error formats
func TestParseErrorResponseWithVariousFormats(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "lowercase fields",
			body: `{"message":"error","errorKey":"KEY"}`,
		},
		{
			name: "extra fields",
			body: `{"Message":"error","ErrorKey":"KEY","ExtraField":"value"}`,
		},
		{
			name: "nested object (invalid for our parser)",
			body: `{"Message":"error","Details":{"Code":123}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ParseErrorResponse([]byte(tt.body))
			// Just verify it doesn't panic
			_ = result
		})
	}
}

// mockReadCloser is a helper for testing body closing
type mockReadCloser struct {
	reader  io.Reader
	onClose func() error
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *mockReadCloser) Close() error {
	if m.onClose != nil {
		return m.onClose()
	}
	return nil
}

// TestErrorResponseStruct tests the ErrorResponse struct directly
func TestErrorResponseStruct(t *testing.T) {
	errResp := &internal.ErrorResponse{
		Message:  "Test error",
		ErrorKey: "TEST_ERROR",
		Field:    "test_field",
	}

	// Marshal and unmarshal to verify JSON tags
	data, err := json.Marshal(errResp)
	if err != nil {
		t.Fatalf("failed to marshal ErrorResponse: %v", err)
	}

	var decoded internal.ErrorResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal ErrorResponse: %v", err)
	}

	if decoded.Message != errResp.Message {
		t.Errorf("expected message %q, got %q", errResp.Message, decoded.Message)
	}
	if decoded.ErrorKey != errResp.ErrorKey {
		t.Errorf("expected error key %q, got %q", errResp.ErrorKey, decoded.ErrorKey)
	}
	if decoded.Field != errResp.Field {
		t.Errorf("expected field %q, got %q", errResp.Field, decoded.Field)
	}
}
