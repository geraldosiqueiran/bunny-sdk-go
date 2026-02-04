package shield_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
	"github.com/geraldo/bunny-sdk-go/shield"
)

func TestNewClient(t *testing.T) {
	client := shield.NewClient("test-api-key")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := shield.NewClient("test-key",
		shield.WithHTTPClient(mock),
		shield.WithUserAgent("custom/1.0"),
		shield.WithBaseURL("https://custom.example.com"),
	)
	if client == nil {
		t.Fatal("NewClient with options returned nil")
	}
}

func TestClient_ServiceAccessors(t *testing.T) {
	client := shield.NewClient("test-key")

	if client.Zones() == nil {
		t.Error("Zones() returned nil")
	}
	if client.WAF() == nil {
		t.Error("WAF() returned nil")
	}
	if client.AccessLists("zone-id") == nil {
		t.Error("AccessLists() returned nil")
	}
	if client.RateLimits() == nil {
		t.Error("RateLimits() returned nil")
	}
	if client.BotDetection("zone-id") == nil {
		t.Error("BotDetection() returned nil")
	}
	if client.UploadScanning("zone-id") == nil {
		t.Error("UploadScanning() returned nil")
	}
	if client.Metrics() == nil {
		t.Error("Metrics() returned nil")
	}
	if client.EventLogs() == nil {
		t.Error("EventLogs() returned nil")
	}
	if client.DDoS() == nil {
		t.Error("DDoS() returned nil")
	}
	if client.Promo() == nil {
		t.Error("Promo() returned nil")
	}
}

// Zone Service Tests

