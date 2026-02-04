package shield

import (
	"context"
	"fmt"
	"net/http"
)

// BotDetectionService provides methods for managing bot detection settings.
type BotDetectionService interface {
	Get(ctx context.Context) (*BotDetectionSettings, error)
	Update(ctx context.Context, req *UpdateBotDetectionRequest) (*BotDetectionSettings, error)
}

type botDetectionService struct {
	client httpClient
	zoneID string
}

func newBotDetectionService(client httpClient, zoneID string) BotDetectionService {
	return &botDetectionService{client: client, zoneID: zoneID}
}

// Get returns bot detection settings for the zone.
func (s *botDetectionService) Get(ctx context.Context) (*BotDetectionSettings, error) {
	path := fmt.Sprintf("/shield/zone/%s/bot-detection", s.zoneID)
	var resp BotDetectionSettings
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates bot detection settings for the zone.
func (s *botDetectionService) Update(ctx context.Context, req *UpdateBotDetectionRequest) (*BotDetectionSettings, error) {
	path := fmt.Sprintf("/shield/zone/%s/bot-detection", s.zoneID)
	var resp BotDetectionSettings
	if err := s.client.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
