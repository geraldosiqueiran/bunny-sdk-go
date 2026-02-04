package shield

import (
	"context"
	"net/http"
)

// PromoService provides methods for promotional information.
type PromoService interface {
	Get(ctx context.Context) (*PromoInfo, error)
}

type promoService struct {
	client httpClient
}

func newPromoService(client httpClient) PromoService {
	return &promoService{client: client}
}

// Get returns promotional information.
func (s *promoService) Get(ctx context.Context) (*PromoInfo, error) {
	var resp PromoInfo
	if err := s.client.do(ctx, http.MethodGet, "/shield/promo", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
