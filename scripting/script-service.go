package scripting

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ScriptService provides methods for managing edge scripts.
type ScriptService interface {
	List(ctx context.Context, opts *ScriptListOptions) (*ScriptListResponse, error)
	Create(ctx context.Context, req *CreateScriptRequest) (*EdgeScript, error)
	Get(ctx context.Context, id int64) (*EdgeScript, error)
	Update(ctx context.Context, id int64, req *UpdateScriptRequest) (*EdgeScript, error)
	Delete(ctx context.Context, id int64, deleteLinkedPullZones bool) error
	GetStatistics(ctx context.Context, id int64, opts *StatisticsOptions) (*ScriptStatistics, error)
	RotateDeploymentKey(ctx context.Context, id int64) error
}

type scriptService struct {
	client httpClient
}

func newScriptService(client httpClient) ScriptService {
	return &scriptService{client: client}
}

// List returns all edge scripts.
func (s *scriptService) List(ctx context.Context, opts *ScriptListOptions) (*ScriptListResponse, error) {
	path := "/compute/script"
	if opts != nil {
		q := url.Values{}
		if opts.Page > 0 {
			q.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			q.Set("perPage", strconv.Itoa(opts.PerPage))
		}
		if opts.Search != "" {
			q.Set("search", opts.Search)
		}
		if opts.IncludeLinkedPullzones {
			q.Set("includeLinkedPullzones", "true")
		}
		if opts.IntegrationID != nil {
			q.Set("integrationId", strconv.FormatInt(*opts.IntegrationID, 10))
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
	}

	var resp ScriptListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new edge script.
func (s *scriptService) Create(ctx context.Context, req *CreateScriptRequest) (*EdgeScript, error) {
	var script EdgeScript
	if err := s.client.do(ctx, http.MethodPost, "/compute/script", req, &script); err != nil {
		return nil, err
	}
	return &script, nil
}

// Get returns a specific edge script by ID.
func (s *scriptService) Get(ctx context.Context, id int64) (*EdgeScript, error) {
	path := fmt.Sprintf("/compute/script/%d", id)
	var script EdgeScript
	if err := s.client.do(ctx, http.MethodGet, path, nil, &script); err != nil {
		return nil, err
	}
	return &script, nil
}

// Update updates an edge script.
func (s *scriptService) Update(ctx context.Context, id int64, req *UpdateScriptRequest) (*EdgeScript, error) {
	path := fmt.Sprintf("/compute/script/%d", id)
	var script EdgeScript
	if err := s.client.do(ctx, http.MethodPost, path, req, &script); err != nil {
		return nil, err
	}
	return &script, nil
}

// Delete deletes an edge script.
func (s *scriptService) Delete(ctx context.Context, id int64, deleteLinkedPullZones bool) error {
	path := fmt.Sprintf("/compute/script/%d", id)
	if deleteLinkedPullZones {
		path += "?deleteLinkedPullZones=true"
	}
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// GetStatistics returns statistics for an edge script.
func (s *scriptService) GetStatistics(ctx context.Context, id int64, opts *StatisticsOptions) (*ScriptStatistics, error) {
	path := fmt.Sprintf("/compute/script/%d/statistics", id)
	if opts != nil {
		q := url.Values{}
		if opts.DateFrom != "" {
			q.Set("dateFrom", opts.DateFrom)
		}
		if opts.DateTo != "" {
			q.Set("dateTo", opts.DateTo)
		}
		if opts.LoadLatest {
			q.Set("loadLatest", "true")
		}
		if opts.Hourly {
			q.Set("hourly", "true")
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
	}

	var stats ScriptStatistics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// RotateDeploymentKey rotates the deployment key for an edge script.
func (s *scriptService) RotateDeploymentKey(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/compute/script/%d/deploymentKey/rotate", id)
	return s.client.do(ctx, http.MethodPost, path, nil, nil)
}
