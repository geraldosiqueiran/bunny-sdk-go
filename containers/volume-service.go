package containers

import (
	"context"
	"fmt"
	"net/http"
)

// VolumeService provides methods for managing application volumes.
type VolumeService interface {
	// List returns all volumes for the application.
	List(ctx context.Context) (*VolumeListResponse, error)

	// Update updates a volume's name or size.
	Update(ctx context.Context, volumeID string, req *UpdateVolumeRequest) (*VolumeUpdateResponse, error)

	// Detach detaches a volume from all pods.
	Detach(ctx context.Context, volumeID string) (*VolumeNameResponse, error)

	// DeleteInstance deletes a specific volume instance.
	DeleteInstance(ctx context.Context, volumeID string, instanceID string) (*VolumeInstanceIDResponse, error)

	// DeleteAll deletes all instances of a volume.
	// Note: All instances must be detached before deletion.
	DeleteAll(ctx context.Context, volumeID string) (*VolumeInstanceIDsResponse, error)
}

type volumeService struct {
	client httpClient
	appID  string
}

func newVolumeService(client httpClient, appID string) VolumeService {
	return &volumeService{
		client: client,
		appID:  appID,
	}
}

// List returns all volumes for the application.
func (s *volumeService) List(ctx context.Context) (*VolumeListResponse, error) {
	path := fmt.Sprintf("/apps/%s/volumes", s.appID)

	var resp VolumeListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a volume's name or size.
func (s *volumeService) Update(ctx context.Context, volumeID string, req *UpdateVolumeRequest) (*VolumeUpdateResponse, error) {
	path := fmt.Sprintf("/apps/%s/volumes/%s", s.appID, volumeID)

	var resp VolumeUpdateResponse
	if err := s.client.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Detach detaches a volume from all pods.
func (s *volumeService) Detach(ctx context.Context, volumeID string) (*VolumeNameResponse, error) {
	path := fmt.Sprintf("/apps/%s/volumes/%s/detach", s.appID, volumeID)

	var resp VolumeNameResponse
	if err := s.client.do(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteInstance deletes a specific volume instance.
func (s *volumeService) DeleteInstance(ctx context.Context, volumeID string, instanceID string) (*VolumeInstanceIDResponse, error) {
	path := fmt.Sprintf("/apps/%s/volumes/%s/instances/%s", s.appID, volumeID, instanceID)

	var resp VolumeInstanceIDResponse
	if err := s.client.do(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAll deletes all instances of a volume.
// Note: All instances must be detached before deletion.
func (s *volumeService) DeleteAll(ctx context.Context, volumeID string) (*VolumeInstanceIDsResponse, error) {
	path := fmt.Sprintf("/apps/%s/volumes/%s", s.appID, volumeID)

	var resp VolumeInstanceIDsResponse
	if err := s.client.do(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
