package shield

import (
	"context"
	"fmt"
	"net/http"
)

// RateLimitService provides methods for managing rate limit rules.
type RateLimitService interface {
	List(ctx context.Context) (*RateLimitListResponse, error)
	Create(ctx context.Context, req *CreateRateLimitRequest) (*RateLimit, error)
	Get(ctx context.Context, rateLimitID string) (*RateLimit, error)
	Update(ctx context.Context, rateLimitID string, req *UpdateRateLimitRequest) (*RateLimit, error)
	Delete(ctx context.Context, rateLimitID string) error
}

type rateLimitService struct {
	client httpClient
}

func newRateLimitService(client httpClient) RateLimitService {
	return &rateLimitService{client: client}
}

// List returns all rate limit rules.
func (s *rateLimitService) List(ctx context.Context) (*RateLimitListResponse, error) {
	var resp RateLimitListResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/rate-limits", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new rate limit rule.
func (s *rateLimitService) Create(ctx context.Context, req *CreateRateLimitRequest) (*RateLimit, error) {
	var rateLimit RateLimit
	if err := s.client.do(ctx, http.MethodPost, "/shield/rate-limit", req, &rateLimit); err != nil {
		return nil, err
	}
	return &rateLimit, nil
}

// Get returns a specific rate limit rule by ID.
func (s *rateLimitService) Get(ctx context.Context, rateLimitID string) (*RateLimit, error) {
	path := fmt.Sprintf("/shield/rate-limit/%s", rateLimitID)
	var rateLimit RateLimit
	if err := s.client.do(ctx, http.MethodGet, path, nil, &rateLimit); err != nil {
		return nil, err
	}
	return &rateLimit, nil
}

// Update updates a rate limit rule (PATCH).
func (s *rateLimitService) Update(ctx context.Context, rateLimitID string, req *UpdateRateLimitRequest) (*RateLimit, error) {
	path := fmt.Sprintf("/shield/rate-limit/%s", rateLimitID)
	var rateLimit RateLimit
	if err := s.client.do(ctx, http.MethodPatch, path, req, &rateLimit); err != nil {
		return nil, err
	}
	return &rateLimit, nil
}

// Delete deletes a rate limit rule.
func (s *rateLimitService) Delete(ctx context.Context, rateLimitID string) error {
	path := fmt.Sprintf("/shield/rate-limit/%s", rateLimitID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}
