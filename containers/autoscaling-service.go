package containers

import (
	"context"
	"fmt"
	"net/http"
)

// AutoscalingService provides methods for managing application autoscaling.
type AutoscalingService interface {
	// Get returns the autoscaling settings for an application.
	Get(ctx context.Context) (*AutoScaling, error)

	// Update updates the autoscaling settings for an application.
	Update(ctx context.Context, req *AutoScaling) error
}

type autoscalingService struct {
	client httpClient
	appID  string
}

func newAutoscalingService(client httpClient, appID string) AutoscalingService {
	return &autoscalingService{
		client: client,
		appID:  appID,
	}
}

// Get returns the autoscaling settings for an application.
func (s *autoscalingService) Get(ctx context.Context) (*AutoScaling, error) {
	path := fmt.Sprintf("/apps/%s/autoscaling", s.appID)

	var settings AutoScaling
	if err := s.client.do(ctx, http.MethodGet, path, nil, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// Update updates the autoscaling settings for an application.
func (s *autoscalingService) Update(ctx context.Context, req *AutoScaling) error {
	path := fmt.Sprintf("/apps/%s/autoscaling", s.appID)
	return s.client.do(ctx, http.MethodPut, path, req, nil)
}
