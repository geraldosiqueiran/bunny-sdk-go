package containers

import (
	"context"
	"net/http"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
)

func TestRegistryService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/registries" {
				t.Errorf("expected /registries, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{
				"items": [{"id": 123, "displayName": "Docker Hub"}],
				"meta": {"totalItems": 1}
			}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Registries().List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestRegistryService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/registries/123" {
				t.Errorf("expected /registries/123, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": 123, "displayName": "My Registry"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	registry, err := client.Registries().Get(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if registry.ID != 123 {
		t.Errorf("expected ID 123, got %d", registry.ID)
	}
}

func TestRegistryService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"id": 456, "status": "Saved"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Registries().Create(context.Background(), &CreateRegistryRequest{
		DisplayName: "New Registry",
		Type:        RegistryTypeDockerHub,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != 456 {
		t.Errorf("expected ID 456, got %d", resp.ID)
	}
}

func TestRegistryService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if req.URL.Path != "/mc/registries/123" {
				t.Errorf("expected /registries/123, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": 123, "status": "Saved"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Registries().Update(context.Background(), 123, &UpdateRegistryRequest{
		DisplayName: "Updated Registry",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != RegistryStatusSaved {
		t.Errorf("expected Saved status, got %s", resp.Status)
	}
}

func TestRegistryService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"status": "Removed"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Registries().Delete(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != RegistryDeleteStatusRemoved {
		t.Errorf("expected Removed, got %s", resp.Status)
	}
}

func TestRegistryService_ListImages(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if req.URL.Path != "/mc/registries/images" {
				t.Errorf("expected /registries/images, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `[{"id": "nginx", "namespace": "library"}]`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	images, err := client.Registries().ListImages(context.Background(), &ListImagesRequest{
		RegistryID: "dockerhub",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(images) != 1 {
		t.Errorf("expected 1 image, got %d", len(images))
	}
}

func TestRegistryService_ListTags(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/registries/tags" {
				t.Errorf("expected /registries/tags, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `[{"name": "latest"}, {"name": "1.0"}]`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	tags, err := client.Registries().ListTags(context.Background(), &ListTagsRequest{
		RegistryID:     "dockerhub",
		ImageName:      "nginx",
		ImageNamespace: "library",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestRegistryService_GetDigest(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/registries/digest" {
				t.Errorf("expected /registries/digest, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"digest": "sha256:abc123"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	digest, err := client.Registries().GetDigest(context.Background(), &GetDigestRequest{
		RegistryID:     "dockerhub",
		ImageName:      "nginx",
		ImageNamespace: "library",
		Tag:            "latest",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if digest.Digest != "sha256:abc123" {
		t.Errorf("expected sha256:abc123, got %s", digest.Digest)
	}
}

func TestRegistryService_GetConfigSuggestions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/registries/config-suggestions" {
				t.Errorf("expected /registries/config-suggestions, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"appName": "nginx", "description": "Web server"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	suggestions, err := client.Registries().GetConfigSuggestions(context.Background(), &GetConfigSuggestionsRequest{
		RegistryID:     "dockerhub",
		ImageName:      "nginx",
		ImageNamespace: "library",
		Tag:            "latest",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suggestions.AppName != "nginx" {
		t.Errorf("expected nginx, got %s", suggestions.AppName)
	}
}

func TestRegistryService_SearchPublicImages(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/registries/public-images/search" {
				t.Errorf("expected /registries/public-images/search, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `[{"id": "nginx"}, {"id": "node"}]`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	images, err := client.Registries().SearchPublicImages(context.Background(), &SearchPublicImagesRequest{
		RegistryID: "dockerhub",
		Prefix:     "n",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(images) != 2 {
		t.Errorf("expected 2 images, got %d", len(images))
	}
}
