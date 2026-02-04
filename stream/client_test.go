package stream_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/geraldo/bunny-sdk-go/internal/testutil"
	"github.com/geraldo/bunny-sdk-go/stream"
)

func TestNewClient(t *testing.T) {
	client := stream.NewClient("test-api-key")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := stream.NewClient("test-key",
		stream.WithHTTPClient(mock),
		stream.WithUserAgent("custom/1.0"),
	)
	if client == nil {
		t.Fatal("NewClient with options returned nil")
	}
}

func TestVideoService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/library/123/videos") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			if req.Header.Get("AccessKey") != "test-key" {
				t.Errorf("expected AccessKey header")
			}

			body := `{"itemsPerPage":10,"currentPage":1,"totalItems":2,"items":[{"videoId":"vid1","title":"Test Video 1"},{"videoId":"vid2","title":"Test Video 2"}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(resp.Items))
	}
	if resp.Items[0].VideoID != "vid1" {
		t.Errorf("expected vid1, got %s", resp.Items[0].VideoID)
	}
}

func TestVideoService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}

			body := `{"videoId":"abc123","title":"Test Video","state":"finished","duration":120}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	video, err := client.Videos(123).Get(context.Background(), "abc123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.VideoID != "abc123" {
		t.Errorf("expected abc123, got %s", video.VideoID)
	}
	if video.Title != "Test Video" {
		t.Errorf("expected 'Test Video', got %s", video.Title)
	}
}

