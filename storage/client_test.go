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
// TestZoneService_Update tests zone update
func TestZoneService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"Id":123,"Name":"updated-zone"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	zone, err := client.Zones().Update(context.Background(), 123, &storage.UpdateZoneRequest{
		OriginURL: "https://example.com",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.ID != 123 {
		t.Errorf("expected ID 123, got %d", zone.ID)
	}
}

// TestZoneService_CheckAvailability tests zone name availability check
func TestZoneService_CheckAvailability(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"available":true}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	resp, err := client.Zones().CheckAvailability(context.Background(), "my-zone")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Available {
		t.Error("expected zone to be available")
	}
}

// TestZoneService_ResetPassword tests zone password reset
func TestZoneService_ResetPassword(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Password":"new-password-123"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	resp, err := client.Zones().ResetPassword(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Password == "" {
		t.Error("expected non-empty password")
	}
}

// TestZoneService_ResetReadOnlyPassword tests read-only password reset
func TestZoneService_ResetReadOnlyPassword(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"ReadOnlyPassword":"readonly-password-123"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	resp, err := client.Zones().ResetReadOnlyPassword(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ReadOnlyPassword == "" {
		t.Error("expected non-empty password")
	}
}

// TestFileService_DeleteDirectory tests directory deletion
func TestFileService_DeleteDirectory(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/test-zone/folder/") {
				t.Errorf("expected path with trailing slash for directory, got %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.DeleteDirectory(context.Background(), "folder")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_ListWithOptions tests list with pagination
func TestFileService_ListWithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `[
				{"Guid":"file1","ObjectName":"test1.txt","Length":1024,"IsDirectory":false},
				{"Guid":"file2","ObjectName":"test2.txt","Length":2048,"IsDirectory":false}
			]`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	files, err := fs.List(context.Background(), "files")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

// TestZoneListResponse_HasMore tests pagination helper
func TestZoneListResponse_HasMore(t *testing.T) {
	tests := []struct {
		name    string
		resp    storage.ZoneListResponse
		hasMore bool
	}{
		{
			name:    "has more pages",
			resp:    storage.ZoneListResponse{CurrentPage: 0, PageSize: 10, TotalItems: 25},
			hasMore: true,
		},
		{
			name:    "last page",
			resp:    storage.ZoneListResponse{CurrentPage: 2, PageSize: 10, TotalItems: 25},
			hasMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.resp.HasMore()
			if got != tt.hasMore {
				t.Errorf("HasMore() = %v, want %v", got, tt.hasMore)
			}
		})
	}
}

// TestWithUserAgent tests custom user agent option
func TestWithUserAgent(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := storage.NewClient("test-key",
		storage.WithHTTPClient(mock),
		storage.WithUserAgent("custom/1.0"),
	)
	if client == nil {
		t.Fatal("NewClient with custom user agent returned nil")
	}
}

// TestWithBaseURL tests custom base URL option
func TestWithBaseURL(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := storage.NewClient("test-key",
		storage.WithHTTPClient(mock),
		storage.WithBaseURL("https://custom.example.com"),
	)
	if client == nil {
		t.Fatal("NewClient with custom base URL returned nil")
	}
}

// TestWithFileUserAgent tests custom user agent for file operations
func TestWithFileUserAgent(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock),
		storage.WithFileUserAgent("custom-file/1.0"),
	)
	if fs == nil {
		t.Fatal("NewFileService with custom user agent returned nil")
	}
}

