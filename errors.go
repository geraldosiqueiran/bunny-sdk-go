package bunny

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Bunny.net API.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"Message"`
	ErrorKey   string `json:"ErrorKey,omitempty"`
	Field      string `json:"Field,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("bunny: %s (status: %d, field: %s)", e.Message, e.StatusCode, e.Field)
	}
	if e.ErrorKey != "" {
		return fmt.Sprintf("bunny: %s (status: %d, key: %s)", e.Message, e.StatusCode, e.ErrorKey)
	}
	return fmt.Sprintf("bunny: %s (status: %d)", e.Message, e.StatusCode)
}

// IsAuthError returns true if the error is an authentication error (401 or 403).
func (e *APIError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

// IsNotFound returns true if the error is a not found error (404).
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsRateLimited returns true if the error is a rate limit error (429).
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsRetryable returns true if the error is retryable (5xx or 429).
func (e *APIError) IsRetryable() bool {
	return e.StatusCode >= 500 || e.StatusCode == http.StatusTooManyRequests
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	*APIError
}

// AuthError represents a 401 Unauthorized or 403 Forbidden error.
type AuthError struct {
	*APIError
}

// RateLimitError represents a 429 Too Many Requests error.
type RateLimitError struct {
	*APIError
}

// newError creates the appropriate error type based on status code.
func newError(statusCode int, message, errorKey, field string) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		Message:    message,
		ErrorKey:   errorKey,
		Field:      field,
	}

	switch statusCode {
	case http.StatusNotFound:
		return &NotFoundError{APIError: apiErr}
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthError{APIError: apiErr}
	case http.StatusTooManyRequests:
		return &RateLimitError{APIError: apiErr}
	default:
		return apiErr
	}
}
