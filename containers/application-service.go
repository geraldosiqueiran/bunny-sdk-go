package containers

import (
	"context"
	"fmt"
	"net/http"
)

// ApplicationService provides methods for managing Magic Containers applications.
type ApplicationService interface {
	// List returns all applications with optional pagination.
	List(ctx context.Context, opts *ListOptions) (*ApplicationListResponse, error)

	// Get returns a specific application by ID.
	Get(ctx context.Context, appID string) (*Application, error)

	// Create creates a new application.
	Create(ctx context.Context, req *CreateApplicationRequest) (*ApplicationIDResponse, error)

	// Update performs a full update of an application (PUT).
	Update(ctx context.Context, appID string, req *UpdateApplicationRequest) (*ApplicationIDResponse, error)

	// Patch performs a partial update of an application (PATCH).
	Patch(ctx context.Context, appID string, req *PatchApplicationRequest) (*ApplicationIDResponse, error)

	// Delete deletes an application.
	Delete(ctx context.Context, appID string) error

	// Deploy deploys an application.
	Deploy(ctx context.Context, appID string) error

	// Undeploy undeploys an application.
	Undeploy(ctx context.Context, appID string) error

	// Restart restarts an application.
	Restart(ctx context.Context, appID string) error

	// GetOverview returns overview metrics for an application.
	GetOverview(ctx context.Context, appID string) (*ApplicationOverview, error)

	// GetStatistics returns statistics for an application.
	GetStatistics(ctx context.Context, appID string, opts *StatisticsOptions) (*ApplicationStatistics, error)
}

type applicationService struct {
	client httpClient
}

func newApplicationService(client httpClient) ApplicationService {
	return &applicationService{client: client}
}

// List returns all applications with optional pagination.
func (s *applicationService) List(ctx context.Context, opts *ListOptions) (*ApplicationListResponse, error) {
	path := "/apps"
	if opts != nil {
		path = buildListPath(path, opts)
	}

	var resp ApplicationListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a specific application by ID.
func (s *applicationService) Get(ctx context.Context, appID string) (*Application, error) {
	path := fmt.Sprintf("/apps/%s", appID)

	var app Application
	if err := s.client.do(ctx, http.MethodGet, path, nil, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

// Create creates a new application.
func (s *applicationService) Create(ctx context.Context, req *CreateApplicationRequest) (*ApplicationIDResponse, error) {
	var resp ApplicationIDResponse
	if err := s.client.do(ctx, http.MethodPost, "/apps", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update performs a full update of an application (PUT).
func (s *applicationService) Update(ctx context.Context, appID string, req *UpdateApplicationRequest) (*ApplicationIDResponse, error) {
	path := fmt.Sprintf("/apps/%s", appID)

	var resp ApplicationIDResponse
	if err := s.client.do(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Patch performs a partial update of an application (PATCH).
func (s *applicationService) Patch(ctx context.Context, appID string, req *PatchApplicationRequest) (*ApplicationIDResponse, error) {
	path := fmt.Sprintf("/apps/%s", appID)

	var resp ApplicationIDResponse
	if err := s.client.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes an application.
func (s *applicationService) Delete(ctx context.Context, appID string) error {
	path := fmt.Sprintf("/apps/%s", appID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// Deploy deploys an application.
func (s *applicationService) Deploy(ctx context.Context, appID string) error {
	path := fmt.Sprintf("/apps/%s/deploy", appID)
	return s.client.do(ctx, http.MethodPost, path, nil, nil)
}

// Undeploy undeploys an application.
func (s *applicationService) Undeploy(ctx context.Context, appID string) error {
	path := fmt.Sprintf("/apps/%s/undeploy", appID)
	return s.client.do(ctx, http.MethodPost, path, nil, nil)
}

// Restart restarts an application.
func (s *applicationService) Restart(ctx context.Context, appID string) error {
	path := fmt.Sprintf("/apps/%s/restart", appID)
	return s.client.do(ctx, http.MethodPost, path, nil, nil)
}

// GetOverview returns overview metrics for an application.
func (s *applicationService) GetOverview(ctx context.Context, appID string) (*ApplicationOverview, error) {
	path := fmt.Sprintf("/apps/%s/overview", appID)

	var overview ApplicationOverview
	if err := s.client.do(ctx, http.MethodGet, path, nil, &overview); err != nil {
		return nil, err
	}
	return &overview, nil
}

// GetStatistics returns statistics for an application.
func (s *applicationService) GetStatistics(ctx context.Context, appID string, opts *StatisticsOptions) (*ApplicationStatistics, error) {
	path := fmt.Sprintf("/apps/%s/statistics", appID)
	if opts != nil {
		path = buildStatisticsPath(path, opts)
	}

	var stats ApplicationStatistics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
