package containers

import (
	"context"
	"fmt"
	"net/http"
)

// LogForwardingService provides methods for managing log forwarding configurations.
type LogForwardingService interface {
	// List returns all log forwarding configurations.
	List(ctx context.Context) (*LogForwardingListResponse, error)

	// Get returns the log forwarding configuration for a specific application.
	Get(ctx context.Context, appID string) (*LogForwardingConfig, error)

	// Create creates a new log forwarding configuration.
	Create(ctx context.Context, req *CreateLogForwardingRequest) (*LogForwardingConfig, error)

	// Update updates an existing log forwarding configuration.
	Update(ctx context.Context, appID string, req *UpdateLogForwardingRequest) (*LogForwardingConfig, error)

	// Delete removes a log forwarding configuration.
	Delete(ctx context.Context, appID string) error
}

type logForwardingService struct {
	client httpClient
}

func newLogForwardingService(client httpClient) LogForwardingService {
	return &logForwardingService{client: client}
}

// List returns all log forwarding configurations.
func (s *logForwardingService) List(ctx context.Context) (*LogForwardingListResponse, error) {
	var resp LogForwardingListResponse
	if err := s.client.do(ctx, http.MethodGet, "/log/forwarding", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns the log forwarding configuration for a specific application.
func (s *logForwardingService) Get(ctx context.Context, appID string) (*LogForwardingConfig, error) {
	path := fmt.Sprintf("/log/forwarding/%s", appID)

	var config LogForwardingConfig
	if err := s.client.do(ctx, http.MethodGet, path, nil, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Create creates a new log forwarding configuration.
func (s *logForwardingService) Create(ctx context.Context, req *CreateLogForwardingRequest) (*LogForwardingConfig, error) {
	var config LogForwardingConfig
	if err := s.client.do(ctx, http.MethodPost, "/log/forwarding", req, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Update updates an existing log forwarding configuration.
func (s *logForwardingService) Update(ctx context.Context, appID string, req *UpdateLogForwardingRequest) (*LogForwardingConfig, error) {
	path := fmt.Sprintf("/log/forwarding/%s", appID)

	var config LogForwardingConfig
	if err := s.client.do(ctx, http.MethodPut, path, req, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Delete removes a log forwarding configuration.
func (s *logForwardingService) Delete(ctx context.Context, appID string) error {
	path := fmt.Sprintf("/log/forwarding/%s", appID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}
