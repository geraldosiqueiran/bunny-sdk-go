package shield

import (
	"context"
	"fmt"
	"net/http"
)

// AccessListService provides methods for managing zone-scoped access lists.
type AccessListService interface {
	Get(ctx context.Context) (*AccessList, error)
	Add(ctx context.Context, req *AddAccessListEntryRequest) (*AccessListEntry, error)
	Update(ctx context.Context, req *UpdateAccessListEntriesRequest) error
	Delete(ctx context.Context, req *DeleteAccessListEntriesRequest) error
	GetEnums(ctx context.Context) (*AccessListEnums, error)
	UpdateConfig(ctx context.Context, req *UpdateAccessListConfigRequest) (*AccessListConfig, error)
}

type accessListService struct {
	client httpClient
	zoneID string
}

func newAccessListService(client httpClient, zoneID string) AccessListService {
	return &accessListService{client: client, zoneID: zoneID}
}

// Get returns access lists for the zone.
func (s *accessListService) Get(ctx context.Context) (*AccessList, error) {
	path := fmt.Sprintf("/shield/zone/%s/access-lists", s.zoneID)
	var resp AccessList
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Add adds an entry to the access list.
func (s *accessListService) Add(ctx context.Context, req *AddAccessListEntryRequest) (*AccessListEntry, error) {
	path := fmt.Sprintf("/shield/zone/%s/access-lists", s.zoneID)
	var entry AccessListEntry
	if err := s.client.do(ctx, http.MethodPost, path, req, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Update updates access list entries (batch PATCH).
func (s *accessListService) Update(ctx context.Context, req *UpdateAccessListEntriesRequest) error {
	path := fmt.Sprintf("/shield/zone/%s/access-lists", s.zoneID)
	return s.client.do(ctx, http.MethodPatch, path, req, nil)
}

// Delete removes entries from the access list (DELETE with body).
func (s *accessListService) Delete(ctx context.Context, req *DeleteAccessListEntriesRequest) error {
	path := fmt.Sprintf("/shield/zone/%s/access-lists", s.zoneID)
	return s.client.do(ctx, http.MethodDelete, path, req, nil)
}

// GetEnums returns available access list types and actions.
func (s *accessListService) GetEnums(ctx context.Context) (*AccessListEnums, error) {
	path := fmt.Sprintf("/shield/zone/%s/access-lists/enums", s.zoneID)
	var resp AccessListEnums
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateConfig updates access list configuration.
func (s *accessListService) UpdateConfig(ctx context.Context, req *UpdateAccessListConfigRequest) (*AccessListConfig, error) {
	path := fmt.Sprintf("/shield/zone/%s/access-lists/configurations", s.zoneID)
	var resp AccessListConfig
	if err := s.client.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
