package scripting

import (
	"context"
	"fmt"
	"net/http"
)

// SecretService provides methods for managing edge script secrets.
type SecretService interface {
	List(ctx context.Context) (*SecretListResponse, error)
	Add(ctx context.Context, req *AddSecretRequest) (*EdgeScriptSecret, error)
	Update(ctx context.Context, secretID int64, req *UpdateSecretRequest) (*EdgeScriptSecret, error)
	Upsert(ctx context.Context, req *UpsertSecretRequest) (*EdgeScriptSecret, error)
	Delete(ctx context.Context, secretID int64) error
}

type secretService struct {
	client   httpClient
	scriptID int64
}

func newSecretService(client httpClient, scriptID int64) SecretService {
	return &secretService{client: client, scriptID: scriptID}
}

// List returns all secrets for the edge script.
func (s *secretService) List(ctx context.Context) (*SecretListResponse, error) {
	path := fmt.Sprintf("/compute/script/%d/secrets", s.scriptID)
	var resp SecretListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Add creates a new secret for the edge script.
func (s *secretService) Add(ctx context.Context, req *AddSecretRequest) (*EdgeScriptSecret, error) {
	path := fmt.Sprintf("/compute/script/%d/secrets", s.scriptID)
	var secret EdgeScriptSecret
	if err := s.client.do(ctx, http.MethodPost, path, req, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

// Update updates an existing secret.
func (s *secretService) Update(ctx context.Context, secretID int64, req *UpdateSecretRequest) (*EdgeScriptSecret, error) {
	path := fmt.Sprintf("/compute/script/%d/secrets/%d", s.scriptID, secretID)
	var secret EdgeScriptSecret
	if err := s.client.do(ctx, http.MethodPost, path, req, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

// Upsert creates or updates a secret by name.
func (s *secretService) Upsert(ctx context.Context, req *UpsertSecretRequest) (*EdgeScriptSecret, error) {
	path := fmt.Sprintf("/compute/script/%d/secrets", s.scriptID)
	var secret EdgeScriptSecret
	// Upsert returns 200 for new, 204 for update - try to decode response
	if err := s.client.do(ctx, http.MethodPut, path, req, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

// Delete deletes a secret.
func (s *secretService) Delete(ctx context.Context, secretID int64) error {
	path := fmt.Sprintf("/compute/script/%d/secrets/%d", s.scriptID, secretID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}
