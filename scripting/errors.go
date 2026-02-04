package scripting

import "fmt"

// APIError represents an error response from the Bunny.net Edge Scripting API.
type APIError struct {
	StatusCode int
	Message    string
	ErrorKey   string
	Field      string
}

// Error returns the error message.
func (e *APIError) Error() string {
	msg := fmt.Sprintf("bunny scripting api: status %d", e.StatusCode)
	if e.Message != "" {
		msg += ": " + e.Message
	}
	if e.ErrorKey != "" {
		msg += " (error_key: " + e.ErrorKey + ")"
	}
	if e.Field != "" {
		msg += " (field: " + e.Field + ")"
	}
	return msg
}

func newAPIError(statusCode int, message, errorKey, field string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		ErrorKey:   errorKey,
		Field:      field,
	}
}
