package stream

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// LibraryService provides methods for managing video libraries.
type LibraryService interface {
	List(ctx context.Context, opts *LibraryListOptions) (*LibraryListResponse, error)
	Get(ctx context.Context, libraryID int64) (*Library, error)
	Create(ctx context.Context, req *CreateLibraryRequest) (*Library, error)
	Update(ctx context.Context, libraryID int64, req *UpdateLibraryRequest) (*Library, error)
	Delete(ctx context.Context, libraryID int64) error
	GetStatistics(ctx context.Context, libraryID int64, opts *StatisticsOptions) (*LibraryStatistics, error)
}

// LibraryListResponse represents a paginated list of libraries.
type LibraryListResponse struct {
	ItemsPerPage int       `json:"itemsPerPage"`
	CurrentPage  int       `json:"currentPage"`
	TotalItems   int       `json:"totalItems"`
	Items        []Library `json:"items"`
}

// HasMore returns true if there are more pages to fetch.
func (r *LibraryListResponse) HasMore() bool {
	return r.CurrentPage*r.ItemsPerPage < r.TotalItems
}

type libraryService struct {
	client httpClient
}

func newLibraryService(client httpClient) LibraryService {
	return &libraryService{client: client}
}

// List returns a paginated list of video libraries.
func (s *libraryService) List(ctx context.Context, opts *LibraryListOptions) (*LibraryListResponse, error) {
	path := "/library"
	if opts != nil {
		path = path + "?" + buildLibraryListQuery(opts)
	}

	var resp LibraryListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single library by ID.
func (s *libraryService) Get(ctx context.Context, libraryID int64) (*Library, error) {
	path := fmt.Sprintf("/library/%d", libraryID)

	var library Library
	if err := s.client.do(ctx, http.MethodGet, path, nil, &library); err != nil {
		return nil, err
	}
	return &library, nil
}

// Create creates a new video library.
func (s *libraryService) Create(ctx context.Context, req *CreateLibraryRequest) (*Library, error) {
	var library Library
	if err := s.client.do(ctx, http.MethodPost, "/library", req, &library); err != nil {
		return nil, err
	}
	return &library, nil
}

// Update updates a library's settings.
func (s *libraryService) Update(ctx context.Context, libraryID int64, req *UpdateLibraryRequest) (*Library, error) {
	path := fmt.Sprintf("/library/%d", libraryID)

	var library Library
	if err := s.client.do(ctx, http.MethodPost, path, req, &library); err != nil {
		return nil, err
	}
	return &library, nil
}

// Delete permanently deletes a library and all its videos.
func (s *libraryService) Delete(ctx context.Context, libraryID int64) error {
	path := fmt.Sprintf("/library/%d", libraryID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// GetStatistics returns statistics for a library.
func (s *libraryService) GetStatistics(ctx context.Context, libraryID int64, opts *StatisticsOptions) (*LibraryStatistics, error) {
	path := fmt.Sprintf("/library/%d/statistics", libraryID)
	if opts != nil {
		path = path + "?" + buildStatisticsQuery(opts)
	}

	var resp LibraryStatistics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildLibraryListQuery(opts *LibraryListOptions) string {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.ItemsPerPage > 0 {
		params.Set("itemsPerPage", strconv.Itoa(opts.ItemsPerPage))
	}
	return params.Encode()
}
