# Code Standards & Implementation Guidelines

## File Naming Conventions

### Go Source Files
- **Kebab-case** for multi-word file names: `waf-rules-service.go`, `bot-detection-service.go`
- **Single-word files**: lowercase: `types.go`, `client.go`, `errors.go`
- **Test files**: Match source file name + `_test`: `video-service_test.go`, `client_test.go`
- **Package files**: One `client.go` and `types.go` per package minimum

**Examples:**
- ✅ `rate-limit-service.go` (clear intent)
- ✅ `log-forwarding-service.go` (descriptive)
- ❌ `rlservice.go` (unclear)
- ❌ `logfwd.go` (abbreviated)

---

## Code Formatting & Style

### General
- Standard `gofmt` formatting (no custom formatting)
- No line length hard limit, but keep under 120 chars for readability
- Blank lines between logical sections in functions
- Import grouping: stdlib, then project packages

### Comments
- Start with function name for exported functions: `// List returns...`
- Use clear, descriptive comments for non-obvious logic
- Comment exported types and interfaces

**Example:**
```go
// VideoService provides methods for video management operations.
type VideoService interface {
  List(ctx context.Context, opts *ListOptions) (*PaginatedResponse[Video], error)
  // ...
}
```

---

## Naming Conventions

### Package Names
- Lowercase, single word preferred: `stream`, `storage`, `shield`, `scripting`, `containers`
- Never use underscores in package names

### Type Names (PascalCase)
- Exported types: `Video`, `Library`, `Zone`, `File`, `WAFRule`
- Interfaces end in suffix: `VideoService`, `FileService`
- Error types: `APIError`, `NotFoundError`, `AuthError`, `RateLimitError`

**Examples:**
```go
type Video struct { /* ... */ }
type VideoService interface { /* ... */ }
type NotFoundError struct { /* ... */ }
```

### Function/Method Names (PascalCase)
- Resource operations: `List`, `Get`, `Create`, `Update`, `Delete`
- Resource checks: `CheckNameAvailability`, `IsActive`
- Helpers: `HasMore`, `IsAuthError`, `IsRetryable`

**Examples:**
```go
func (c *Client) Videos(libraryID int64) VideoService
func (s *videoService) List(ctx context.Context, opts *ListOptions) (*PaginatedResponse[Video], error)
func (e *APIError) IsAuthError() bool
```

### Struct Field Names (PascalCase)
- Exported fields use PascalCase
- JSON tags use camelCase (snake_case when matching Bunny.net API)
- Optional fields: use pointers

**Examples:**
```go
type Video struct {
  VideoID    int64     `json:"videoId"`
  Title      string    `json:"title"`
  Resolution *string   `json:"resolution,omitempty"`
  CreatedAt  time.Time `json:"dateCreated"`
}
```

### Constant Names (PascalCase)
- Exported constants use PascalCase
- Grouped with `const ( )`

**Examples:**
```go
const (
  RegionFalkenstein StorageRegion = "storage"
  RegionNewYork     StorageRegion = "ny"
  RegionLosAngeles  StorageRegion = "la"
)
```

### Variable Names (camelCase)
- Short names in tight scopes acceptable: `i`, `resp`, `opts`
- Longer names for wider scopes: `response`, `options`, `libraryID`

---

## Package Structure Conventions

### Minimal Package Layout
Every package should have:
- `client.go` - Client type and service accessors
- `types.go` - All data structures
- `{resource}-service.go` - Specific service implementations
- `errors.go` - Package-specific error types (if needed)
- `*_test.go` - Test files

### Client Pattern
```go
type Client struct {
  apiKey     string
  httpClient HTTPClient
  userAgent  string
  baseURL    string
}

func NewClient(apiKey string, opts ...Option) *Client {
  c := &Client{
    apiKey:     apiKey,
    httpClient: http.DefaultClient,
    userAgent:  defaultUserAgent,
    baseURL:    defaultBaseURL,
  }
  for _, opt := range opts {
    opt(c)
  }
  return c
}

// Service accessors
func (c *Client) Videos(libraryID int64) VideoService {
  return &videoService{client: c, libraryID: libraryID}
}
```

