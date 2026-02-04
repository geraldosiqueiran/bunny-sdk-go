package scripting

import (
	"context"
	"fmt"
	"net/http"
)

// VariableService provides methods for managing edge script variables.
type VariableService interface {
	Add(ctx context.Context, req *AddVariableRequest) (*EdgeScriptVariable, error)
	Get(ctx context.Context, variableID int64) (*EdgeScriptVariable, error)
	Update(ctx context.Context, variableID int64, req *UpdateVariableRequest) (*EdgeScriptVariable, error)
	Upsert(ctx context.Context, req *UpsertVariableRequest) (*EdgeScriptVariable, error)
	Delete(ctx context.Context, variableID int64) error
}

type variableService struct {
	client   httpClient
	scriptID int64
}

func newVariableService(client httpClient, scriptID int64) VariableService {
	return &variableService{client: client, scriptID: scriptID}
}

// Add creates a new variable for the edge script.
func (s *variableService) Add(ctx context.Context, req *AddVariableRequest) (*EdgeScriptVariable, error) {
	path := fmt.Sprintf("/compute/script/%d/variables/add", s.scriptID)
	var variable EdgeScriptVariable
	if err := s.client.do(ctx, http.MethodPost, path, req, &variable); err != nil {
		return nil, err
	}
	return &variable, nil
}

// Get returns a specific variable by ID.
func (s *variableService) Get(ctx context.Context, variableID int64) (*EdgeScriptVariable, error) {
	path := fmt.Sprintf("/compute/script/%d/variables/%d", s.scriptID, variableID)
	var variable EdgeScriptVariable
	if err := s.client.do(ctx, http.MethodGet, path, nil, &variable); err != nil {
		return nil, err
	}
	return &variable, nil
}

// Update updates an existing variable.
func (s *variableService) Update(ctx context.Context, variableID int64, req *UpdateVariableRequest) (*EdgeScriptVariable, error) {
	path := fmt.Sprintf("/compute/script/%d/variables/%d", s.scriptID, variableID)
	var variable EdgeScriptVariable
	if err := s.client.do(ctx, http.MethodPost, path, req, &variable); err != nil {
		return nil, err
	}
	return &variable, nil
}

// Upsert creates or updates a variable by name.
func (s *variableService) Upsert(ctx context.Context, req *UpsertVariableRequest) (*EdgeScriptVariable, error) {
	path := fmt.Sprintf("/compute/script/%d/variables", s.scriptID)
	var variable EdgeScriptVariable
	// Upsert returns 200 for new, 204 for update - try to decode response
	if err := s.client.do(ctx, http.MethodPut, path, req, &variable); err != nil {
		return nil, err
	}
	return &variable, nil
}

// Delete deletes a variable.
func (s *variableService) Delete(ctx context.Context, variableID int64) error {
	path := fmt.Sprintf("/compute/script/%d/variables/%d", s.scriptID, variableID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}
