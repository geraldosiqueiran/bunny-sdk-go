package containers

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerTemplateService provides methods for managing container templates.
type ContainerTemplateService interface {
	// Get returns a specific container template.
	Get(ctx context.Context, containerID string) (*ContainerTemplate, error)

	// Create adds a new container template to the application.
	Create(ctx context.Context, req *CreateContainerTemplateRequest) (*ContainerTemplate, error)

	// Patch performs a partial update of a container template.
	Patch(ctx context.Context, containerID string, req *PatchContainerTemplateRequest) (*ContainerTemplate, error)

	// Delete removes a container template from the application.
	Delete(ctx context.Context, containerID string) error

	// SetEnvironmentVariables sets environment variables for a container.
	SetEnvironmentVariables(ctx context.Context, containerID string, envVars SetEnvironmentVariablesRequest) (*ContainerTemplate, error)
}

type containerTemplateService struct {
	client httpClient
	appID  string
}

func newContainerTemplateService(client httpClient, appID string) ContainerTemplateService {
	return &containerTemplateService{
		client: client,
		appID:  appID,
	}
}

// Get returns a specific container template.
func (s *containerTemplateService) Get(ctx context.Context, containerID string) (*ContainerTemplate, error) {
	path := fmt.Sprintf("/apps/%s/containers/%s", s.appID, containerID)

	var template ContainerTemplate
	if err := s.client.do(ctx, http.MethodGet, path, nil, &template); err != nil {
		return nil, err
	}
	return &template, nil
}

// Create adds a new container template to the application.
func (s *containerTemplateService) Create(ctx context.Context, req *CreateContainerTemplateRequest) (*ContainerTemplate, error) {
	path := fmt.Sprintf("/apps/%s/containers", s.appID)

	var template ContainerTemplate
	if err := s.client.do(ctx, http.MethodPost, path, req, &template); err != nil {
		return nil, err
	}
	return &template, nil
}

// Patch performs a partial update of a container template.
func (s *containerTemplateService) Patch(ctx context.Context, containerID string, req *PatchContainerTemplateRequest) (*ContainerTemplate, error) {
	path := fmt.Sprintf("/apps/%s/containers/%s", s.appID, containerID)

	var template ContainerTemplate
	if err := s.client.do(ctx, http.MethodPatch, path, req, &template); err != nil {
		return nil, err
	}
	return &template, nil
}

// Delete removes a container template from the application.
func (s *containerTemplateService) Delete(ctx context.Context, containerID string) error {
	path := fmt.Sprintf("/apps/%s/containers/%s", s.appID, containerID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// SetEnvironmentVariables sets environment variables for a container.
func (s *containerTemplateService) SetEnvironmentVariables(ctx context.Context, containerID string, envVars SetEnvironmentVariablesRequest) (*ContainerTemplate, error) {
	path := fmt.Sprintf("/apps/%s/containers/%s/env", s.appID, containerID)

	var template ContainerTemplate
	if err := s.client.do(ctx, http.MethodPut, path, envVars, &template); err != nil {
		return nil, err
	}
	return &template, nil
}
