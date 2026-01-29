package storage_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
	"github.com/geraldo/bunny-sdk-go/storage"
)

func TestNewClient(t *testing.T) {
	client := storage.NewClient("test-api-key")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestZoneService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if req.Header.Get("AccessKey") != "test-key" {
				t.Error("expected AccessKey header")
			}

			body := `{"Items":[{"Id":123,"Name":"test-zone","Region":"de"}],"TotalItems":1,"CurrentPage":0,"PageSize":100}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	resp, err := client.Zones().List(context.Background(), nil)

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

func TestZoneService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Id":123,"Name":"test-zone","Region":"de","StorageUsed":1024}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	zone, err := client.Zones().Get(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.ID != 123 {
		t.Errorf("expected ID 123, got %d", zone.ID)
	}
}

func TestZoneService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}

			body := `{"Id":456,"Name":"new-zone","Region":"ny","Password":"secret123"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	zone, err := client.Zones().Create(context.Background(), &storage.CreateZoneRequest{
		Name:   "new-zone",
		Region: "ny",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.Name != "new-zone" {
		t.Errorf("expected new-zone, got %s", zone.Name)
	}
}

func TestZoneService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, `{"success":true}`), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	err := client.Zones().Delete(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFileService_Upload(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if req.Header.Get("AccessKey") != "zone-pass" {
				t.Error("expected AccessKey header with zone password")
			}
			if !strings.Contains(req.URL.String(), "test-zone/files/test.txt") {
				t.Errorf("unexpected URL: %s", req.URL.String())
			}
			return testutil.NewMockResponse(201, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.Upload(context.Background(), "files/test.txt", strings.NewReader("content"), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFileService_Download(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, "file content here"), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	reader, err := fs.Download(context.Background(), "files/test.txt")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer reader.Close()

	content, _ := io.ReadAll(reader)
	if string(content) != "file content here" {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestFileService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.HasSuffix(req.URL.Path, "/") {
				t.Error("list path should end with /")
			}

			body := `[{"Guid":"file1","ObjectName":"test.txt","Length":1024,"IsDirectory":false}]`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	files, err := fs.List(context.Background(), "files")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

func TestFileService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(200, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.Delete(context.Background(), "files/test.txt")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegionBaseURL(t *testing.T) {
	tests := []struct {
		region   storage.Region
		expected string
	}{
		{storage.RegionFalkenstein, "https://storage.bunnycdn.com"},
		{storage.RegionNewYork, "https://ny.storage.bunnycdn.com"},
		{storage.RegionSingapore, "https://sg.storage.bunnycdn.com"},
		{storage.RegionLondon, "https://uk.storage.bunnycdn.com"},
	}

	for _, tt := range tests {
		t.Run(string(tt.region), func(t *testing.T) {
			got := storage.RegionBaseURL(tt.region)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