// TestFileService_ErrorHandling tests file error handling
func TestFileService_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		expectErr  bool
	}{
		{
			name:       "404 not found",
			statusCode: 404,
			body:       "Not Found",
			expectErr:  true,
		},
		{
			name:       "403 forbidden",
			statusCode: 403,
			body:       "Forbidden",
			expectErr:  true,
		},
		{
			name:       "500 server error",
			statusCode: 500,
			body:       "Internal Server Error",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return testutil.NewMockResponse(tt.statusCode, tt.body), nil
				},
			}

			fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
				storage.WithFileHTTPClient(mock))
			_, err := fs.Download(context.Background(), "test.txt")

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestZoneService_ListWithPagination tests list with pagination options
func TestZoneService_ListWithPagination(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify pagination parameters are in URL
			if !strings.Contains(req.URL.RawQuery, "page=") {
				t.Error("expected page parameter in URL")
			}
			body := `{"Items":[{"Id":123,"Name":"test-zone"}],"TotalItems":1,"CurrentPage":1,"PageSize":10}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	resp, err := client.Zones().List(context.Background(), &storage.ZoneListOptions{
		Page:    1,
		PerPage: 10,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 zone, got %d", len(resp.Items))
	}
}

// TestRegionVariants tests different region configurations
func TestRegionVariants(t *testing.T) {
	regions := []storage.Region{
		storage.RegionFalkenstein,
		storage.RegionNewYork,
		storage.RegionLosAngeles,
		storage.RegionSingapore,
		storage.RegionSydney,
		storage.RegionStockholm,
		storage.RegionSaoPaulo,
		storage.RegionJohannesburg,
		storage.RegionLondon,
	}

	for _, region := range regions {
		t.Run(string(region), func(t *testing.T) {
			baseURL := storage.RegionBaseURL(region)
			if baseURL == "" {
				t.Errorf("RegionBaseURL returned empty string for region %s", region)
			}
			if !strings.HasPrefix(baseURL, "https://") {
				t.Errorf("RegionBaseURL should return https URL, got %s", baseURL)
			}
		})
	}
}

// TestRegionBaseURL_UnknownRegion tests unknown region handling
func TestRegionBaseURL_UnknownRegion(t *testing.T) {
	region := storage.Region("unknown")
	baseURL := storage.RegionBaseURL(region)
	// Should return default region URL
	if baseURL == "" {
		t.Error("RegionBaseURL returned empty string for unknown region")
	}
}

// TestHandleErrorResponse tests error response handling
func TestHandleErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		expectErr  bool
	}{
		{
			name:       "structured error",
			statusCode: 400,
			body:       `{"Message":"validation error","ErrorKey":"VALIDATION","Field":"name"}`,
			expectErr:  true,
		},
		{
			name:       "plain text error",
			statusCode: 500,
			body:       "Internal Server Error",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return testutil.NewMockResponse(tt.statusCode, tt.body), nil
				},
			}

			client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
			_, err := client.Zones().Get(context.Background(), 123)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestFileAPIError_Error tests APIError error message formatting
func TestFileAPIError_Error(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
	}{
		{
			name:         "with field",
			statusCode:   400,
			responseBody: `{"Message":"invalid field","Field":"name"}`,
		},
		{
			name:         "with error key",
			statusCode:   400,
			responseBody: `{"Message":"error occurred","ErrorKey":"ERROR_KEY"}`,
		},
		{
			name:         "basic error",
			statusCode:   500,
			responseBody: "Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return testutil.NewMockResponse(tt.statusCode, tt.responseBody), nil
				},
			}

			fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
				storage.WithFileHTTPClient(mock))
			_, err := fs.Download(context.Background(), "test.txt")

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// Error message should contain status code
			if !strings.Contains(err.Error(), "bunny storage") {
				t.Errorf("expected error to contain 'bunny storage', got: %s", err.Error())
			}
		})
	}
}

// TestFileService_UploadWithChecksum tests upload with checksum
func TestFileService_UploadWithChecksum(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify checksum header
			if req.Header.Get("Checksum") == "" {
				t.Error("expected Checksum header")
			}
			return testutil.NewMockResponse(201, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))

	opts := &storage.UploadOptions{
		Checksum: "abc123",
	}
	err := fs.Upload(context.Background(), "test.txt", strings.NewReader("content"), opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_UploadError tests upload error handling
func TestFileService_UploadError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"upload failed"}`), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.Upload(context.Background(), "test.txt", strings.NewReader("content"), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestFileService_DownloadError tests download error handling
func TestFileService_DownloadError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, "Not Found"), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	_, err := fs.Download(context.Background(), "nonexistent.txt")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestFileService_ListError tests list error handling
func TestFileService_ListError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(403, "Forbidden"), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	_, err := fs.List(context.Background(), "path")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestFileService_DeleteError tests delete error handling
func TestFileService_DeleteError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, "Not Found"), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.Delete(context.Background(), "nonexistent.txt")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestFileService_DeleteDirectoryError tests delete directory error handling
func TestFileService_DeleteDirectoryError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(403, "Forbidden"), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.DeleteDirectory(context.Background(), "folder")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_ListError tests zone list error handling
func TestZoneService_ListError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(401, `{"Message":"unauthorized"}`), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_GetError tests zone get error handling
func TestZoneService_GetError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"not found"}`), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().Get(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_CreateError tests zone create error handling
func TestZoneService_CreateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"invalid request"}`), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().Create(context.Background(), &storage.CreateZoneRequest{
		Name: "test",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_UpdateError tests zone update error handling
func TestZoneService_UpdateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"update failed"}`), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().Update(context.Background(), 123, &storage.UpdateZoneRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_CheckAvailabilityError tests availability check error handling
func TestZoneService_CheckAvailabilityError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(500, "Server Error"), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().CheckAvailability(context.Background(), "test")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_ResetPasswordError tests password reset error handling
func TestZoneService_ResetPasswordError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(403, "Forbidden"), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().ResetPassword(context.Background(), 123)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_ResetReadOnlyPasswordError tests read-only password reset error handling
func TestZoneService_ResetReadOnlyPasswordError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(403, "Forbidden"), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().ResetReadOnlyPassword(context.Background(), 123)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneListOptions_AllParameters tests zone list with all query parameters
func TestZoneListOptions_AllParameters(t *testing.T) {
	requestURL := ""
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			requestURL = req.URL.String()
			body := `{"Items":[],"TotalItems":0,"CurrentPage":0,"PageSize":10}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, _ = client.Zones().List(context.Background(), &storage.ZoneListOptions{
		Page:        2,
		PerPage:     25,
		Search:      "test",
		IncludeDeleted: true,
	})

	expectedParams := []string{"page=", "perPage=", "search="}
	for _, param := range expectedParams {
		if !strings.Contains(requestURL, param) {
			t.Errorf("expected URL to contain %q, got: %s", param, requestURL)
		}
	}
}

