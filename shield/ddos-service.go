package shield

import (
	"context"
	"net/http"
)

// DDoSService provides methods for DDoS protection information.
type DDoSService interface {
	GetEnums(ctx context.Context) (*DDoSEnums, error)
}

type ddosService struct {
	client httpClient
}

func newDDoSService(client httpClient) DDoSService {
	return &ddosService{client: client}
}

// GetEnums returns DDoS protection enumeration values.
func (s *ddosService) GetEnums(ctx context.Context) (*DDoSEnums, error) {
	var resp DDoSEnums
	if err := s.client.do(ctx, http.MethodGet, "/shield/ddos/enums", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
