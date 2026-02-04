package scripting

import (
	"context"
	"fmt"
	"net/http"
)

// CodeService provides methods for managing edge script code.
type CodeService interface {
	Get(ctx context.Context) (*EdgeScriptCode, error)
	Set(ctx context.Context, req *UpdateCodeRequest) error
}

type codeService struct {
	client   httpClient
	scriptID int64
}

func newCodeService(client httpClient, scriptID int64) CodeService {
	return &codeService{client: client, scriptID: scriptID}
}

// Get returns the code for the edge script.
func (s *codeService) Get(ctx context.Context) (*EdgeScriptCode, error) {
	path := fmt.Sprintf("/compute/script/%d/code", s.scriptID)
	var code EdgeScriptCode
	if err := s.client.do(ctx, http.MethodGet, path, nil, &code); err != nil {
		return nil, err
	}
	return &code, nil
}

// Set sets the code for the edge script.
func (s *codeService) Set(ctx context.Context, req *UpdateCodeRequest) error {
	path := fmt.Sprintf("/compute/script/%d/code", s.scriptID)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}
