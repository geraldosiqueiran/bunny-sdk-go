package containers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// RegionService provides methods for managing regions.
type RegionService interface {
	// List returns all available regions.
	List(ctx context.Context, opts *ListOptions) (*RegionListResponse, error)

	// GetOptimal returns the optimal region based on location.
	GetOptimal(ctx context.Context, cdnServerToken string) (*OptimalRegionResponse, error)
}

type regionService struct {
	client httpClient
}

func newRegionService(client httpClient) RegionService {
	return &regionService{client: client}
}

// List returns all available regions.
func (s *regionService) List(ctx context.Context, opts *ListOptions) (*RegionListResponse, error) {
	path := "/regions"
	if opts != nil {
		path = buildListPath(path, opts)
	}

	var resp RegionListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOptimal returns the optimal region based on location.
func (s *regionService) GetOptimal(ctx context.Context, cdnServerToken string) (*OptimalRegionResponse, error) {
	path := "/regions/optimal"
	if cdnServerToken != "" {
		params := url.Values{}
		params.Set("cdnServerToken", cdnServerToken)
		path = path + "?" + params.Encode()
	}

	var resp OptimalRegionResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RegionSettingsService provides methods for managing application region settings.
type RegionSettingsService interface {
	// Get returns the region settings for an application.
	Get(ctx context.Context) (*RegionSettings, error)

	// Update updates the region settings for an application.
	Update(ctx context.Context, req *UpdateRegionSettingsRequest) error
}

type regionSettingsService struct {
	client httpClient
	appID  string
}

func newRegionSettingsService(client httpClient, appID string) RegionSettingsService {
	return &regionSettingsService{
		client: client,
		appID:  appID,
	}
}

// Get returns the region settings for an application.
func (s *regionSettingsService) Get(ctx context.Context) (*RegionSettings, error) {
	path := fmt.Sprintf("/apps/%s/region-settings", s.appID)

	var settings RegionSettings
	if err := s.client.do(ctx, http.MethodGet, path, nil, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// Update updates the region settings for an application.
func (s *regionSettingsService) Update(ctx context.Context, req *UpdateRegionSettingsRequest) error {
	path := fmt.Sprintf("/apps/%s/region-settings", s.appID)
	return s.client.do(ctx, http.MethodPut, path, req, nil)
}
