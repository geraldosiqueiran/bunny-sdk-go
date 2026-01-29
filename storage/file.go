package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/geraldo/bunny-sdk-go/internal"
)

// FileService provides methods for file operations in a storage zone.
type FileService interface {
	Upload(ctx context.Context, path string, reader io.Reader, opts *UploadOptions) error
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	List(ctx context.Context, path string) ([]File, error)
	Delete(ctx context.Context, path string) error
	DeleteDirectory(ctx context.Context, path string) error
}

type fileService struct {
	httpClient HTTPClient
	baseURL    string
	zoneName   string
	accessKey  string
	userAgent  string
}

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewFileService creates a new FileService for file operations.
// Use the storage zone password (NOT global API key) as accessKey.
func NewFileService(zoneName, accessKey string, region Region, opts ...FileServiceOption) FileService {
	fs := &fileService{
		httpClient: http.DefaultClient,
		baseURL:    RegionBaseURL(region),
		zoneName:   zoneName,
		accessKey:  accessKey,
		userAgent:  "bunny-sdk-go/1.0",
	}
	for _, opt := range opts {
		opt(fs)
	}
	return fs
}

// FileServiceOption is a functional option for configuring FileService.
type FileServiceOption func(*fileService)

// WithFileHTTPClient sets a custom HTTP client for file operations.
func WithFileHTTPClient(hc HTTPClient) FileServiceOption {
	return func(fs *fileService) {
		fs.httpClient = hc
	}
}

// WithFileUserAgent sets a custom user agent for file operations.
func WithFileUserAgent(ua string) FileServiceOption {
	return func(fs *fileService) {
		fs.userAgent = ua
	}
}

// Upload uploads a file to the storage zone.
// The path should not include the zone name (e.g., "documents/report.pdf").
// Directories are created automatically.
func (s *fileService) Upload(ctx context.Context, path string, reader io.Reader, opts *UploadOptions) error {
	fullURL := s.buildURL(path)

	req, err := internal.NewRequest(ctx, http.MethodPut, fullURL, reader)
	if err != nil {
		return err
	}

	s.setHeaders(req)
	if opts != nil {
		if opts.Checksum != "" {
			req.Header.Set("Checksum", opts.Checksum)
		}
		if opts.ContentType != "" {
			req.Header.Set("Content-Type", opts.ContentType)
		} else {
			req.Header.Set("Content-Type", "application/octet-stream")
		}
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return s.handleError(resp)
	}

	return nil
}

// Download downloads a file from the storage zone.
// The caller is responsible for closing the returned ReadCloser.
func (s *fileService) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	fullURL := s.buildURL(path)

	req, err := internal.NewRequest(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, s.handleError(resp)
	}

	return resp.Body, nil
}

// List lists files and directories at the given path.
// The path should be a directory path (trailing slash is added automatically).
func (s *fileService) List(ctx context.Context, path string) ([]File, error) {
	// Ensure trailing slash for directory listing
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	fullURL := s.buildURL(path)

	req, err := internal.NewRequest(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.handleError(resp)
	}

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return files, nil
}

// Delete deletes a file from the storage zone.
func (s *fileService) Delete(ctx context.Context, path string) error {
	// Ensure no trailing slash for file deletion
	path = strings.TrimSuffix(path, "/")

	fullURL := s.buildURL(path)

	req, err := internal.NewRequest(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return err
	}

	s.setHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return s.handleError(resp)
	}

	return nil
}

// DeleteDirectory deletes a directory and all its contents recursively.
func (s *fileService) DeleteDirectory(ctx context.Context, path string) error {
	// Ensure trailing slash for directory deletion
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	fullURL := s.buildURL(path)

	req, err := internal.NewRequest(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return err
	}

	s.setHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete directory failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return s.handleError(resp)
	}

	return nil
}

func (s *fileService) buildURL(path string) string {
	// Remove leading slash from path
	path = strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s/%s/%s", s.baseURL, s.zoneName, path)
}

func (s *fileService) setHeaders(req *http.Request) {
	req.Header.Set("AccessKey", s.accessKey)
	req.Header.Set("User-Agent", s.userAgent)
}

func (s *fileService) handleError(resp *http.Response) error {
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

// APIError represents an error from the Bunny.net Storage API.
type APIError struct {
	StatusCode int
	Message    string
	ErrorKey   string
	Field      string
}

func (e *APIError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("bunny storage: %s (status: %d, field: %s)", e.Message, e.StatusCode, e.Field)
	}
	if e.ErrorKey != "" {
		return fmt.Sprintf("bunny storage: %s (status: %d, key: %s)", e.Message, e.StatusCode, e.ErrorKey)
	}
	return fmt.Sprintf("bunny storage: %s (status: %d)", e.Message, e.StatusCode)
}

func newAPIError(statusCode int, message, errorKey, field string) error {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		ErrorKey:   errorKey,
		Field:      field,
	}
}
