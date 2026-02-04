package shield

import (
	"context"
	"fmt"
	"net/http"
)

// ZoneService provides methods for managing Shield zones.
type ZoneService interface {
	List(ctx context.Context) (*ZoneListResponse, error)
	Create(ctx context.Context, req *CreateZoneRequest) (*ShieldZone, error)
	Get(ctx context.Context, zoneID string) (*ShieldZone, error)
	Update(ctx context.Context, zoneID string, req *UpdateZoneRequest) (*ShieldZone, error)
	GetByPullZone(ctx context.Context, pullZoneID int64) (*ShieldZone, error)
	GetPullZoneMapping(ctx context.Context) (*PullZoneMappingResponse, error)
}

type zoneService struct {
	client httpClient
}

func newZoneService(client httpClient) ZoneService {
	return &zoneService{client: client}
}

// List returns all Shield zones.
func (s *zoneService) List(ctx context.Context) (*ZoneListResponse, error) {
	var resp ZoneListResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/zones", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new Shield zone.
func (s *zoneService) Create(ctx context.Context, req *CreateZoneRequest) (*ShieldZone, error) {
	var zone ShieldZone
	if err := s.client.do(ctx, http.MethodPost, "/shield/zone", req, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// Get returns a specific Shield zone by ID.
func (s *zoneService) Get(ctx context.Context, zoneID string) (*ShieldZone, error) {
	path := fmt.Sprintf("/shield/zone/%s", zoneID)
	var zone ShieldZone
	if err := s.client.do(ctx, http.MethodGet, path, nil, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// Update updates a Shield zone (PATCH).
func (s *zoneService) Update(ctx context.Context, zoneID string, req *UpdateZoneRequest) (*ShieldZone, error) {
	path := fmt.Sprintf("/shield/zone/%s", zoneID)
	var zone ShieldZone
	if err := s.client.do(ctx, http.MethodPatch, path, req, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetByPullZone returns the Shield zone associated with a Pull zone ID.
func (s *zoneService) GetByPullZone(ctx context.Context, pullZoneID int64) (*ShieldZone, error) {
	path := fmt.Sprintf("/shield/zone/pullzone/%d", pullZoneID)
	var zone ShieldZone
	if err := s.client.do(ctx, http.MethodGet, path, nil, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// GetPullZoneMapping returns the Shield zone to Pull zone mapping.
func (s *zoneService) GetPullZoneMapping(ctx context.Context) (*PullZoneMappingResponse, error) {
	var resp PullZoneMappingResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/zones/pullzone-mapping", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
