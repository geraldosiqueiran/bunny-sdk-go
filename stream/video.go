package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/geraldo/bunny-sdk-go/internal"
)

// VideoService provides methods for managing videos in a library.
type VideoService interface {
	List(ctx context.Context, opts *VideoListOptions) (*VideoListResponse, error)
	Get(ctx context.Context, videoID string) (*Video, error)
	Create(ctx context.Context, req *CreateVideoRequest) (*Video, error)
	Update(ctx context.Context, videoID string, req *UpdateVideoRequest) (*Video, error)
	Delete(ctx context.Context, videoID string) error
	Upload(ctx context.Context, videoID string, reader io.Reader) error
	FetchFromURL(ctx context.Context, req *FetchVideoRequest) (*FetchVideoResponse, error)
	Reencode(ctx context.Context, videoID string, req *ReencodeRequest) error
	AddCaption(ctx context.Context, videoID string, req *AddCaptionRequest) error
	DeleteCaption(ctx context.Context, videoID, srclang string) error
	SetThumbnail(ctx context.Context, videoID string, req *SetThumbnailRequest) (*SetThumbnailResponse, error)
	GetHeatmap(ctx context.Context, videoID string) (*HeatmapData, error)
	GetStatistics(ctx context.Context, videoID string, opts *StatisticsOptions) (*VideoStatistics, error)
	GetPlaybackInfo(ctx context.Context, videoID string) (*PlaybackInfo, error)
}

// VideoListResponse represents a paginated list of videos.
type VideoListResponse struct {
	ItemsPerPage int     `json:"itemsPerPage"`
	CurrentPage  int     `json:"currentPage"`
	TotalItems   int     `json:"totalItems"`
	Items        []Video `json:"items"`
}

// HasMore returns true if there are more pages to fetch.
func (r *VideoListResponse) HasMore() bool {
	return r.CurrentPage*r.ItemsPerPage < r.TotalItems
}

type videoService struct {
	client    httpClient
	libraryID int64
}

type httpClient interface {
	do(ctx context.Context, method, path string, body any, result any) error
	doRaw(ctx context.Context, method, path string, body io.Reader, contentType string) error
}

func newVideoService(client httpClient, libraryID int64) VideoService {
	return &videoService{client: client, libraryID: libraryID}
}

