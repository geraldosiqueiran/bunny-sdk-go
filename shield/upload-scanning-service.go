package shield

import (
	"context"
	"fmt"
	"net/http"
)

// UploadScanningService provides methods for managing upload scanning settings.
type UploadScanningService interface {
	Get(ctx context.Context) (*UploadScanningConfig, error)
	Update(ctx context.Context, req *UpdateUploadScanningRequest) (*UploadScanningConfig, error)
}

type uploadScanningService struct {
	client httpClient
	zoneID string
}

func newUploadScanningService(client httpClient, zoneID string) UploadScanningService {
	return &uploadScanningService{client: client, zoneID: zoneID}
}

// Get returns upload scanning configuration for the zone.
func (s *uploadScanningService) Get(ctx context.Context) (*UploadScanningConfig, error) {
	path := fmt.Sprintf("/shield/zone/%s/upload-scanning", s.zoneID)
	var resp UploadScanningConfig
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates upload scanning configuration for the zone.
func (s *uploadScanningService) Update(ctx context.Context, req *UpdateUploadScanningRequest) (*UploadScanningConfig, error) {
	path := fmt.Sprintf("/shield/zone/%s/upload-scanning", s.zoneID)
	var resp UploadScanningConfig
	if err := s.client.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
