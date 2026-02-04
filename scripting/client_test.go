package scripting_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
	"github.com/geraldo/bunny-sdk-go/scripting"
)

func TestNewClient(t *testing.T) {
	client := scripting.NewClient("test-api-key")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := scripting.NewClient("test-key",
		scripting.WithHTTPClient(mock),
		scripting.WithUserAgent("custom/1.0"),
		scripting.WithBaseURL("https://custom.example.com"),
	)
	if client == nil {
		t.Fatal("NewClient with options returned nil")
	}
}

func TestClient_ServiceAccessors(t *testing.T) {
	client := scripting.NewClient("test-key")

	if client.Scripts() == nil {
		t.Error("Scripts() returned nil")
	}
	if client.Code(123) == nil {
		t.Error("Code() returned nil")
	}
	if client.Releases(123) == nil {
		t.Error("Releases() returned nil")
	}
	if client.Secrets(123) == nil {
		t.Error("Secrets() returned nil")
	}
	if client.Variables(123) == nil {
		t.Error("Variables() returned nil")
	}
}

// Script Service Tests

func TestScriptService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if req.Header.Get("AccessKey") != "test-key" {
				t.Error("expected AccessKey header")
			}
			if !strings.Contains(req.URL.Path, "/compute/script") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}

			body := `{"Items":[{"Id":12345,"Name":"test-script","ScriptType":"CDN"}],"CurrentPage":1,"TotalItems":1,"HasMoreItems":false}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	resp, err := client.Scripts().List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].ID != 12345 {
		t.Errorf("expected ID 12345, got %d", resp.Items[0].ID)
	}
}

func TestScriptService_List_WithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			q := req.URL.RawQuery
			if !strings.Contains(q, "page=2") {
				t.Error("expected page parameter")
			}
			if !strings.Contains(q, "perPage=10") {
				t.Error("expected perPage parameter")
			}
			if !strings.Contains(q, "search=test") {
				t.Error("expected search parameter")
			}
			if !strings.Contains(q, "includeLinkedPullzones=true") {
				t.Error("expected includeLinkedPullzones parameter")
			}
			body := `{"Items":[],"CurrentPage":2,"TotalItems":0,"HasMoreItems":false}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().List(context.Background(), &scripting.ScriptListOptions{
		Page:                   2,
		PerPage:                10,
		Search:                 "test",
		IncludeLinkedPullzones: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptService_List_WithIntegrationID(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "integrationId=999") {
				t.Error("expected integrationId parameter")
			}
			return testutil.NewMockResponse(200, `{"Items":[],"CurrentPage":1,"TotalItems":0,"HasMoreItems":false}`), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	integrationID := int64(999)
	_, err := client.Scripts().List(context.Background(), &scripting.ScriptListOptions{
		IntegrationID: &integrationID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Id":67890,"Name":"new-script","ScriptType":"CDN","DefaultHostname":"new-script.b-cdn.net"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	script, err := client.Scripts().Create(context.Background(), &scripting.CreateScriptRequest{
		Name:                 "new-script",
		ScriptType:           scripting.ScriptTypeCDN,
		CreateLinkedPullZone: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if script.ID != 67890 {
		t.Errorf("expected ID 67890, got %d", script.ID)
	}
}

func TestScriptService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":12345,"Name":"test-script","ScriptType":"CDN"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	script, err := client.Scripts().Get(context.Background(), 12345)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if script.ID != 12345 {
		t.Errorf("expected ID 12345, got %d", script.ID)
	}
}

func TestScriptService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":12345,"Name":"updated-script","ScriptType":"Middleware"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	script, err := client.Scripts().Update(context.Background(), 12345, &scripting.UpdateScriptRequest{
		Name:       "updated-script",
		ScriptType: scripting.ScriptTypeMiddleware,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if script.ScriptType != scripting.ScriptTypeMiddleware {
		t.Errorf("expected Middleware, got %s", script.ScriptType)
	}
}

func TestScriptService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Scripts().Delete(context.Background(), 12345, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptService_Delete_WithLinkedPullZones(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "deleteLinkedPullZones=true") {
				t.Error("expected deleteLinkedPullZones parameter")
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Scripts().Delete(context.Background(), 12345, true)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptService_GetStatistics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345/statistics") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"TotalRequestsServed":1500000,"TotalCpuUsed":45000.5,"TotalMonthlyCost":12.50}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	stats, err := client.Scripts().GetStatistics(context.Background(), 12345, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.TotalRequestsServed != 1500000 {
		t.Errorf("expected 1500000, got %d", stats.TotalRequestsServed)
	}
}