### Service Interface Pattern
```go
type VideoService interface {
  List(ctx context.Context, opts *ListOptions) (*PaginatedResponse[Video], error)
  Get(ctx context.Context, videoID string) (*Video, error)
  Create(ctx context.Context, req *CreateVideoRequest) (*Video, error)
  Update(ctx context.Context, videoID string, req *UpdateVideoRequest) (*Video, error)
  Delete(ctx context.Context, videoID string) error
}

// Private implementation
type videoService struct {
  client    *Client
  libraryID int64
}

func (s *videoService) List(ctx context.Context, opts *ListOptions) (*PaginatedResponse[Video], error) {
  // Implementation
}
```

### Types Pattern
- Group all types in `types.go`
- Organize by logical resource
- Include JSON tags matching Bunny.net API

**Example:**
```go
// Video types
type Video struct { /* ... */ }
type CreateVideoRequest struct { /* ... */ }
type UpdateVideoRequest struct { /* ... */ }

// Library types
type Library struct { /* ... */ }
type CreateLibraryRequest struct { /* ... */ }
```

---

## Error Handling Patterns

### Error Creation
Use `newError()` helper to create typed errors:

```go
if resp.StatusCode == 404 {
  return newError(resp.StatusCode, "not found", "", "")
}
```

Error type automatically determined by status code:
- 404 → NotFoundError
- 401/403 → AuthError
- 429 → RateLimitError
- Other → APIError

### Error Checking
```go
// Check error type with type assertion
if apiErr, ok := err.(*stream.APIError); ok {
  if apiErr.IsAuthError() {
    // Handle authentication error
  }
  if apiErr.IsRetryable() {
    // Implement retry logic
  }
}

// Or use specific types
if _, ok := err.(*stream.NotFoundError); ok {
  // Handle not found
}
```

### Error Fields
Include when available:
- `StatusCode` - HTTP status code
- `Message` - Human-readable error message
- `ErrorKey` - Bunny.net error key (optional)
- `Field` - Field causing validation error (optional)

---

## Type Design Patterns

### Optional Fields
Use **pointers for optional fields** in request/update types:

```go
type UpdateVideoRequest struct {
  Title      *string `json:"title,omitempty"`
  Description *string `json:"description,omitempty"`
  // Omitted fields not sent to API
}
```

### Request/Response Types
- `Create{Resource}Request` - POST request body
- `Update{Resource}Request` - PATCH request body
- `{Resource}` - Response type

### Pagination
```go
type PaginatedResponse[T any] struct {
  Items        []T   `json:"items"`
  TotalItems   int64 `json:"totalItems"`
  CurrentPage  int   `json:"currentPage"`
  ItemsPerPage int   `json:"itemsPerPage"`
}

func (p *PaginatedResponse[T]) HasMore() bool {
  return p.CurrentPage*int64(p.ItemsPerPage) < p.TotalItems
}
```

### Listing Options
```go
type ListOptions struct {
  Page        int
  ItemsPerPage int
  Search      string
  OrderBy     string
}
```

---

## Testing Patterns

### Table-Driven Tests
```go
func TestVideoService_List(t *testing.T) {
  tests := []struct {
    name    string
    opts    *ListOptions
    want    *PaginatedResponse[Video]
    wantErr bool
  }{
    {
      name: "list with pagination",
      opts: &ListOptions{Page: 1, ItemsPerPage: 10},
      want: &PaginatedResponse[Video]{Items: []Video{...}},
    },
    {
      name: "empty list",
      want: &PaginatedResponse[Video]{Items: []Video{}},
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Test implementation
    })
  }
}
```

### Mock HTTP Client
```go
mock := &testutil.MockHTTPClient{
  DoFunc: func(req *http.Request) (*http.Response, error) {
    return testutil.NewMockResponse(200, `{"videoId":"123"}`), nil
  },
}

client := stream.NewClient("key", stream.WithHTTPClient(mock))
```

