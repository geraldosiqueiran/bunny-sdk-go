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
