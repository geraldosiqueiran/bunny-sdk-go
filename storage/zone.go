package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// httpClient is the internal interface for making API requests.
type httpClient interface {
	do(ctx context.Context, method, path string, body any, result any) error
}

// ZoneService provides methods for managing storage zones.
type ZoneService interface {
	List(ctx context.Context, opts *ZoneListOptions) (*ZoneListResponse, error)
	Get(ctx context.Context, zoneID int64) (*Zone, error)
	Create(ctx context.Context, req *CreateZoneRequest) (*Zone, error)
	Update(ctx context.Context, zoneID int64, req *UpdateZoneRequest) (*Zone, error)
	Delete(ctx context.Context, zoneID int64) error
	CheckAvailability(ctx context.Context, name string) (*AvailabilityResponse, error)
	ResetPassword(ctx context.Context, zoneID int64) (*ResetPasswordResponse, error)
	ResetReadOnlyPassword(ctx context.Context, zoneID int64) (*ResetReadOnlyPasswordResponse, error)
}

type zoneService struct {
	client httpClient
}

func newZoneService(client httpClient) ZoneService {
	return &zoneService{client: client}
}

// List returns a paginated list of storage zones.
func (s *zoneService) List(ctx context.Context, opts *ZoneListOptions) (*ZoneListResponse, error) {
	path := "/storagezone"
	if opts != nil {
		path = path + "?" + buildZoneListQuery(opts)
	}

	var resp ZoneListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single storage zone by ID.
func (s *zoneService) Get(ctx context.Context, zoneID int64) (*Zone, error) {
	path := fmt.Sprintf("/storagezone/%d", zoneID)

	var zone Zone
	if err := s.client.do(ctx, http.MethodGet, path, nil, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// Create creates a new storage zone.
func (s *zoneService) Create(ctx context.Context, req *CreateZoneRequest) (*Zone, error) {
	var zone Zone
	if err := s.client.do(ctx, http.MethodPost, "/storagezone", req, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// Update updates a storage zone's settings.
func (s *zoneService) Update(ctx context.Context, zoneID int64, req *UpdateZoneRequest) (*Zone, error) {
	path := fmt.Sprintf("/storagezone/%d", zoneID)

	var zone Zone
	if err := s.client.do(ctx, http.MethodPost, path, req, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// Delete permanently deletes a storage zone.
func (s *zoneService) Delete(ctx context.Context, zoneID int64) error {
	path := fmt.Sprintf("/storagezone/%d", zoneID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// CheckAvailability checks if a storage zone name is available.
func (s *zoneService) CheckAvailability(ctx context.Context, name string) (*AvailabilityResponse, error) {
	path := fmt.Sprintf("/storagezone/checkavailability/%s", url.PathEscape(name))

	var resp AvailabilityResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetPassword resets the main password for a storage zone.
func (s *zoneService) ResetPassword(ctx context.Context, zoneID int64) (*ResetPasswordResponse, error) {
	path := fmt.Sprintf("/storagezone/%d/resetPassword", zoneID)

	var resp ResetPasswordResponse
	// API requires empty JSON body
	if err := s.client.do(ctx, http.MethodPost, path, struct{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetReadOnlyPassword resets the read-only password for a storage zone.
func (s *zoneService) ResetReadOnlyPassword(ctx context.Context, zoneID int64) (*ResetReadOnlyPasswordResponse, error) {
	path := fmt.Sprintf("/storagezone/%d/resetReadOnlyPassword", zoneID)

	var resp ResetReadOnlyPasswordResponse
	// API requires empty JSON body
	if err := s.client.do(ctx, http.MethodPost, path, struct{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildZoneListQuery(opts *ZoneListOptions) string {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("perPage", strconv.Itoa(opts.PerPage))
	}
	if opts.IncludeDeleted {
		params.Set("includeDeleted", "true")
	}
	if opts.Search != "" {
		params.Set("search", opts.Search)
	}
	return params.Encode()
}
