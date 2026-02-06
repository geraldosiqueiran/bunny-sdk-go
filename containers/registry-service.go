package containers

import (
	"context"
	"fmt"
	"net/http"
)

// RegistryService provides methods for managing container registries.
type RegistryService interface {
	// List returns all container registries.
	List(ctx context.Context) (*ContainerRegistryListResponse, error)

	// Get returns a specific registry by ID.
	Get(ctx context.Context, registryID int64) (*ContainerRegistry, error)

	// Create adds a new container registry.
	Create(ctx context.Context, req *CreateRegistryRequest) (*RegistryOperationResponse, error)

	// Update updates an existing registry.
	Update(ctx context.Context, registryID int64, req *UpdateRegistryRequest) (*RegistryOperationResponse, error)

	// Delete removes a registry.
	Delete(ctx context.Context, registryID int64) (*RegistryDeleteResponse, error)

	// ListImages lists container images in a registry.
	ListImages(ctx context.Context, req *ListImagesRequest) ([]ContainerImage, error)

	// ListTags lists tags for a container image.
	ListTags(ctx context.Context, req *ListTagsRequest) ([]ImageTag, error)

	// GetDigest gets the digest for a specific image tag.
	GetDigest(ctx context.Context, req *GetDigestRequest) (*ImageDigest, error)

	// GetConfigSuggestions gets configuration suggestions for an image.
	GetConfigSuggestions(ctx context.Context, req *GetConfigSuggestionsRequest) (*ConfigSuggestions, error)

	// SearchPublicImages searches for public container images.
	SearchPublicImages(ctx context.Context, req *SearchPublicImagesRequest) ([]ContainerImage, error)
}

type registryService struct {
	client httpClient
}

func newRegistryService(client httpClient) RegistryService {
	return &registryService{client: client}
}

// List returns all container registries.
func (s *registryService) List(ctx context.Context) (*ContainerRegistryListResponse, error) {
	var resp ContainerRegistryListResponse
	if err := s.client.do(ctx, http.MethodGet, "/registries", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a specific registry by ID.
func (s *registryService) Get(ctx context.Context, registryID int64) (*ContainerRegistry, error) {
	path := fmt.Sprintf("/registries/%d", registryID)

	var registry ContainerRegistry
	if err := s.client.do(ctx, http.MethodGet, path, nil, &registry); err != nil {
		return nil, err
	}
	return &registry, nil
}

// Create adds a new container registry.
func (s *registryService) Create(ctx context.Context, req *CreateRegistryRequest) (*RegistryOperationResponse, error) {
	var resp RegistryOperationResponse
	if err := s.client.do(ctx, http.MethodPost, "/registries", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates an existing registry.
func (s *registryService) Update(ctx context.Context, registryID int64, req *UpdateRegistryRequest) (*RegistryOperationResponse, error) {
	path := fmt.Sprintf("/registries/%d", registryID)

	var resp RegistryOperationResponse
	if err := s.client.do(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete removes a registry.
func (s *registryService) Delete(ctx context.Context, registryID int64) (*RegistryDeleteResponse, error) {
	path := fmt.Sprintf("/registries/%d", registryID)

	var resp RegistryDeleteResponse
	if err := s.client.do(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListImages lists container images in a registry.
func (s *registryService) ListImages(ctx context.Context, req *ListImagesRequest) ([]ContainerImage, error) {
	var images []ContainerImage
	if err := s.client.do(ctx, http.MethodPost, "/registries/images", req, &images); err != nil {
		return nil, err
	}
	return images, nil
}

// ListTags lists tags for a container image.
func (s *registryService) ListTags(ctx context.Context, req *ListTagsRequest) ([]ImageTag, error) {
	var tags []ImageTag
	if err := s.client.do(ctx, http.MethodPost, "/registries/tags", req, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// GetDigest gets the digest for a specific image tag.
func (s *registryService) GetDigest(ctx context.Context, req *GetDigestRequest) (*ImageDigest, error) {
	var digest ImageDigest
	if err := s.client.do(ctx, http.MethodPost, "/registries/digest", req, &digest); err != nil {
		return nil, err
	}
	return &digest, nil
}

// GetConfigSuggestions gets configuration suggestions for an image.
func (s *registryService) GetConfigSuggestions(ctx context.Context, req *GetConfigSuggestionsRequest) (*ConfigSuggestions, error) {
	var suggestions ConfigSuggestions
	if err := s.client.do(ctx, http.MethodPost, "/registries/config-suggestions", req, &suggestions); err != nil {
		return nil, err
	}
	return &suggestions, nil
}

// SearchPublicImages searches for public container images.
func (s *registryService) SearchPublicImages(ctx context.Context, req *SearchPublicImagesRequest) ([]ContainerImage, error) {
	var images []ContainerImage
	if err := s.client.do(ctx, http.MethodPost, "/registries/public-images/search", req, &images); err != nil {
		return nil, err
	}
	return images, nil
}
