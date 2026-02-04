package scripting

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ReleaseService provides methods for managing edge script releases.
type ReleaseService interface {
	List(ctx context.Context, opts *ReleaseListOptions) (*ReleaseListResponse, error)
	GetActive(ctx context.Context) (*EdgeScriptRelease, error)
	Publish(ctx context.Context, req *PublishReleaseRequest) error
	PublishByUUID(ctx context.Context, uuid string, req *PublishReleaseRequest) error
}

type releaseService struct {
	client   httpClient
	scriptID int64
}

func newReleaseService(client httpClient, scriptID int64) ReleaseService {
	return &releaseService{client: client, scriptID: scriptID}
}

// List returns all releases for the edge script.
func (s *releaseService) List(ctx context.Context, opts *ReleaseListOptions) (*ReleaseListResponse, error) {
	path := fmt.Sprintf("/compute/script/%d/releases", s.scriptID)
	if opts != nil {
		q := url.Values{}
		if opts.Page > 0 {
			q.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			q.Set("perPage", strconv.Itoa(opts.PerPage))
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
	}

	var resp ReleaseListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetActive returns the active release for the edge script.
func (s *releaseService) GetActive(ctx context.Context) (*EdgeScriptRelease, error) {
	path := fmt.Sprintf("/compute/script/%d/releases/active", s.scriptID)
	var release EdgeScriptRelease
	if err := s.client.do(ctx, http.MethodGet, path, nil, &release); err != nil {
		return nil, err
	}
	return &release, nil
}

// Publish publishes a new release for the edge script.
func (s *releaseService) Publish(ctx context.Context, req *PublishReleaseRequest) error {
	path := fmt.Sprintf("/compute/script/%d/publish", s.scriptID)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}

// PublishByUUID publishes a specific release by UUID.
func (s *releaseService) PublishByUUID(ctx context.Context, uuid string, req *PublishReleaseRequest) error {
	path := fmt.Sprintf("/compute/script/%d/publish/%s", s.scriptID, uuid)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}