// List returns a paginated list of videos in the library.
func (s *videoService) List(ctx context.Context, opts *VideoListOptions) (*VideoListResponse, error) {
	path := fmt.Sprintf("/library/%d/videos", s.libraryID)
	if opts != nil {
		path = path + "?" + buildVideoListQuery(opts)
	}

	var resp VideoListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single video by ID.
func (s *videoService) Get(ctx context.Context, videoID string) (*Video, error) {
	path := fmt.Sprintf("/library/%d/videos/%s", s.libraryID, videoID)

	var video Video
	if err := s.client.do(ctx, http.MethodGet, path, nil, &video); err != nil {
		return nil, err
	}
	return &video, nil
}

// Create creates a new video entry (before uploading the actual video file).
func (s *videoService) Create(ctx context.Context, req *CreateVideoRequest) (*Video, error) {
	path := fmt.Sprintf("/library/%d/videos", s.libraryID)

	var video Video
	if err := s.client.do(ctx, http.MethodPost, path, req, &video); err != nil {
		return nil, err
	}
	return &video, nil
}

// Update updates video metadata.
func (s *videoService) Update(ctx context.Context, videoID string, req *UpdateVideoRequest) (*Video, error) {
	path := fmt.Sprintf("/library/%d/videos/%s", s.libraryID, videoID)

	var video Video
	if err := s.client.do(ctx, http.MethodPost, path, req, &video); err != nil {
		return nil, err
	}
	return &video, nil
}

// Delete permanently deletes a video.
func (s *videoService) Delete(ctx context.Context, videoID string) error {
	path := fmt.Sprintf("/library/%d/videos/%s", s.libraryID, videoID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// Upload uploads a video file to an existing video entry.
func (s *videoService) Upload(ctx context.Context, videoID string, reader io.Reader) error {
	path := fmt.Sprintf("/library/%d/videos/%s", s.libraryID, videoID)
	return s.client.doRaw(ctx, http.MethodPut, path, reader, "application/octet-stream")
}

// FetchFromURL fetches a video from a remote URL.
func (s *videoService) FetchFromURL(ctx context.Context, req *FetchVideoRequest) (*FetchVideoResponse, error) {
	path := fmt.Sprintf("/library/%d/videos/fetch", s.libraryID)

	var resp FetchVideoResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Reencode triggers re-encoding of a video with optional resolution settings.
func (s *videoService) Reencode(ctx context.Context, videoID string, req *ReencodeRequest) error {
	path := fmt.Sprintf("/library/%d/videos/%s/reencode", s.libraryID, videoID)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}

// AddCaption adds a caption track to a video.
func (s *videoService) AddCaption(ctx context.Context, videoID string, req *AddCaptionRequest) error {
	path := fmt.Sprintf("/library/%d/videos/%s/captions", s.libraryID, videoID)
	return s.client.do(ctx, http.MethodPost, path, req, nil)
}

// DeleteCaption removes a caption track from a video.
func (s *videoService) DeleteCaption(ctx context.Context, videoID, srclang string) error {
	path := fmt.Sprintf("/library/%d/videos/%s/captions/%s", s.libraryID, videoID, srclang)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// SetThumbnail sets the video thumbnail from a specific timestamp.
func (s *videoService) SetThumbnail(ctx context.Context, videoID string, req *SetThumbnailRequest) (*SetThumbnailResponse, error) {
	path := fmt.Sprintf("/library/%d/videos/%s/thumbnail", s.libraryID, videoID)

	var resp SetThumbnailResponse
	if err := s.client.do(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetHeatmap returns engagement heatmap data for a video.
func (s *videoService) GetHeatmap(ctx context.Context, videoID string) (*HeatmapData, error) {
	path := fmt.Sprintf("/library/%d/videos/%s/heatmap", s.libraryID, videoID)

	var resp HeatmapData
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStatistics returns statistics for a video.
func (s *videoService) GetStatistics(ctx context.Context, videoID string, opts *StatisticsOptions) (*VideoStatistics, error) {
	path := fmt.Sprintf("/library/%d/videos/%s/statistics", s.libraryID, videoID)
	if opts != nil {
		path = path + "?" + buildStatisticsQuery(opts)
	}

	var resp VideoStatistics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPlaybackInfo returns playback URLs for a video.
func (s *videoService) GetPlaybackInfo(ctx context.Context, videoID string) (*PlaybackInfo, error) {
	path := fmt.Sprintf("/library/%d/videos/%s/play", s.libraryID, videoID)

	var resp PlaybackInfo
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildVideoListQuery(opts *VideoListOptions) string {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.ItemsPerPage > 0 {
		params.Set("itemsPerPage", strconv.Itoa(opts.ItemsPerPage))
	}
	if opts.Search != "" {
		params.Set("search", opts.Search)
	}
	if opts.Collection != "" {
		params.Set("collection", opts.Collection)
	}
	if opts.OrderBy != "" {
		params.Set("orderBy", opts.OrderBy)
	}
	if opts.IncludeThumbnails {
		params.Set("includeThumbnails", "true")
	}
	return params.Encode()
}

func buildStatisticsQuery(opts *StatisticsOptions) string {
	params := url.Values{}
	if opts.DateFrom != "" {
		params.Set("dateFrom", opts.DateFrom)
	}
	if opts.DateTo != "" {
		params.Set("dateTo", opts.DateTo)
	}
	return params.Encode()
}

// doRequestJSON is a helper for JSON requests used by the stream client.
func doRequestJSON(ctx context.Context, httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}, baseURL, apiKey, userAgent, method, path string, body any, result any) error {
	fullURL := baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := internal.NewRequest(ctx, method, fullURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return handleErrorResponse(resp)
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	} else {
		resp.Body.Close()
	}

	return nil
}

// doRequestRaw is a helper for raw body requests (file uploads).
func doRequestRaw(ctx context.Context, httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}, baseURL, apiKey, userAgent, method, path string, body io.Reader, contentType string) error {
	fullURL := baseURL + path

	req, err := internal.NewRequest(ctx, method, fullURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("User-Agent", userAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return handleErrorResponse(resp)
	}

	return nil
}

func handleErrorResponse(resp *http.Response) error {
	body, err := internal.ReadResponseBody(resp)
	if err != nil {
		return newAPIError(resp.StatusCode, "failed to read error response", "", "")
	}

	errResp := internal.ParseErrorResponse(body)
	if errResp != nil && errResp.Message != "" {
		return newAPIError(resp.StatusCode, errResp.Message, errResp.ErrorKey, errResp.Field)
	}

	return newAPIError(resp.StatusCode, string(body), "", "")
}

// APIError represents an error from the Bunny.net API.
type APIError struct {
	StatusCode int
	Message    string
	ErrorKey   string
	Field      string
}

func (e *APIError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("bunny stream: %s (status: %d, field: %s)", e.Message, e.StatusCode, e.Field)
	}
	if e.ErrorKey != "" {
		return fmt.Sprintf("bunny stream: %s (status: %d, key: %s)", e.Message, e.StatusCode, e.ErrorKey)
	}
	return fmt.Sprintf("bunny stream: %s (status: %d)", e.Message, e.StatusCode)
}

func newAPIError(statusCode int, message, errorKey, field string) error {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		ErrorKey:   errorKey,
		Field:      field,
	}
}
