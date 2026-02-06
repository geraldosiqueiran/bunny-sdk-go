package containers

import (
	"context"
	"fmt"
	"net/http"
)

// LimitsService provides methods for retrieving user limits.
type LimitsService interface {
	// Get returns the current user's limits.
	Get(ctx context.Context) (*UserLimits, error)
}

type limitsService struct {
	client httpClient
}

func newLimitsService(client httpClient) LimitsService {
	return &limitsService{client: client}
}

// Get returns the current user's limits.
func (s *limitsService) Get(ctx context.Context) (*UserLimits, error) {
	var limits UserLimits
	if err := s.client.do(ctx, http.MethodGet, "/limits", nil, &limits); err != nil {
		return nil, err
	}
	return &limits, nil
}

// NodeService provides methods for listing nodes.
type NodeService interface {
	// List returns all available nodes.
	List(ctx context.Context, opts *ListOptions) (*NodeListResponse, error)
}

type nodeService struct {
	client httpClient
}

func newNodeService(client httpClient) NodeService {
	return &nodeService{client: client}
}

// List returns all available nodes.
func (s *nodeService) List(ctx context.Context, opts *ListOptions) (*NodeListResponse, error) {
	path := "/nodes"
	if opts != nil {
		path = buildListPath(path, opts)
	}

	var resp NodeListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PodService provides methods for managing pods.
type PodService interface {
	// Recreate triggers recreation of a specific pod.
	Recreate(ctx context.Context, podID string) error
}

type podService struct {
	client httpClient
	appID  string
}

func newPodService(client httpClient, appID string) PodService {
	return &podService{
		client: client,
		appID:  appID,
	}
}

// Recreate triggers recreation of a specific pod.
func (s *podService) Recreate(ctx context.Context, podID string) error {
	path := fmt.Sprintf("/apps/%s/pods/%s/recreate", s.appID, podID)
	return s.client.do(ctx, http.MethodPost, path, nil, nil)
}
