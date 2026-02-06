package containers

import (
	"context"
	"fmt"
	"net/http"
)

// EndpointService provides methods for managing application endpoints.
type EndpointService interface {
	// List returns all endpoints for the application.
	List(ctx context.Context) (*EndpointListResponse, error)

	// Create adds a new endpoint to a container.
	Create(ctx context.Context, containerID string, req *EndpointRequest) (*EndpointIDResponse, error)

	// Update updates an existing endpoint.
	Update(ctx context.Context, endpointID string, req *EndpointRequest) error

	// Delete removes an endpoint.
	Delete(ctx context.Context, endpointID string) error
}

type endpointService struct {
	client httpClient
	appID  string
}

func newEndpointService(client httpClient, appID string) EndpointService {
	return &endpointService{
		client: client,
		appID:  appID,
	}
}

// List returns all endpoints for the application.
func (s *endpointService) List(ctx context.Context) (*EndpointListResponse, error) {
	path := fmt.Sprintf("/apps/%s/endpoints", s.appID)

	var resp EndpointListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create adds a new endpoint to a container.
func (s *endpointService) Create(ctx context.Context, containerID string, req *EndpointRequest) (*EndpointIDResponse, error) {
	path := fmt.Sprintf("/apps/%s/containers/%s/endpoints", s.appID, containerID)

	var resp EndpointIDResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates an existing endpoint.
func (s *endpointService) Update(ctx context.Context, endpointID string, req *EndpointRequest) error {
	path := fmt.Sprintf("/apps/%s/endpoints/%s", s.appID, endpointID)
	return s.client.do(ctx, http.MethodPut, path, req, nil)
}

// Delete removes an endpoint.
func (s *endpointService) Delete(ctx context.Context, endpointID string) error {
	path := fmt.Sprintf("/apps/%s/endpoints/%s", s.appID, endpointID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}