func TestZoneService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if req.Header.Get("AccessKey") != "test-key" {
				t.Error("expected AccessKey header")
			}
			if !strings.Contains(req.URL.Path, "/shield/zones") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}

			body := `{"Items":[{"Id":"zone-123","Name":"test-zone","HostNames":["example.com"]}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Zones().List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].Name != "test-zone" {
		t.Errorf("expected test-zone, got %s", resp.Items[0].Name)
	}
}

func TestZoneService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Id":"zone-456","Name":"new-zone","HostNames":["new.example.com"]}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	zone, err := client.Zones().Create(context.Background(), &shield.CreateZoneRequest{
		Name:      "new-zone",
		HostNames: []string{"new.example.com"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.Name != "new-zone" {
		t.Errorf("expected new-zone, got %s", zone.Name)
	}
}

func TestZoneService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/shield/zone/zone-123") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":"zone-123","Name":"test-zone"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	zone, err := client.Zones().Get(context.Background(), "zone-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.ID != "zone-123" {
		t.Errorf("expected zone-123, got %s", zone.ID)
	}
}

func TestZoneService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			body := `{"Id":"zone-123","Name":"updated-zone"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	zone, err := client.Zones().Update(context.Background(), "zone-123", &shield.UpdateZoneRequest{
		Name: "updated-zone",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.Name != "updated-zone" {
		t.Errorf("expected updated-zone, got %s", zone.Name)
	}
}

func TestZoneService_GetByPullZone(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/shield/zone/pullzone/12345") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Id":"zone-123","Name":"test-zone"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	zone, err := client.Zones().GetByPullZone(context.Background(), 12345)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.ID != "zone-123" {
		t.Errorf("expected zone-123, got %s", zone.ID)
	}
}

func TestZoneService_GetPullZoneMapping(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"ShieldZoneId":"zone-123","PullZoneId":12345}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Zones().GetPullZoneMapping(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

// WAF Service Tests

func TestWAFService_ListRules(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"Id":"rule-1","Name":"SQL Injection","RuleType":"Predefined","IsActive":true}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().ListRules(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 rule, got %d", len(resp.Items))
	}
}

func TestWAFService_CreateCustomRule(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Id":"custom-rule-1","Name":"Block Admin","Pattern":"/admin/*","Action":"Block"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	rule, err := client.WAF().CreateCustomRule(context.Background(), &shield.CreateCustomRuleRequest{
		Name:    "Block Admin",
		Pattern: "/admin/*",
		Action:  "Block",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Name != "Block Admin" {
		t.Errorf("expected Block Admin, got %s", rule.Name)
	}
}

func TestWAFService_UpdateCustomRule(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			body := `{"Id":"custom-rule-1","Name":"Block Admin","IsActive":false}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	isActive := false
	rule, err := client.WAF().UpdateCustomRule(context.Background(), "custom-rule-1", &shield.UpdateCustomRuleRequest{
		IsActive: &isActive,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.IsActive {
		t.Error("expected IsActive to be false")
	}
}

func TestWAFService_ReplaceCustomRule(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			body := `{"Id":"custom-rule-1","Name":"Updated Rule","IsActive":true}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	rule, err := client.WAF().ReplaceCustomRule(context.Background(), "custom-rule-1", &shield.ReplaceCustomRuleRequest{
		Name:     "Updated Rule",
		IsActive: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Name != "Updated Rule" {
		t.Errorf("expected Updated Rule, got %s", rule.Name)
	}
}

func TestWAFService_DeleteCustomRule(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	err := client.WAF().DeleteCustomRule(context.Background(), "custom-rule-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWAFService_GetProfiles(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"Id":"profile-1","Name":"Standard","IsDefault":true}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().GetProfiles(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 profile, got %d", len(resp.Items))
	}
}

// Access List Service Tests

func TestAccessListService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/shield/zone/zone-123/access-lists") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"Allowed":[{"Type":"IP","Value":"192.168.1.1"}],"Blocked":[],"Challenged":[]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.AccessLists("zone-123").Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Allowed) != 1 {
		t.Errorf("expected 1 allowed entry, got %d", len(resp.Allowed))
	}
}

func TestAccessListService_Add(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Type":"IP","Value":"192.168.1.100","Action":"Allow"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	entry, err := client.AccessLists("zone-123").Add(context.Background(), &shield.AddAccessListEntryRequest{
		Type:   "IP",
		Value:  "192.168.1.100",
		Action: "Allow",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Value != "192.168.1.100" {
		t.Errorf("expected 192.168.1.100, got %s", entry.Value)
	}
}

func TestAccessListService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	err := client.AccessLists("zone-123").Delete(context.Background(), &shield.DeleteAccessListEntriesRequest{
		Entries: []shield.AccessListEntryIdentifier{
			{Type: "IP", Value: "192.168.1.100"},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Rate Limit Service Tests

func TestRateLimitService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"Id":"rl-1","Name":"API Rate Limit","Path":"/api/*","RequestsPerSecond":100}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.RateLimits().List(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 rate limit, got %d", len(resp.Items))
	}
}

func TestRateLimitService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Id":"rl-2","Name":"Login Rate Limit","Path":"/login","RequestsPerSecond":5}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	rateLimit, err := client.RateLimits().Create(context.Background(), &shield.CreateRateLimitRequest{
		Name:              "Login Rate Limit",
		Path:              "/login",
		RequestsPerSecond: 5,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rateLimit.Name != "Login Rate Limit" {
		t.Errorf("expected Login Rate Limit, got %s", rateLimit.Name)
	}
}

func TestRateLimitService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	err := client.RateLimits().Delete(context.Background(), "rl-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Metrics Service Tests

func TestMetricsService_GetOverview(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"TotalRequests":1000000,"BlockedRequests":5000,"AllowedRequests":995000}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetOverview(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalRequests != 1000000 {
		t.Errorf("expected 1000000, got %d", resp.TotalRequests)
	}
}

func TestMetricsService_GetOverview_WithDateRange(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "from=") {
				t.Error("expected from parameter")
			}
			if !strings.Contains(req.URL.RawQuery, "to=") {
				t.Error("expected to parameter")
			}
			body := `{"TotalRequests":500000,"BlockedRequests":2500}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetOverview(context.Background(), &shield.DateRangeOptions{
		From: "2026-02-01",
		To:   "2026-02-04",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalRequests != 500000 {
		t.Errorf("expected 500000, got %d", resp.TotalRequests)
	}
}

// Event Logs Service Tests

func TestEventLogsService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"Id":"log-1","ZoneId":"zone-123","EventType":"WAFRuleTriggered","SourceIP":"192.168.1.1","Action":"Block"}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.EventLogs().List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 log, got %d", len(resp.Items))
	}
}

func TestEventLogsService_List_WithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "zoneId=") {
				t.Error("expected zoneId parameter")
			}
			if !strings.Contains(req.URL.RawQuery, "limit=") {
				t.Error("expected limit parameter")
			}
			body := `{"Items":[],"TotalCount":0}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, err := client.EventLogs().List(context.Background(), &shield.EventLogListOptions{
		ZoneID: "zone-123",
		Limit:  50,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Bot Detection Service Tests

func TestBotDetectionService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/shield/zone/zone-123/bot-detection") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"IsEnabled":true,"DetectionLevel":"Medium","Action":"Block"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.BotDetection("zone-123").Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsEnabled {
		t.Error("expected IsEnabled to be true")
	}
}

func TestBotDetectionService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			body := `{"IsEnabled":true,"DetectionLevel":"High"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.BotDetection("zone-123").Update(context.Background(), &shield.UpdateBotDetectionRequest{
		DetectionLevel: "High",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.DetectionLevel != "High" {
		t.Errorf("expected High, got %s", resp.DetectionLevel)
	}
}

// DDoS Service Tests

func TestDDoSService_GetEnums(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Profiles":[{"Id":"low","Name":"Low"}],"Triggers":["RequestSpike"]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.DDoS().GetEnums(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(resp.Profiles))
	}
}

// Promo Service Tests

func TestPromoService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"CurrentPromos":[{"Id":"promo-1","Title":"Shield Pro Discount"}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Promo().Get(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.CurrentPromos) != 1 {
		t.Errorf("expected 1 promo, got %d", len(resp.CurrentPromos))
	}
}

// Error Handling Tests

func TestZoneService_Error(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(401, `{"Message":"unauthorized"}`), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, err := client.Zones().List(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "bunny shield api") {
		t.Errorf("expected bunny shield api error, got: %s", err.Error())
	}
}

func TestZoneService_NotFound(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"Zone not found","ErrorKey":"ZONE_NOT_FOUND"}`), nil
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, err := client.Zones().Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "ZONE_NOT_FOUND") {
		t.Errorf("expected ZONE_NOT_FOUND error, got: %s", err.Error())
	}
}

func TestNetworkError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		},
	}

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, err := client.Zones().List(context.Background())

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

	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, err := client.Zones().List(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("expected decode error, got: %s", err.Error())
	}
}

// Additional WAF Service Tests

func TestWAFService_ListCustomRules(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"Id":"cr-1","Name":"Custom Rule"}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().ListCustomRules(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1, got %d", len(resp.Items))
	}
}

func TestWAFService_GetCustomRule(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Id":"cr-1","Name":"Custom Rule"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	rule, err := client.WAF().GetCustomRule(context.Background(), "cr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "cr-1" {
		t.Errorf("expected cr-1, got %s", rule.ID)
	}
}

func TestWAFService_GetEngineConfig(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"ProfileId":"p1","IsEnabled":true}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().GetEngineConfig(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsEnabled {
		t.Error("expected IsEnabled true")
	}
}

func TestWAFService_GetEnums(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"RuleActions":["Block","Allow"]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().GetEnums(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.RuleActions) != 2 {
		t.Errorf("expected 2, got %d", len(resp.RuleActions))
	}
}

func TestWAFService_GetTriggeredRules(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"RuleId":"r1","TriggerCount":10}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().GetTriggeredRules(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1, got %d", len(resp.Items))
	}
}

func TestWAFService_SubmitTriggeredRuleReview(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"ReviewId":"rev-1","RuleId":"r1","Action":"Whitelist"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().SubmitTriggeredRuleReview(context.Background(), &shield.TriggeredRuleReviewRequest{
		RuleID: "r1", Action: "Whitelist",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ReviewID != "rev-1" {
		t.Errorf("expected rev-1, got %s", resp.ReviewID)
	}
}

func TestWAFService_GetAIRecommendation(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Recommendations":[{"RuleId":"r1","Confidence":0.95}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().GetAIRecommendation(context.Background(), "r1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Recommendations) != 1 {
		t.Errorf("expected 1, got %d", len(resp.Recommendations))
	}
}

func TestWAFService_GetPlanSegmentation(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Plans":[{"PlanId":"starter","AvailableRules":25}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.WAF().GetPlanSegmentation(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Plans) != 1 {
		t.Errorf("expected 1, got %d", len(resp.Plans))
	}
}

// Additional Access List Tests

func TestAccessListService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, ""), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	err := client.AccessLists("zone-123").Update(context.Background(), &shield.UpdateAccessListEntriesRequest{
		Updates: []shield.AccessListEntryUpdate{{Type: "IP", Value: "1.1.1.1", Action: "Block"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAccessListService_GetEnums(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Types":["IP","CIDR"],"Actions":["Allow","Block"]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.AccessLists("zone-123").GetEnums(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Types) != 2 {
		t.Errorf("expected 2, got %d", len(resp.Types))
	}
}

func TestAccessListService_UpdateConfig(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"DefaultAction":"Allow","IsEnabled":true}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.AccessLists("zone-123").UpdateConfig(context.Background(), &shield.UpdateAccessListConfigRequest{
		DefaultAction: "Allow",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.DefaultAction != "Allow" {
		t.Errorf("expected Allow, got %s", resp.DefaultAction)
	}
}

// Additional Rate Limit Tests

func TestRateLimitService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Id":"rl-1","Name":"Test","Path":"/api"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.RateLimits().Get(context.Background(), "rl-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "rl-1" {
		t.Errorf("expected rl-1, got %s", resp.ID)
	}
}

func TestRateLimitService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}
			body := `{"Id":"rl-1","Name":"Updated"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.RateLimits().Update(context.Background(), "rl-1", &shield.UpdateRateLimitRequest{
		Name: "Updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != "Updated" {
		t.Errorf("expected Updated, got %s", resp.Name)
	}
}

// Upload Scanning Tests

func TestUploadScanningService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"IsEnabled":true,"ScanLevel":"Full"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.UploadScanning("zone-123").Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsEnabled {
		t.Error("expected IsEnabled true")
	}
}

func TestUploadScanningService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"IsEnabled":true,"ScanLevel":"Basic"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.UploadScanning("zone-123").Update(context.Background(), &shield.UpdateUploadScanningRequest{
		ScanLevel: "Basic",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ScanLevel != "Basic" {
		t.Errorf("expected Basic, got %s", resp.ScanLevel)
	}
}

// Additional Metrics Tests

func TestMetricsService_GetOverviewDetailed(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Zones":[{"ZoneId":"z1","TotalRequests":1000}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetOverviewDetailed(context.Background(), &shield.MetricsDetailedOptions{
		From: "2026-01-01", ZoneID: "z1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Zones) != 1 {
		t.Errorf("expected 1, got %d", len(resp.Zones))
	}
}

func TestMetricsService_GetWAFRuleMetrics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"RuleId":"r1","TriggerCount":100}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetWAFRuleMetrics(context.Background(), "r1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TriggerCount != 100 {
		t.Errorf("expected 100, got %d", resp.TriggerCount)
	}
}

func TestMetricsService_GetRateLimitMetrics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"RuleId":"rl-1","BlockedRequests":50}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetRateLimitMetrics(context.Background(), "rl-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.BlockedRequests != 50 {
		t.Errorf("expected 50, got %d", resp.BlockedRequests)
	}
}

func TestMetricsService_GetAllRateLimitMetrics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Items":[{"RuleId":"rl-1"}],"TotalCount":1}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetAllRateLimitMetrics(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1, got %d", len(resp.Items))
	}
}

func TestMetricsService_GetBotDetectionMetrics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"ZoneId":"z1","TotalBotRequests":500}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetBotDetectionMetrics(context.Background(), "z1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalBotRequests != 500 {
		t.Errorf("expected 500, got %d", resp.TotalBotRequests)
	}
}

func TestMetricsService_GetUploadScanningMetrics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"ZoneId":"z1","TotalScannedFiles":100}`
			return testutil.NewMockResponse(200, body), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	resp, err := client.Metrics().GetUploadScanningMetrics(context.Background(), "z1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalScannedFiles != 100 {
		t.Errorf("expected 100, got %d", resp.TotalScannedFiles)
	}
}

// Error with Field test

func TestAPIError_WithField(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"Invalid","ErrorKey":"ERR","Field":"name"}`), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, err := client.Zones().Create(context.Background(), &shield.CreateZoneRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "field: name") {
		t.Errorf("expected field in error, got: %s", err.Error())
	}
}

// Event log with all options

func TestEventLogsService_List_AllOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			q := req.URL.RawQuery
			if !strings.Contains(q, "from=") || !strings.Contains(q, "to=") || !strings.Contains(q, "offset=") {
				t.Error("missing query params")
			}
			return testutil.NewMockResponse(200, `{"Items":[],"TotalCount":0}`), nil
		},
	}
	client := shield.NewClient("test-key", shield.WithHTTPClient(mock))
	_, _ = client.EventLogs().List(context.Background(), &shield.EventLogListOptions{
		ZoneID: "z1", From: "2026-01-01", To: "2026-02-01", Limit: 10, Offset: 5,
	})
}

// Error tests for full coverage

func errMock() *testutil.MockHTTPClient {
	return &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(500, `{"Message":"error"}`), nil
		},
	}
}

func TestZoneService_CreateError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Zones().Create(context.Background(), &shield.CreateZoneRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestZoneService_UpdateError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Zones().Update(context.Background(), "z", &shield.UpdateZoneRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestZoneService_GetByPullZoneError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Zones().GetByPullZone(context.Background(), 123)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestZoneService_GetPullZoneMappingError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Zones().GetPullZoneMapping(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_ListRulesError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().ListRules(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_ListCustomRulesError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().ListCustomRules(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_CreateCustomRuleError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().CreateCustomRule(context.Background(), &shield.CreateCustomRuleRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetCustomRuleError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetCustomRule(context.Background(), "r")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_UpdateCustomRuleError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().UpdateCustomRule(context.Background(), "r", &shield.UpdateCustomRuleRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_ReplaceCustomRuleError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().ReplaceCustomRule(context.Background(), "r", &shield.ReplaceCustomRuleRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetProfilesError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetProfiles(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetEngineConfigError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetEngineConfig(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetEnumsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetEnums(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetTriggeredRulesError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetTriggeredRules(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_SubmitTriggeredRuleReviewError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().SubmitTriggeredRuleReview(context.Background(), &shield.TriggeredRuleReviewRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetAIRecommendationError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetAIRecommendation(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWAFService_GetPlanSegmentationError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.WAF().GetPlanSegmentation(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccessListService_GetError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.AccessLists("z").Get(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccessListService_AddError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.AccessLists("z").Add(context.Background(), &shield.AddAccessListEntryRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccessListService_GetEnumsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.AccessLists("z").GetEnums(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccessListService_UpdateConfigError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.AccessLists("z").UpdateConfig(context.Background(), &shield.UpdateAccessListConfigRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRateLimitService_ListError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.RateLimits().List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRateLimitService_CreateError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.RateLimits().Create(context.Background(), &shield.CreateRateLimitRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRateLimitService_GetError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.RateLimits().Get(context.Background(), "r")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRateLimitService_UpdateError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.RateLimits().Update(context.Background(), "r", &shield.UpdateRateLimitRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBotDetectionService_GetError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.BotDetection("z").Get(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBotDetectionService_UpdateError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.BotDetection("z").Update(context.Background(), &shield.UpdateBotDetectionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUploadScanningService_GetError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.UploadScanning("z").Get(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUploadScanningService_UpdateError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.UploadScanning("z").Update(context.Background(), &shield.UpdateUploadScanningRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetOverviewError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetOverview(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetOverviewDetailedError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetOverviewDetailed(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetWAFRuleMetricsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetWAFRuleMetrics(context.Background(), "r", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetRateLimitMetricsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetRateLimitMetrics(context.Background(), "r", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetAllRateLimitMetricsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetAllRateLimitMetrics(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetBotDetectionMetricsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetBotDetectionMetrics(context.Background(), "z", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsService_GetUploadScanningMetricsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Metrics().GetUploadScanningMetrics(context.Background(), "z", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEventLogsService_ListError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.EventLogs().List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDDoSService_GetEnumsError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.DDoS().GetEnums(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPromoService_GetError(t *testing.T) {
	client := shield.NewClient("k", shield.WithHTTPClient(errMock()))
	_, err := client.Promo().Get(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// Edge case tests for query builders

func TestMetricsService_EmptyDateRange(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, `{}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, _ = client.Metrics().GetOverview(context.Background(), &shield.DateRangeOptions{})
}

func TestMetricsService_OnlyFromDate(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "from=") {
				t.Error("expected from")
			}
			return testutil.NewMockResponse(200, `{}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, _ = client.Metrics().GetOverview(context.Background(), &shield.DateRangeOptions{From: "2026-01-01"})
}

func TestMetricsService_DetailedOnlyZoneID(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "zoneId=") {
				t.Error("expected zoneId")
			}
			return testutil.NewMockResponse(200, `{"Zones":[]}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, _ = client.Metrics().GetOverviewDetailed(context.Background(), &shield.MetricsDetailedOptions{ZoneID: "z1"})
}

func TestEventLogs_EmptyOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, `{"Items":[],"TotalCount":0}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, _ = client.EventLogs().List(context.Background(), &shield.EventLogListOptions{})
}

func TestMetricsService_DetailedEmpty(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, `{"Zones":[]}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, _ = client.Metrics().GetOverviewDetailed(context.Background(), &shield.MetricsDetailedOptions{})
}

func TestMetricsService_DetailedTo(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "to=") {
				t.Error("expected to")
			}
			return testutil.NewMockResponse(200, `{"Zones":[]}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, _ = client.Metrics().GetOverviewDetailed(context.Background(), &shield.MetricsDetailedOptions{To: "2026-02-01"})
}


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
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, err := client.Zones().List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestErrorResponse_EmptyMessage(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Valid JSON but empty Message field
			return testutil.NewMockResponse(400, `{"ErrorKey":"ERR"}`), nil
		},
	}
	client := shield.NewClient("k", shield.WithHTTPClient(mock))
	_, err := client.Zones().List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