### Error Testing
```go
func TestVideoService_NotFound(t *testing.T) {
  mock := &testutil.MockHTTPClient{
    DoFunc: func(req *http.Request) (*http.Response, error) {
      return testutil.NewMockResponse(404, `{"Message":"Not found"}`), nil
    },
  }

  client := stream.NewClient("key", stream.WithHTTPClient(mock))
  _, err := client.Videos(1).Get(context.Background(), "nonexistent")

  if !apiErr.IsNotFound() {
    t.Errorf("expected NotFoundError, got %v", err)
  }
}
```

---

## Authentication Patterns

### Root Package
```go
// Set auth header based on API key type
setAuthHeader(req, apiKey)
// Adds: Authorization: Bearer {apiKey}
```

### Stream Package
```go
client := stream.NewClient("library-api-key")
```

### Storage Zone Management
```go
client := storage.NewClient("global-api-key")
```

### Storage File Operations
```go
fileService := storage.NewFileService("zone-name", "zone-password", storage.RegionFalkenstein)
```

### Other Services (Shield, Scripting, Containers)
```go
shieldClient := shield.NewClient("global-api-key")
scriptClient := scripting.NewClient("global-api-key")
containerClient := containers.NewClient("global-api-key")
```

---

## Context Usage

All operations accept context for:
- Request cancellation
- Timeout enforcement
- Deadline propagation

**Pattern:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

video, err := client.Videos(123).Get(ctx, "video-id")
```

---

## JSON Tag Conventions

Match **exact Bunny.net API field names**:

```go
type Video struct {
  // camelCase fields (matching API)
  VideoID      int64     `json:"videoId"`
  LibraryID    int64     `json:"libraryId"`
  Title        string    `json:"title"`

  // snake_case fields (if API uses them)
  DateCreated  time.Time `json:"dateCreated"`

  // Omitempty for optional response fields
  Description  *string   `json:"description,omitempty"`

  // Ignored on request (only response)
  Created      time.Time `json:"dateCreated" `
}
```

---

## Import Organization

```go
import (
  // Standard library
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "strconv"

  // Project packages
  "github.com/geraldo/bunny-sdk-go/internal"
)
```

---

## Code Examples Checklist

When adding examples:
- ✅ Include imports
- ✅ Show error handling
- ✅ Use realistic IDs/values
- ✅ Demonstrate context usage
- ✅ Include defer statements for cleanup
- ✅ Test that examples compile

---

## Performance Considerations

### Memory
- Use io.Reader for streaming uploads/downloads
- Avoid copying response bodies
- Close response bodies properly: `defer resp.Body.Close()`

### HTTP
- Reuse http.DefaultClient (connection pooling)
- Support context for timeouts
- Use connection pooling via HTTP keep-alive

### JSON
- Use json.Decoder for response streams
- Use json.Marshal for request bodies
- Custom time parsing for flexibility

---

## Code Review Checklist

Before committing:
- [ ] Follows naming conventions (PascalCase types, camelCase functions)
- [ ] Error handling includes typed errors
- [ ] Context passed through all operations
- [ ] Response bodies closed
- [ ] JSON tags match Bunny.net API
- [ ] Tests included (table-driven)
- [ ] Comments on exported items
- [ ] No external dependencies

---

## Modularization Guidelines

When code files exceed 200 lines:
1. Identify logical separation points
2. Move related methods to separate service file
3. Keep shared types in `types.go`
4. Maintain single responsibility per file

**Example - Before:**
```
video.go (500 lines)
  - Video CRUD
  - Captions CRUD
  - Analytics
  - Transcription
```

**Example - After:**
```
video.go (200 lines) - Video CRUD
captions.go (100 lines) - Caption operations
analytics.go (100 lines) - Analytics
transcription.go (100 lines) - Transcription
```

---

## Consistency Checklist

For all new implementations:
- [ ] Package has client.go with NewClient()
- [ ] All public types in types.go
- [ ] Service interfaces in {resource}-service.go files
- [ ] Error handling with typed errors
- [ ] Full test coverage with table-driven tests
- [ ] Mock HTTP client support
- [ ] Context support on all operations
- [ ] Follows JSON tag conventions
- [ ] README example for package