func TestVideoService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Error("expected Content-Type: application/json")
			}

			body := `{"videoId":"new-vid","title":"New Video","state":"created"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	video, err := client.Videos(123).Create(context.Background(), &stream.CreateVideoRequest{
		Title: "New Video",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.VideoID != "new-vid" {
		t.Errorf("expected new-vid, got %s", video.VideoID)
	}
}

func TestVideoService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).Delete(context.Background(), "abc123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVideoService_Upload(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if req.Header.Get("Content-Type") != "application/octet-stream" {
				t.Error("expected Content-Type: application/octet-stream")
			}
			return testutil.NewMockResponse(200, `{"videoId":"abc123","state":"processing"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).Upload(context.Background(), "abc123", strings.NewReader("video content"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLibraryService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"itemsPerPage":10,"currentPage":1,"totalItems":1,"items":[{"libraryId":123,"name":"Test Library"}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Libraries().List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestCollectionService_List(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"itemsPerPage":10,"currentPage":1,"totalItems":1,"items":[{"guid":"col1","name":"Test Collection"}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Collections(123).List(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestErrorHandling(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"Message":"Video not found"}`
			return testutil.NewMockResponse(404, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected 404 in error, got: %s", err.Error())
	}
}
// TestVideoService_Update tests video update
func TestVideoService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			body := `{"videoId":"vid1","title":"Updated Title"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	video, err := client.Videos(123).Update(context.Background(), "vid1", &stream.UpdateVideoRequest{
		Title: "Updated Title",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", video.Title)
	}
}

// TestVideoService_FetchFromURL tests video fetching from URL
func TestVideoService_FetchFromURL(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"success":true,"message":"Video is being fetched","statusCode":200}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).FetchFromURL(context.Background(), &stream.FetchVideoRequest{
		URL: "https://example.com/video.mp4",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

// TestVideoService_Reencode tests video re-encoding
func TestVideoService_Reencode(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).Reencode(context.Background(), "vid1", &stream.ReencodeRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoService_AddCaption tests adding video captions
func TestVideoService_AddCaption(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).AddCaption(context.Background(), "vid1", &stream.AddCaptionRequest{
		SrcLang: "en",
		Label:   "English",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoService_DeleteCaption tests deleting video captions
func TestVideoService_DeleteCaption(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).DeleteCaption(context.Background(), "vid1", "en")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoService_SetThumbnail tests setting video thumbnail
func TestVideoService_SetThumbnail(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"success":true,"message":"Thumbnail updated","statusCode":200}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).SetThumbnail(context.Background(), "vid1", &stream.SetThumbnailRequest{
		ThumbnailTime: 5000, // 5 seconds
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

// TestVideoService_GetHeatmap tests getting video heatmap
func TestVideoService_GetHeatmap(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"videoId":"vid1","heatmapData":[{"time":0,"plays":100},{"time":10,"plays":200},{"time":20,"plays":150}]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	heatmap, err := client.Videos(123).GetHeatmap(context.Background(), "vid1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(heatmap.HeatmapData) != 3 {
		t.Errorf("expected 3 data points, got %d", len(heatmap.HeatmapData))
	}
}

// TestVideoService_GetStatistics tests getting video statistics
func TestVideoService_GetStatistics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"viewsChart":{"01/15":100,"01/16":150}}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	stats, err := client.Videos(123).GetStatistics(context.Background(), "vid1", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats == nil {
		t.Error("expected stats to be non-nil")
	}
}

// TestVideoService_GetPlaybackInfo tests getting playback info
func TestVideoService_GetPlaybackInfo(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"videoLibraryId":123,"guid":"vid1","videoId":"vid1"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	info, err := client.Videos(123).GetPlaybackInfo(context.Background(), "vid1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.VideoID == "" {
		t.Error("expected non-empty video ID")
	}
}

// TestLibraryService_Get tests getting a single library
func TestLibraryService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"libraryId":123,"name":"Test Library"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	lib, err := client.Libraries().Get(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lib.LibraryID != 123 {
		t.Errorf("expected library ID 123, got %d", lib.LibraryID)
	}
}

// TestLibraryService_Create tests creating a library
func TestLibraryService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"libraryId":456,"name":"New Library"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	lib, err := client.Libraries().Create(context.Background(), &stream.CreateLibraryRequest{
		Name: "New Library",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lib.Name != "New Library" {
		t.Errorf("expected name 'New Library', got %s", lib.Name)
	}
}

// TestLibraryService_Update tests updating a library
func TestLibraryService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"libraryId":123,"name":"Updated Library"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	lib, err := client.Libraries().Update(context.Background(), 123, &stream.UpdateLibraryRequest{
		Name: "Updated Library",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lib.Name != "Updated Library" {
		t.Errorf("expected name 'Updated Library', got %s", lib.Name)
	}
}

// TestLibraryService_Delete tests deleting a library
func TestLibraryService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Libraries().Delete(context.Background(), 123)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestLibraryService_GetStatistics tests getting library statistics
func TestLibraryService_GetStatistics(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "dateFrom=") {
				t.Error("expected dateFrom parameter")
			}
			if !strings.Contains(query, "dateTo=") {
				t.Error("expected dateTo parameter")
			}
			body := `{"viewsChart":{"01/15":1000,"01/16":1500}}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	stats, err := client.Libraries().GetStatistics(context.Background(), 123, &stream.StatisticsOptions{
		DateFrom: "2024-01-01",
		DateTo:   "2024-01-31",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats == nil {
		t.Error("expected stats to be non-nil")
	}
}

// TestCollectionService_Get tests getting a single collection
func TestCollectionService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"guid":"col1","name":"Test Collection"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	coll, err := client.Collections(123).Get(context.Background(), "col1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coll.GUID != "col1" {
		t.Errorf("expected GUID col1, got %s", coll.GUID)
	}
}

// TestCollectionService_Create tests creating a collection
func TestCollectionService_Create(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"guid":"col2","name":"New Collection"}`
			return testutil.NewMockResponse(201, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	coll, err := client.Collections(123).Create(context.Background(), &stream.CreateCollectionRequest{
		Name: "New Collection",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coll.Name != "New Collection" {
		t.Errorf("expected name 'New Collection', got %s", coll.Name)
	}
}

// TestCollectionService_Update tests updating a collection
func TestCollectionService_Update(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"guid":"col1","name":"Updated Collection"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	coll, err := client.Collections(123).Update(context.Background(), "col1", &stream.UpdateCollectionRequest{
		Name: "Updated Collection",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coll.Name != "Updated Collection" {
		t.Errorf("expected name 'Updated Collection', got %s", coll.Name)
	}
}

// TestCollectionService_Delete tests deleting a collection
func TestCollectionService_Delete(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", req.Method)
			}
			return testutil.NewMockResponse(204, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Collections(123).Delete(context.Background(), "col1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoListResponse_HasMore tests pagination helper
func TestVideoListResponse_HasMore(t *testing.T) {
	tests := []struct {
		name    string
		resp    stream.VideoListResponse
		hasMore bool
	}{
		{
			name:    "has more pages",
			resp:    stream.VideoListResponse{CurrentPage: 1, ItemsPerPage: 10, TotalItems: 25},
			hasMore: true,
		},
		{
			name:    "last page",
			resp:    stream.VideoListResponse{CurrentPage: 3, ItemsPerPage: 10, TotalItems: 25},
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

// TestWithBaseAPIURL tests custom base API URL option
func TestWithBaseAPIURL(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := stream.NewClient("test-key",
		stream.WithHTTPClient(mock),
		stream.WithBaseAPIURL("https://custom-api.example.com"),
	)
	if client == nil {
		t.Fatal("NewClient with custom base API URL returned nil")
	}
}

// TestWithStreamAPIURL tests custom stream API URL option
func TestWithStreamAPIURL(t *testing.T) {
	mock := &testutil.MockHTTPClient{}
	client := stream.NewClient("test-key",
		stream.WithHTTPClient(mock),
		stream.WithStreamAPIURL("https://custom-stream.example.com"),
	)
	if client == nil {
		t.Fatal("NewClient with custom stream API URL returned nil")
	}
}

// TestBuildVideoListQuery tests query builder for video list
func TestBuildVideoListQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify query parameters are in URL
			query := req.URL.RawQuery
			if !strings.Contains(query, "page=2") {
				t.Error("expected page parameter")
			}
			if !strings.Contains(query, "itemsPerPage=50") {
				t.Error("expected itemsPerPage parameter")
			}
			body := `{"itemsPerPage":50,"currentPage":2,"totalItems":100,"items":[]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).List(context.Background(), &stream.VideoListOptions{
		Page:         2,
		ItemsPerPage: 50,
		
		OrderBy:      "title",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestBuildStatisticsQuery tests query builder for statistics
func TestBuildStatisticsQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "dateFrom=") {
				t.Error("expected dateFrom parameter")
			}
			if !strings.Contains(query, "dateTo=") {
				t.Error("expected dateTo parameter")
			}
			body := `{"viewsChart":{}}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).GetStatistics(context.Background(), "vid1", &stream.StatisticsOptions{
		DateFrom: "2024-01-01",
		DateTo:   "2024-01-31",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestBuildLibraryListQuery tests query builder for library list
func TestBuildLibraryListQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "page=") {
				t.Error("expected page parameter")
			}
			body := `{"itemsPerPage":10,"currentPage":1,"totalItems":0,"items":[]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Libraries().List(context.Background(), &stream.LibraryListOptions{
		Page:         1,
		ItemsPerPage: 10,
		
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestBuildCollectionListQuery tests query builder for collection list
func TestBuildCollectionListQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "page=") {
				t.Error("expected page parameter")
			}
			if !strings.Contains(query, "orderBy=") {
				t.Error("expected orderBy parameter")
			}
			if !strings.Contains(query, "includeThumbnails=") {
				t.Error("expected includeThumbnails parameter")
			}
			body := `{"itemsPerPage":10,"currentPage":1,"totalItems":0,"items":[]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Collections(123).List(context.Background(), &stream.CollectionListOptions{
		Page:              1,
		ItemsPerPage:      10,
		OrderBy:           "name",
		IncludeThumbnails: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestLibraryListResponse_HasMore tests pagination helper
func TestLibraryListResponse_HasMore(t *testing.T) {
	resp := stream.LibraryListResponse{CurrentPage: 1, ItemsPerPage: 10, TotalItems: 25}
	if !resp.HasMore() {
		t.Error("expected HasMore to be true")
	}
	resp2 := stream.LibraryListResponse{CurrentPage: 3, ItemsPerPage: 10, TotalItems: 25}
	if resp2.HasMore() {
		t.Error("expected HasMore to be false")
	}
}

// TestCollectionListResponse_HasMore tests pagination helper
func TestCollectionListResponse_HasMore(t *testing.T) {
	resp := stream.CollectionListResponse{CurrentPage: 1, ItemsPerPage: 10, TotalItems: 25}
	if !resp.HasMore() {
		t.Error("expected HasMore to be true")
	}
	resp2 := stream.CollectionListResponse{CurrentPage: 3, ItemsPerPage: 10, TotalItems: 25}
	if resp2.HasMore() {
		t.Error("expected HasMore to be false")
	}
}

// TestVideoService_ListError tests video list error handling
func TestVideoService_ListError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(401, `{"Message":"unauthorized"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_CreateError tests video create error handling
func TestVideoService_CreateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"invalid request"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Create(context.Background(), &stream.CreateVideoRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_UpdateError tests video update error handling
func TestVideoService_UpdateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"video not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Update(context.Background(), "vid1", &stream.UpdateVideoRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_FetchFromURLError tests fetch error handling
func TestVideoService_FetchFromURLError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"invalid URL"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).FetchFromURL(context.Background(), &stream.FetchVideoRequest{URL: "invalid"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_SetThumbnailError tests set thumbnail error handling
func TestVideoService_SetThumbnailError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"invalid thumbnail"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).SetThumbnail(context.Background(), "vid1", &stream.SetThumbnailRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_GetHeatmapError tests get heatmap error handling
func TestVideoService_GetHeatmapError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"heatmap not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).GetHeatmap(context.Background(), "vid1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_GetStatisticsError tests get statistics error handling
func TestVideoService_GetStatisticsError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(403, `{"Message":"forbidden"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).GetStatistics(context.Background(), "vid1", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_GetPlaybackInfoError tests get playback info error handling
func TestVideoService_GetPlaybackInfoError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"video not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).GetPlaybackInfo(context.Background(), "vid1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestLibraryService_ListError tests library list error handling
func TestLibraryService_ListError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(401, `{"Message":"unauthorized"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Libraries().List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestLibraryService_GetError tests library get error handling
func TestLibraryService_GetError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"library not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Libraries().Get(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestLibraryService_CreateError tests library create error handling
func TestLibraryService_CreateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"invalid request"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Libraries().Create(context.Background(), &stream.CreateLibraryRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestLibraryService_UpdateError tests library update error handling
func TestLibraryService_UpdateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"library not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Libraries().Update(context.Background(), 123, &stream.UpdateLibraryRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestLibraryService_GetStatisticsError tests get statistics error handling
func TestLibraryService_GetStatisticsError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(403, `{"Message":"forbidden"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Libraries().GetStatistics(context.Background(), 123, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestCollectionService_ListError tests collection list error handling
func TestCollectionService_ListError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(401, `{"Message":"unauthorized"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Collections(123).List(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestCollectionService_GetError tests collection get error handling
func TestCollectionService_GetError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"collection not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Collections(123).Get(context.Background(), "col1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestCollectionService_CreateError tests collection create error handling
func TestCollectionService_CreateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"invalid request"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Collections(123).Create(context.Background(), &stream.CreateCollectionRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestCollectionService_UpdateError tests collection update error handling
func TestCollectionService_UpdateError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"collection not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Collections(123).Update(context.Background(), "col1", &stream.UpdateCollectionRequest{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoAPIError_Error tests APIError formatting
func TestVideoAPIError_Error(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expectInErr  string
	}{
		{
			name:         "with field",
			statusCode:   400,
			responseBody: `{"Message":"invalid field","Field":"title"}`,
			expectInErr:  "field: title",
		},
		{
			name:         "with error key",
			statusCode:   400,
			responseBody: `{"Message":"error occurred","ErrorKey":"ERROR_KEY"}`,
			expectInErr:  "key: ERROR_KEY",
		},
		{
			name:         "basic error",
			statusCode:   500,
			responseBody: `{"Message":"server error"}`,
			expectInErr:  "server error",
		},
		{
			name:         "plain text error",
			statusCode:   503,
			responseBody: "Service Unavailable",
			expectInErr:  "Service Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &testutil.MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return testutil.NewMockResponse(tt.statusCode, tt.responseBody), nil
				},
			}

			client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
			_, err := client.Videos(123).Get(context.Background(), "vid1")

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), "bunny stream") {
				t.Errorf("expected error to contain 'bunny stream', got: %s", err.Error())
			}
			if !strings.Contains(err.Error(), tt.expectInErr) {
				t.Errorf("expected error to contain %q, got: %s", tt.expectInErr, err.Error())
			}
		})
	}
}

// TestVideoListOptions_AllParameters tests video list with all query parameters
func TestVideoListOptions_AllParameters(t *testing.T) {
	requestURL := ""
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			requestURL = req.URL.String()
			body := `{"itemsPerPage":25,"currentPage":2,"totalItems":100,"items":[]}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, _ = client.Videos(123).List(context.Background(), &stream.VideoListOptions{
		Page:              2,
		ItemsPerPage:      25,
		Search:            "test",
		OrderBy:           "date",
		Collection:        "col1",
		IncludeThumbnails: true,
	})

	expectedParams := []string{"page=", "itemsPerPage=", "orderBy=", "search=", "collection=", "includeThumbnails="}
	for _, param := range expectedParams {
		if !strings.Contains(requestURL, param) {
			t.Errorf("expected URL to contain %q, got: %s", param, requestURL)
		}
	}
}

// TestVideoService_NetworkError tests network error handling
func TestVideoService_NetworkError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Get(context.Background(), "vid1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_UploadRawData tests upload with raw data
func TestVideoService_UploadRawData(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			// Verify Content-Type header for raw upload
			if req.Header.Get("Content-Type") != "application/octet-stream" {
				t.Errorf("expected Content-Type: application/octet-stream, got %s", req.Header.Get("Content-Type"))
			}
			return testutil.NewMockResponse(200, `{"videoId":"vid1","state":"processing"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	content := "raw video data content"
	err := client.Videos(123).Upload(context.Background(), "vid1", strings.NewReader(content))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoService_UploadError tests upload error handling
func TestVideoService_UploadError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(413, `{"Message":"file too large"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).Upload(context.Background(), "vid1", strings.NewReader("data"))

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_InvalidJSONResponse tests invalid JSON response handling
func TestVideoService_InvalidJSONResponse(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(200, "not valid json {{{"), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Get(context.Background(), "vid1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("expected decode error, got: %s", err.Error())
	}
}

// TestVideoService_EmptyResponse tests empty response handling
func TestVideoService_EmptyResponse(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(http.StatusNoContent, ""), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	err := client.Videos(123).Delete(context.Background(), "vid1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoService_AddOutputCodec tests adding output codec
func TestVideoService_AddOutputCodec(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/outputs/0") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"videoId":"vid1","title":"Test"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	video, err := client.Videos(123).AddOutputCodec(context.Background(), "vid1", stream.OutputCodecX264)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.VideoID != "vid1" {
		t.Errorf("expected vid1, got %s", video.VideoID)
	}
}

// TestVideoService_AddOutputCodecError tests add output codec error handling
func TestVideoService_AddOutputCodecError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(400, `{"Message":"Invalid codec"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).AddOutputCodec(context.Background(), "vid1", stream.OutputCodecAV1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestVideoService_CleanupUnconfiguredResolutions tests cleanup resolutions
func TestVideoService_CleanupUnconfiguredResolutions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/resolutions/cleanup") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			if !strings.Contains(req.URL.RawQuery, "dryRun=true") {
				t.Errorf("expected dryRun query param, got: %s", req.URL.RawQuery)
			}
			body := `{"success":true,"message":"Cleanup complete","statusCode":200}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).CleanupUnconfiguredResolutions(context.Background(), "vid1", &stream.CleanupResolutionsOptions{
		DryRun: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected success to be true")
	}
}

// TestVideoService_GetHeatmapData tests getting heatmap data
func TestVideoService_GetHeatmapData(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/play/heatmap") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"showHeatmap":true,"enableDRM":false}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	data, err := client.Videos(123).GetHeatmapData(context.Background(), "vid1", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !data.ShowHeatmap {
		t.Error("expected showHeatmap to be true")
	}
}

// TestVideoService_GetStorageSizeInfo tests getting storage size info
func TestVideoService_GetStorageSizeInfo(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/storage") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"success":true,"statusCode":200,"data":{"thumbnails":1024,"originals":5000}}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).GetStorageSizeInfo(context.Background(), "vid1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Thumbnails != 1024 {
		t.Errorf("expected thumbnails 1024, got %d", resp.Data.Thumbnails)
	}
}

// TestVideoService_Repackage tests video repackage
func TestVideoService_Repackage(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/repackage") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"videoId":"vid1","title":"Repackaged"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	video, err := client.Videos(123).Repackage(context.Background(), "vid1", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if video.VideoID != "vid1" {
		t.Errorf("expected vid1, got %s", video.VideoID)
	}
}

// TestVideoService_Transcribe tests video transcription
func TestVideoService_Transcribe(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/transcribe") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"success":true,"message":"Transcription queued","statusCode":200}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).Transcribe(context.Background(), "vid1", &stream.TranscribeRequest{
		TargetLanguages: []string{"en", "es"},
	}, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected success to be true")
	}
}

// TestVideoService_TranscribeWithForce tests transcribe with force option
func TestVideoService_TranscribeWithForce(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "force=true") {
				t.Errorf("expected force query param, got: %s", req.URL.RawQuery)
			}
			body := `{"success":true,"message":"Transcription queued","statusCode":200}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.Videos(123).Transcribe(context.Background(), "vid1", nil, &stream.TranscribeOptions{Force: true})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestVideoService_TriggerSmartActions tests triggering smart actions
func TestVideoService_TriggerSmartActions(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/smart") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"success":true,"message":"Smart actions queued","statusCode":202}`
			return testutil.NewMockResponse(202, body), nil
		},
	}

	generateTitle := true
	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).TriggerSmartActions(context.Background(), "vid1", &stream.SmartActionsRequest{
		GenerateTitle: &generateTitle,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected success to be true")
	}
}

// TestVideoService_GetResolutionsInfo tests getting resolutions info
func TestVideoService_GetResolutionsInfo(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if !strings.HasSuffix(req.URL.Path, "/resolutions") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			body := `{"success":true,"statusCode":200,"data":{"videoId":"vid1","hasOriginal":true}}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.Videos(123).GetResolutionsInfo(context.Background(), "vid1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Data.HasOriginal {
		t.Error("expected hasOriginal to be true")
	}
}

// TestOEmbedService_Get tests oEmbed get
func TestOEmbedService_Get(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}
			if !strings.Contains(req.URL.Path, "/OEmbed") {
				t.Errorf("unexpected path: %s", req.URL.Path)
			}
			if !strings.Contains(req.URL.RawQuery, "url=") {
				t.Errorf("expected url query param, got: %s", req.URL.RawQuery)
			}
			body := `{"version":"1.0","type":"video","title":"Test Video","width":1920,"height":1080}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	resp, err := client.OEmbed().Get(context.Background(), &stream.OEmbedOptions{
		URL:      "https://iframe.mediadelivery.net/embed/123/abc",
		MaxWidth: 1920,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Type != "video" {
		t.Errorf("expected type video, got %s", resp.Type)
	}
}

// TestOEmbedService_GetError tests oEmbed error handling
func TestOEmbedService_GetError(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return testutil.NewMockResponse(404, `{"Message":"Video not found"}`), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, err := client.OEmbed().Get(context.Background(), &stream.OEmbedOptions{URL: "invalid"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestBuildCleanupResolutionsQuery tests cleanup resolutions query builder
func TestBuildCleanupResolutionsQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "resolutionsToDelete=720p") {
				t.Error("expected resolutionsToDelete parameter")
			}
			if !strings.Contains(query, "deleteNonConfiguredResolutions=true") {
				t.Error("expected deleteNonConfiguredResolutions parameter")
			}
			if !strings.Contains(query, "deleteOriginal=true") {
				t.Error("expected deleteOriginal parameter")
			}
			if !strings.Contains(query, "deleteMp4Files=true") {
				t.Error("expected deleteMp4Files parameter")
			}
			body := `{"success":true,"statusCode":200}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, _ = client.Videos(123).CleanupUnconfiguredResolutions(context.Background(), "vid1", &stream.CleanupResolutionsOptions{
		ResolutionsToDelete:            "720p",
		DeleteNonConfiguredResolutions: true,
		DeleteOriginal:                 true,
		DeleteMp4Files:                 true,
	})
}

// TestBuildHeatmapDataQuery tests heatmap data query builder
func TestBuildHeatmapDataQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "token=abc123") {
				t.Error("expected token parameter")
			}
			if !strings.Contains(query, "expires=") {
				t.Error("expected expires parameter")
			}
			body := `{"showHeatmap":true}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, _ = client.Videos(123).GetHeatmapData(context.Background(), "vid1", &stream.HeatmapDataOptions{
		Token:   "abc123",
		Expires: 1234567890,
	})
}

// TestBuildRepackageQuery tests repackage query builder
func TestBuildRepackageQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "keepOriginalFiles=true") {
				t.Error("expected keepOriginalFiles parameter")
			}
			body := `{"videoId":"vid1"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, _ = client.Videos(123).Repackage(context.Background(), "vid1", &stream.RepackageOptions{
		KeepOriginalFiles: true,
	})
}

// TestBuildOEmbedQuery tests oEmbed query builder
func TestBuildOEmbedQuery(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			query := req.URL.RawQuery
			if !strings.Contains(query, "maxWidth=1280") {
				t.Error("expected maxWidth parameter")
			}
			if !strings.Contains(query, "maxHeight=720") {
				t.Error("expected maxHeight parameter")
			}
			if !strings.Contains(query, "token=secret") {
				t.Error("expected token parameter")
			}
			body := `{"type":"video"}`
			return testutil.NewMockResponse(200, body), nil
		},
	}

	client := stream.NewClient("test-key", stream.WithHTTPClient(mock))
	_, _ = client.OEmbed().Get(context.Background(), &stream.OEmbedOptions{
		URL:       "https://example.com/video",
		MaxWidth:  1280,
		MaxHeight: 720,
		Token:     "secret",
		Expires:   9999999999,
	})
}