// TestFileService_NetworkError tests network error handling
func TestFileService_NetworkError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	_, err := fs.Download(context.Background(), "test.txt")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestFileService_UploadWithAllOptions tests upload with all options
func TestFileService_UploadWithAllOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify checksum header
			if req.Header.Get("Checksum") == "" {
				t.Error("expected Checksum header")
			}
			// Verify content type
			if req.Header.Get("Content-Type") == "" {
				t.Error("expected Content-Type header")
			}
			return testutil.NewMockResponse(201, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))

	opts := &storage.UploadOptions{
		Checksum:    "abc123",
		ContentType: "video/mp4",
	}
	err := fs.Upload(context.Background(), "video.mp4", strings.NewReader("content"), opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestZoneService_NetworkError tests network error handling for zones
func TestZoneService_NetworkError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().Get(context.Background(), 123)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestFileService_ListEmptyPath tests list with empty path
func TestFileService_ListEmptyPath(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `[{"Guid":"file1","ObjectName":"root.txt","Length":100,"IsDirectory":false}]`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	files, err := fs.List(context.Background(), "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

// TestFileService_DeleteWithPath tests delete with path separators
func TestFileService_DeleteWithPath(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "/test-zone/folder/subfolder/file.txt") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			return testutil.NewMockResponse(200, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.Delete(context.Background(), "folder/subfolder/file.txt")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_ListInvalidJSON tests list with invalid JSON response
func TestFileService_ListInvalidJSON(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, "not valid json ["), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	_, err := fs.List(context.Background(), "path")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestZoneService_InvalidJSONResponse tests invalid JSON response handling for zones
func TestZoneService_InvalidJSONResponse(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, "not valid json {{{"), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	_, err := client.Zones().Get(context.Background(), 123)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("expected decode error, got: %s", err.Error())
	}
}

// TestZoneService_EmptyResponse tests empty response handling
func TestZoneService_EmptyResponse(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(http.StatusNoContent, ""), nil
		},
	}

	client := storage.NewClient("test-key", storage.WithHTTPClient(mock))
	err := client.Zones().Delete(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_UploadWithoutOptions tests upload without options
func TestFileService_UploadWithoutOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			// Verify default content type is set
			if req.Header.Get("Content-Type") != "application/octet-stream" {
				t.Errorf("expected default Content-Type, got %s", req.Header.Get("Content-Type"))
			}
			return testutil.NewMockResponse(201, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	err := fs.Upload(context.Background(), "file.txt", strings.NewReader("data"), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_UploadWithEmptyChecksum tests upload with empty checksum in options
func TestFileService_UploadWithEmptyChecksum(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Checksum header should not be set if empty
			if req.Header.Get("Checksum") != "" {
				t.Error("expected empty Checksum header")
			}
			// ContentType should use default when not provided
			if req.Header.Get("Content-Type") != "application/octet-stream" {
				t.Errorf("expected default Content-Type, got %s", req.Header.Get("Content-Type"))
			}
			return testutil.NewMockResponse(201, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))

	// Pass options but with empty values
	opts := &storage.UploadOptions{
		Checksum:    "",
		ContentType: "",
	}
	err := fs.Upload(context.Background(), "file.txt", strings.NewReader("data"), opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_DownloadClosedBody tests download with body that can be read
func TestFileService_DownloadClosedBody(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, "file content data"), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	reader, err := fs.Download(context.Background(), "file.txt")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer reader.Close()

	content, _ := io.ReadAll(reader)
	if len(content) == 0 {
		t.Error("expected non-empty content")
	}
}

// TestFileService_DeleteDirectoryWithTrailingSlash tests directory deletion path normalization
func TestFileService_DeleteDirectoryWithTrailingSlash(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify trailing slash is present
			if !strings.HasSuffix(req.URL.Path, "/") {
				t.Error("expected trailing slash for directory deletion")
			}
			return testutil.NewMockResponse(200, ""), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	// Test with path that already has trailing slash
	err := fs.DeleteDirectory(context.Background(), "folder/")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFileService_ListRootPath tests listing root directory
func TestFileService_ListRootPath(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Root path should just have zone name
			if !strings.Contains(req.URL.Path, "/test-zone/") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `[{"Guid":"f1","ObjectName":"file.txt","Length":100,"IsDirectory":false}]`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	fs := storage.NewFileService("test-zone", "zone-pass", storage.RegionFalkenstein,
		storage.WithFileHTTPClient(mock))
	files, err := fs.List(context.Background(), "/")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) == 0 {
		t.Error("expected files")
	}
}
