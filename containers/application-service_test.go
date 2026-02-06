package containers

import (
	"context"
	"net/http"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
)

func TestApplicationService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if req.URL.Path != "/mc/apps" {
				t.Errorf("expected /mc/apps, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{
				"items": [{"id": "app1", "name": "Test App", "status": "Active"}],
				"meta": {"totalItems": 1}
			}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Applications().List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].ID != "app1" {
		t.Errorf("expected app1, got %s", resp.Items[0].ID)
	}
}

func TestApplicationService_ListWithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("limit") != "10" {
				t.Errorf("expected limit=10, got %s", req.URL.Query().Get("limit"))
			}
			return testutil.NewMockResponse(200, `{"items":[],"meta":{"totalItems":0}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	_, err := client.Applications().List(context.Background(), &ListOptions{Limit: 10})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123" {
				t.Errorf("expected /apps/app123, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{
				"id": "app123",
				"name": "My App",
				"status": "Active",
				"runtimeType": "Shared"
			}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	app, err := client.Applications().Get(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.ID != "app123" {
		t.Errorf("expected app123, got %s", app.ID)
	}
}

func TestApplicationService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			return testutil.NewMockResponse(201, `{"id": "new-app-id"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Applications().Create(context.Background(), &CreateApplicationRequest{
		Name:        "New App",
		RuntimeType: RuntimeTypeShared,
		AutoScaling: AutoScaling{Min: 1, Max: 5},
		RegionSettings: CreateRegionSettingsRequest{
			AllowedRegionIds: []string{"DE"},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "new-app-id" {
		t.Errorf("expected new-app-id, got %s", resp.ID)
	}
}

func TestApplicationService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if req.URL.Path != "/mc/apps/app123" {
				t.Errorf("expected /apps/app123, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{"id": "app123"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Applications().Update(context.Background(), "app123", &UpdateApplicationRequest{
		Name:        "Updated App",
		RuntimeType: RuntimeTypeShared,
		AutoScaling: AutoScaling{Min: 2, Max: 10},
		RegionSettings: CreateRegionSettingsRequest{
			AllowedRegionIds: []string{"DE", "US"},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "app123" {
		t.Errorf("expected app123, got %s", resp.ID)
	}
}

func TestApplicationService_Patch(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"id": "app123"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	resp, err := client.Applications().Patch(context.Background(), "app123", &PatchApplicationRequest{
		Name: "Patched App",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "app123" {
		t.Errorf("expected app123, got %s", resp.ID)
	}
}

func TestApplicationService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Applications().Delete(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationService_Deploy(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/deploy" {
				t.Errorf("expected /apps/app123/deploy, got %s", req.URL.Path)
			}
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Applications().Deploy(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationService_Undeploy(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/undeploy" {
				t.Errorf("expected /apps/app123/undeploy, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Applications().Undeploy(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationService_Restart(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/restart" {
				t.Errorf("expected /apps/app123/restart, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, ``), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	err := client.Applications().Restart(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationService_GetOverview(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/mc/apps/app123/overview" {
				t.Errorf("expected /apps/app123/overview, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, `{
				"status": "Active",
				"desiredInstances": 3,
				"activeInstances": {"value": 3}
			}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	overview, err := client.Applications().GetOverview(context.Background(), "app123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.Status != ApplicationStatusActive {
		t.Errorf("expected Active status, got %s", overview.Status)
	}
}

func TestApplicationService_GetStatistics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("fromDate") != "2026-01-01" {
				t.Errorf("expected fromDate=2026-01-01")
			}
			if req.URL.Query().Get("granularity") != "Hourly" {
				t.Errorf("expected granularity=Hourly")
			}
			return testutil.NewMockResponse(200, `{"latencyChart":{}}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	_, err := client.Applications().GetStatistics(context.Background(), "app123", &StatisticsOptions{
		FromDate:    "2026-01-01",
		Granularity: StatisticsGranularityHourly,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationService_ErrorHandling(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"Not Found"}`), nil
		},
	}

	client := NewClient("key", WithHTTPClient(mock))
	_, err := client.Applications().Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}
