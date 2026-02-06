package containers

import (
	"context"
	"net/http"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
)

// =============================================================================
// Container Template Service Tests
// =============================================================================

func TestContainerTemplateService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/containers/container456" {
				t.Errorf("expected /apps/app123/containers/container456, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "container456", "name": "web"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	template, err := client.ContainerTemplates("app123").Get(context.Background(), "container456")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.ID != "container456" {
		t.Errorf("expected container456, got %s", template.ID)
	}
}

func TestContainerTemplateService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if req.URL.Path != "/mc/apps/app123/containers" {
				t.Errorf("expected /apps/app123/containers, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "new-container", "name": "api"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	template, err := client.ContainerTemplates("app123").Create(context.Background(), &CreateContainerTemplateRequest{
		Name:           "api",
		ImageName:      "node",
		ImageNamespace: "library",
		ImageTag:       "18",
		ImageRegistryID: "dockerhub",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.Name != "api" {
		t.Errorf("expected api, got %s", template.Name)
	}
}

func TestContainerTemplateService_Patch(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"id": "container456", "name": "updated"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	template, err := client.ContainerTemplates("app123").Patch(context.Background(), "container456", &PatchContainerTemplateRequest{
		Name: "updated",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.Name != "updated" {
		t.Errorf("expected updated, got %s", template.Name)
	}
}

func TestContainerTemplateService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.ContainerTemplates("app123").Delete(context.Background(), "container456")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContainerTemplateService_SetEnvironmentVariables(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if req.URL.Path != "/mc/apps/app123/containers/container456/env" {
				t.Errorf("expected /apps/app123/containers/container456/env, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "container456"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	_, err := client.ContainerTemplates("app123").SetEnvironmentVariables(context.Background(), "container456", SetEnvironmentVariablesRequest{
		"NODE_ENV": "production",
		"PORT":     "3000",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// Endpoint Service Tests
// =============================================================================

func TestEndpointService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/endpoints" {
				t.Errorf("expected /apps/app123/endpoints, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"items": [{"id": "ep1"}], "meta": {"totalItems": 1}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Endpoints("app123").List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestEndpointService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/containers/container456/endpoints" {
				t.Errorf("expected /apps/app123/containers/container456/endpoints, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "new-endpoint"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Endpoints("app123").Create(context.Background(), "container456", &EndpointRequest{
		DisplayName: "web",
		CDN: &CDNEndpointConfig{
			IsSslEnabled: true,
			PortMappings: []PortMapping{{ContainerPort: 80}},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "new-endpoint" {
		t.Errorf("expected new-endpoint, got %s", resp.ID)
	}
}

func TestEndpointService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Endpoints("app123").Update(context.Background(), "ep1", &EndpointRequest{
		DisplayName: "updated",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEndpointService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Endpoints("app123").Delete(context.Background(), "ep1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// Autoscaling Service Tests
// =============================================================================

func TestAutoscalingService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/autoscaling" {
				t.Errorf("expected /apps/app123/autoscaling, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"min": 1, "max": 5}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	settings, err := client.Autoscaling("app123").Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.Min != 1 || settings.Max != 5 {
		t.Errorf("expected min=1, max=5, got min=%d, max=%d", settings.Min, settings.Max)
	}
}

func TestAutoscalingService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Autoscaling("app123").Update(context.Background(), &AutoScaling{Min: 2, Max: 10})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// Region Service Tests
// =============================================================================

func TestRegionService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/regions" {
				t.Errorf("expected /regions, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"items": [{"id": "DE", "name": "Germany"}], "meta": {"totalItems": 1}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Regions().List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestRegionService_GetOptimal(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/regions/optimal" {
				t.Errorf("expected /regions/optimal, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"region": {"id": "DE", "name": "Germany"}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Regions().GetOptimal(context.Background(), "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Region == nil || resp.Region.ID != "DE" {
		t.Errorf("expected DE region")
	}
}

// =============================================================================
// Region Settings Service Tests
// =============================================================================

func TestRegionSettingsService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/region-settings" {
				t.Errorf("expected /apps/app123/region-settings, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"allowedRegionIds": ["DE", "US"]}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	settings, err := client.RegionSettings("app123").Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings.AllowedRegionIds) != 2 {
		t.Errorf("expected 2 allowed regions, got %d", len(settings.AllowedRegionIds))
	}
}

func TestRegionSettingsService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.RegionSettings("app123").Update(context.Background(), &UpdateRegionSettingsRequest{
		AllowedRegionIds: []string{"DE", "US", "UK"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// Limits Service Tests
// =============================================================================

func TestLimitsService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/limits" {
				t.Errorf("expected /limits, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"maxNumberOfApplications": 10}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	limits, err := client.Limits().Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limits.MaxNumberOfApplications != 10 {
		t.Errorf("expected 10, got %d", limits.MaxNumberOfApplications)
	}
}

// =============================================================================
// Node Service Tests
// =============================================================================

func TestNodeService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/nodes" {
				t.Errorf("expected /nodes, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"items": ["node1", "node2"], "meta": {"totalItems": 2}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Nodes().List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(resp.Items))
	}
}

// =============================================================================
// Pod Service Tests
// =============================================================================

func TestPodService_Recreate(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if req.URL.Path != "/mc/apps/app123/pods/pod456/recreate" {
				t.Errorf("expected /apps/app123/pods/pod456/recreate, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Pods("app123").Recreate(context.Background(), "pod456")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// Volume Service Tests
// =============================================================================

func TestVolumeService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/volumes" {
				t.Errorf("expected /apps/app123/volumes, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"items": [{"id": "vol1", "name": "data", "size": 10}], "meta": {"totalItems": 1}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Volumes("app123").List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestVolumeService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"name": "updated", "size": 20}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Volumes("app123").Update(context.Background(), "vol1", &UpdateVolumeRequest{Size: 20})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 20 {
		t.Errorf("expected size 20, got %d", resp.Size)
	}
}

func TestVolumeService_Detach(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/volumes/vol1/detach" {
				t.Errorf("expected /apps/app123/volumes/vol1/detach, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"name": "data"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Volumes("app123").Detach(context.Background(), "vol1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != "data" {
		t.Errorf("expected data, got %s", resp.Name)
	}
}

func TestVolumeService_DeleteInstance(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/volumes/vol1/instances/inst1" {
				t.Errorf("expected /apps/app123/volumes/vol1/instances/inst1, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "inst1"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Volumes("app123").DeleteInstance(context.Background(), "vol1", "inst1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "inst1" {
		t.Errorf("expected inst1, got %s", resp.ID)
	}
}

func TestVolumeService_DeleteAll(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"ids": ["inst1", "inst2"]}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Volumes("app123").DeleteAll(context.Background(), "vol1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.IDs) != 2 {
		t.Errorf("expected 2 IDs, got %d", len(resp.IDs))
	}
}

// =============================================================================
// Log Forwarding Service Tests
// =============================================================================

func TestLogForwardingService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/log/forwarding" {
				t.Errorf("expected /log/forwarding, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"items": [{"id": "lf1", "app": "app123"}]}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.LogForwarding().List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestLogForwardingService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/log/forwarding/app123" {
				t.Errorf("expected /log/forwarding/app123, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "lf1", "app": "app123", "type": "SyslogUdp"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	config, err := client.LogForwarding().Get(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.Type != LogForwardingTypeSyslogUDP {
		t.Errorf("expected SyslogUdp, got %s", config.Type)
	}
}

func TestLogForwardingService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"id": "new-lf", "app": "app123"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	config, err := client.LogForwarding().Create(context.Background(), &CreateLogForwardingRequest{
		App:      "app123",
		Type:     LogForwardingTypeSyslogUDP,
		Endpoint: "logs.example.com",
		Port:     514,
		Format:   LogForwardingFormatRfc5424,
		Enabled:  true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.App != "app123" {
		t.Errorf("expected app123, got %s", config.App)
	}
}

func TestLogForwardingService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"id": "lf1", "enabled": false}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	_, err := client.LogForwarding().Update(context.Background(), "app123", &UpdateLogForwardingRequest{
		App:      "app123",
		Type:     LogForwardingTypeSyslogUDP,
		Endpoint: "logs.example.com",
		Port:     514,
		Format:   LogForwardingFormatRfc5424,
		Enabled:  false,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogForwardingService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.LogForwarding().Delete(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
