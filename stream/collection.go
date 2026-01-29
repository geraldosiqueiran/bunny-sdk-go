package stream

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// CollectionService provides methods for managing collections within a library.
type CollectionService interface {
	List(ctx context.Context, opts *CollectionListOptions) (*CollectionListResponse, error)
	Get(ctx context.Context, collectionID string) (*Collection, error)
	Create(ctx context.Context, req *CreateCollectionRequest) (*Collection, error)
	Update(ctx context.Context, collectionID string, req *UpdateCollectionRequest) (*Collection, error)
	Delete(ctx context.Context, collectionID string) error
}

// CollectionListResponse represents a paginated list of collections.
type CollectionListResponse struct {
	ItemsPerPage int          `json:"itemsPerPage"`
	CurrentPage  int          `json:"currentPage"`
	TotalItems   int          `json:"totalItems"`
	Items        []Collection `json:"items"`
}

// HasMore returns true if there are more pages to fetch.
func (r *CollectionListResponse) HasMore() bool {
	return r.CurrentPage*r.ItemsPerPage < r.TotalItems
}

type collectionService struct {
	client    httpClient
	libraryID int64
}

func newCollectionService(client httpClient, libraryID int64) CollectionService {
	return &collectionService{client: client, libraryID: libraryID}
}

// List returns a paginated list of collections in the library.
func (s *collectionService) List(ctx context.Context, opts *CollectionListOptions) (*CollectionListResponse, error) {
	path := fmt.Sprintf("/library/%d/collections", s.libraryID)
	if opts != nil {
		path = path + "?" + buildCollectionListQuery(opts)
	}

	var resp CollectionListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single collection by ID.
func (s *collectionService) Get(ctx context.Context, collectionID string) (*Collection, error) {
	path := fmt.Sprintf("/library/%d/collections/%s", s.libraryID, collectionID)

	var collection Collection
	if err := s.client.do(ctx, http.MethodGet, path, nil, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

// Create creates a new collection in the library.
func (s *collectionService) Create(ctx context.Context, req *CreateCollectionRequest) (*Collection, error) {
	path := fmt.Sprintf("/library/%d/collections", s.libraryID)

	var collection Collection
	if err := s.client.do(ctx, http.MethodPost, path, req, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

// Update updates a collection's name.
func (s *collectionService) Update(ctx context.Context, collectionID string, req *UpdateCollectionRequest) (*Collection, error) {
	path := fmt.Sprintf("/library/%d/collections/%s", s.libraryID, collectionID)

	var collection Collection
	if err := s.client.do(ctx, http.MethodPost, path, req, &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

// Delete permanently deletes a collection.
func (s *collectionService) Delete(ctx context.Context, collectionID string) error {
	path := fmt.Sprintf("/library/%d/collections/%s", s.libraryID, collectionID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

func buildCollectionListQuery(opts *CollectionListOptions) string {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.ItemsPerPage > 0 {
		params.Set("itemsPerPage", strconv.Itoa(opts.ItemsPerPage))
	}
	if opts.OrderBy != "" {
		params.Set("orderBy", opts.OrderBy)
	}
	if opts.IncludeThumbnails {
		params.Set("includeThumbnails", "true")
	}
	return params.Encode()
}