func TestScriptService_GetStatistics_WithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			q := req.URL.RawQuery
			if !strings.Contains(q, "dateFrom=2026-01-01") {
				t.Error("expected dateFrom parameter")
			}
			if !strings.Contains(q, "dateTo=2026-02-01") {
				t.Error("expected dateTo parameter")
			}
			if !strings.Contains(q, "loadLatest=true") {
				t.Error("expected loadLatest parameter")
			}
			if !strings.Contains(q, "hourly=true") {
				t.Error("expected hourly parameter")
			}
			body := `{"TotalRequestsServed":500000}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().GetStatistics(context.Background(), 12345, &scripting.StatisticsOptions{
		DateFrom:   "2026-01-01",
		DateTo:     "2026-02-01",
		LoadLatest: true,
		Hourly:     true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptService_RotateDeploymentKey(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/deploymentKey/rotate") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Scripts().RotateDeploymentKey(context.Background(), 12345)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Code Service Tests

func TestCodeService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345/code") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Code":"export default { fetch() { return new Response('Hello'); } }","LastModified":"2026-02-04T10:30:00Z"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	code, err := client.Code(12345).Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code.Code == nil || *code.Code == "" {
		t.Error("expected code content")
	}
}

func TestCodeService_Set(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/code") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Code(12345).Set(context.Background(), &scripting.UpdateCodeRequest{
		Code: "export default { fetch() { return new Response('Updated'); } }",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Release Service Tests

func TestReleaseService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345/releases") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Items":[{"Id":67890,"Status":"Live","Note":"Initial release"}],"CurrentPage":1,"TotalItems":1,"HasMoreItems":false}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	resp, err := client.Releases(12345).List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestReleaseService_List_WithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			q := req.URL.RawQuery
			if !strings.Contains(q, "page=2") {
				t.Error("expected page parameter")
			}
			if !strings.Contains(q, "perPage=50") {
				t.Error("expected perPage parameter")
			}
			body := `{"Items":[],"CurrentPage":2,"TotalItems":0,"HasMoreItems":false}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Releases(12345).List(context.Background(), &scripting.ReleaseListOptions{
		Page:    2,
		PerPage: 50,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReleaseService_GetActive(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345/releases/active") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":67890,"Status":"Live","Code":"export default {}"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	release, err := client.Releases(12345).GetActive(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if release.ID != 67890 {
		t.Errorf("expected ID 67890, got %d", release.ID)
	}
}

func TestReleaseService_Publish(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/publish") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Releases(12345).Publish(context.Background(), &scripting.PublishReleaseRequest{
		Note: "New release",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReleaseService_PublishByUUID(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/publish/abc-123-def") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Releases(12345).PublishByUUID(context.Background(), "abc-123-def", &scripting.PublishReleaseRequest{
		Note: "Rollback release",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Secret Service Tests

func TestSecretService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345/secrets") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Secrets":[{"Id":111,"Name":"API_KEY","LastModified":"2026-02-04T10:30:00Z"}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	resp, err := client.Secrets(12345).List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(resp.Secrets))
	}
}

func TestSecretService_Add(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Id":222,"Name":"NEW_SECRET","LastModified":"2026-02-04T10:30:00Z"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	secret, err := client.Secrets(12345).Add(context.Background(), &scripting.AddSecretRequest{
		Name:   "NEW_SECRET",
		Secret: "secret-value",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secret.ID != 222 {
		t.Errorf("expected ID 222, got %d", secret.ID)
	}
}

func TestSecretService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/secrets/111") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":111,"Name":"API_KEY","LastModified":"2026-02-04T11:00:00Z"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	secret, err := client.Secrets(12345).Update(context.Background(), 111, &scripting.UpdateSecretRequest{
		Secret: "new-secret-value",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secret.ID != 111 {
		t.Errorf("expected ID 111, got %d", secret.ID)
	}
}

func TestSecretService_Upsert(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			body := `{"Id":333,"Name":"UPSERTED_SECRET","LastModified":"2026-02-04T10:30:00Z"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	secret, err := client.Secrets(12345).Upsert(context.Background(), &scripting.UpsertSecretRequest{
		Name:   "UPSERTED_SECRET",
		Secret: "secret-value",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secret.ID != 333 {
		t.Errorf("expected ID 333, got %d", secret.ID)
	}
}

func TestSecretService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/secrets/111") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Secrets(12345).Delete(context.Background(), 111)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Variable Service Tests

func TestVariableService_Add(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/variables/add") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":444,"Name":"LOG_LEVEL","Required":false,"DefaultValue":"info"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	variable, err := client.Variables(12345).Add(context.Background(), &scripting.AddVariableRequest{
		Name:         "LOG_LEVEL",
		Required:     false,
		DefaultValue: "info",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if variable.ID != 444 {
		t.Errorf("expected ID 444, got %d", variable.ID)
	}
}

func TestVariableService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/compute/script/12345/variables/444") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":444,"Name":"LOG_LEVEL","Required":false,"DefaultValue":"info"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	variable, err := client.Variables(12345).Get(context.Background(), 444)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if variable.ID != 444 {
		t.Errorf("expected ID 444, got %d", variable.ID)
	}
}

func TestVariableService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/variables/444") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":444,"Name":"LOG_LEVEL","Required":true,"DefaultValue":"debug"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	required := true
	variable, err := client.Variables(12345).Update(context.Background(), 444, &scripting.UpdateVariableRequest{
		DefaultValue: "debug",
		Required:     &required,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !variable.Required {
		t.Error("expected Required to be true")
	}
}

func TestVariableService_Upsert(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/variables") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":555,"Name":"UPSERTED_VAR","Required":false,"DefaultValue":"value"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	variable, err := client.Variables(12345).Upsert(context.Background(), &scripting.UpsertVariableRequest{
		Name:         "UPSERTED_VAR",
		DefaultValue: "value",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if variable.ID != 555 {
		t.Errorf("expected ID 555, got %d", variable.ID)
	}
}

func TestVariableService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/compute/script/12345/variables/444") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	err := client.Variables(12345).Delete(context.Background(), 444)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Error Handling Tests

func TestScriptService_Error(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(401, `{"Message":"unauthorized"}`), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "bunny scripting api") {
		t.Errorf("expected bunny scripting api error, got: %s", err.Error())
	}
}

func TestScriptService_NotFound(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"Script not found","ErrorKey":"SCRIPT_NOT_FOUND"}`), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().Get(context.Background(), 99999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "SCRIPT_NOT_FOUND") {
		t.Errorf("expected SCRIPT_NOT_FOUND error, got: %s", err.Error())
	}
}

func TestNetworkError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInvalidJSONResponse(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, "not valid json {{{"), nil
		},
	}

	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("expected decode error, got: %s", err.Error())
	}
}

func TestAPIError_WithField(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"Invalid","ErrorKey":"ERR","Field":"name"}`), nil
		},
	}
	client := scripting.NewClient("test-key", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().Create(context.Background(), &scripting.CreateScriptRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "field: name") {
		t.Errorf("expected field in error, got: %s", err.Error())
	}
}

// Error tests for service methods

func errMock() *testutil.MockHTTPClient {
	return &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(500, `{"Message":"error"}`), nil
		},
	}
}

func TestScriptService_CreateError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Scripts().Create(context.Background(), &scripting.CreateScriptRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScriptService_UpdateError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Scripts().Update(context.Background(), 1, &scripting.UpdateScriptRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScriptService_GetStatisticsError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Scripts().GetStatistics(context.Background(), 1, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCodeService_GetError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Code(1).Get(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReleaseService_ListError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Releases(1).List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReleaseService_GetActiveError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Releases(1).GetActive(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretService_ListError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Secrets(1).List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretService_AddError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Secrets(1).Add(context.Background(), &scripting.AddSecretRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretService_UpdateError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Secrets(1).Update(context.Background(), 1, &scripting.UpdateSecretRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSecretService_UpsertError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Secrets(1).Upsert(context.Background(), &scripting.UpsertSecretRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVariableService_AddError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Variables(1).Add(context.Background(), &scripting.AddVariableRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVariableService_GetError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Variables(1).Get(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVariableService_UpdateError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Variables(1).Update(context.Background(), 1, &scripting.UpdateVariableRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVariableService_UpsertError(t *testing.T) {
	client := scripting.NewClient("k", scripting.WithHTTPClient(errMock()))
	_, err := client.Variables(1).Upsert(context.Background(), &scripting.UpsertVariableRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// Additional edge case tests

type badReader struct{}

func (badReader) Read(p []byte) (n int, err error) { return 0, context.Canceled }
func (badReader) Close() error                     { return nil }

func TestReadBodyError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       badReader{},
				Header:     make(http.Header),
			}, nil
		},
	}
	client := scripting.NewClient("k", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestErrorResponse_EmptyMessage(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"ErrorKey":"ERR"}`), nil
		},
	}
	client := scripting.NewClient("k", scripting.WithHTTPClient(mock))
	_, err := client.Scripts().List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
