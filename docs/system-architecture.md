# System Architecture

## High-Level Overview

The Bunny.net Go SDK provides a layered architecture with clean separation between public API, service implementations, and HTTP transport layer. Zero external dependencies with idiomatic Go patterns throughout.

```
┌─────────────────────────────────────────────────────────┐
│         Application Code (User)                         │
└────────────────────────┬────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
┌───────▼────────┐ ┌────▼──────────┐ ┌──▼────────────┐
│  Stream API    │ │  Storage API  │ │ Shield API    │
│  (video.go)    │ │  (zone.go)    │ │ (waf*.go)     │
└────────┬───────┘ └────┬──────────┘ └──┬────────────┘
         │              │               │
    ┌────▼──────────────▼───────────────▼─────┐
    │      Service Layer                      │
    │  (Interfaces + Implementations)         │
    │  - VideoService                         │
    │  - ZoneService                          │
    │  - WAFService                           │
    │  - etc.                                 │
    └────┬──────────────────────────────────────┘
         │
    ┌────▼──────────────────────────────────┐
    │      Client Adapter Layer             │
    │  - Auth header injection              │
    │  - URL building                       │
    │  - Request/response marshaling        │
    │  - Error handling                     │
    └────┬──────────────────────────────────┘
         │
    ┌────▼──────────────────────────────────┐
    │    HTTP Transport Layer               │
    │  - http.Client interface              │
    │  - Mock support via HTTPClient        │
    └────┬──────────────────────────────────┘
         │
    ┌────▼──────────────────────────────────┐
    │   Bunny.net API Endpoints             │
    │  - api.bunny.net (management)         │
    │  - video.bunnycdn.com (stream)        │
    │  - {region}.bunnycdn.com (storage)    │
    └──────────────────────────────────────┘
```

---

## Client Hierarchy

### Root Package Clients

#### 1. Client (Management API)
Used for management operations across Storage Zones, Shield, Scripting, Containers.

```
NewClient(globalAPIKey)
  ├── Zones()              → ZoneService (storage)
  ├── Shields()            → ShieldService (shield)
  ├── Scripts()            → ScriptService (scripting)
  └── Applications()       → ApplicationService (containers)
```

**Base URL:** `https://api.bunny.net`

**Auth:** Global API Key (account-level)

#### 2. StreamClient (Stream API)
Dedicated client for video management.

```
NewStreamClient(libraryAPIKey)
  └── Videos(libraryID)     → VideoService
  ├── Libraries()           → LibraryService
  ├── Collections()         → CollectionService
  └── OEmbed()              → OEmbedService
```

**Base URL:** `https://video.bunnycdn.com`

**Auth:** Stream Library API Key (library-specific)

#### 3. StorageClient (Edge Storage)
Dedicated client for file operations.

```
NewStorageClient(zoneName, zonePassword, region)
  └── Upload(), Download(), List(), Delete()
```

**Base URL:** `https://{region}.bunnycdn.com`

**Auth:** Zone Password

---

## Request Flow

### Typical Operation Flow

```
┌─────────────────────────────────────────┐
│  1. User calls service method           │
│     videos.Get(ctx, videoID)            │
└────────────┬────────────────────────────┘
             │
┌────────────▼────────────────────────────┐
│  2. Method validates inputs             │
│     - Check context not nil             │
│     - Validate parameter ranges         │
└────────────┬────────────────────────────┘
             │
┌────────────▼────────────────────────────┐
│  3. Build request path                  │
│     "/library/{id}/videos/{videoID}"    │
└────────────┬────────────────────────────┘
             │
┌────────────▼────────────────────────────┐
│  4. Marshal request body (if needed)    │
│     json.Marshal(requestStruct)         │
└────────────┬────────────────────────────┘
             │
┌────────────▼────────────────────────────┐
│  5. Call doRequest()                    │
│     - Create HTTP request               │
│     - Set auth header                   │
│     - Set User-Agent header             │
│     - Set Content-Type/Accept           │
└────────────┬────────────────────────────┘
             │
┌────────────▼────────────────────────────┐
│  6. Execute HTTP request                │
│     httpClient.Do(req)                  │
└────────────┬────────────────────────────┘
             │
        ┌────┴────┐
        │          │
    ┌───▼──┐  ┌───▼──┐
    │ 2xx  │  │ 4xx+ │
    └───┬──┘  └───┬──┘
        │         │
   ┌────▼──┐  ┌──▼────────────────┐
   │ 7a.   │  │  7b. Parse error  │
   │ Parse │  │  - Read body      │
   │ body  │  │  - Unmarshal JSON │
   │ and   │  │  - Create typed   │
   │return │  │    error          │
   └───────┘  └──────────┬────────┘
                         │
              ┌──────────▼────────────┐
              │  8. Return error      │
              │  - NotFoundError      │
              │  - AuthError          │
              │  - RateLimitError     │
              │  - APIError           │
              └───────────────────────┘
```

---

## Authentication Flow

### Header Injection Pattern

```
Request created with:
  - Method (GET, POST, PATCH, DELETE)
  - URL (baseURL + path)
  - Body (JSON-encoded or nil)

↓

setAuthHeader(request, apiKey) adds:
  Authorization: Bearer {apiKey}

↓

Additional headers added:
  User-Agent: bunny-sdk-go/1.0
  Content-Type: application/json
  Accept: application/json

↓

Request sent to Bunny.net endpoint
```

### Authentication by Service

```
┌──────────────────────┬──────────────────┬──────────────────┐
│ Service              │ Key Type         │ Location         │
├──────────────────────┼──────────────────┼──────────────────┤
│ Stream API           │ Library API Key  │ Library settings │
│ Storage Management   │ Global API Key   │ Account settings │
│ Storage File Ops     │ Zone Password    │ Zone settings    │
│ Shield/WAF           │ Global API Key   │ Account settings │
│ Edge Scripting       │ Global API Key   │ Account settings │
│ Magic Containers     │ Global API Key   │ Account settings │
└──────────────────────┴──────────────────┴──────────────────┘
```

---

## Error Handling Flow

### Error Detection & Transformation

```
HTTP Response received
         │
         ▼
Status code >= 400?
    │         └─ No  → Parse successful response
    │              → Return result
    Yes
    │
    ▼
Read response body
    │
    ▼
Try to unmarshal JSON error:
  {
    "Message": "error message",
    "ErrorKey": "validation_error",
    "Field": "title"
  }
    │
    ▼
Determine error type by status code:
    │
    ├─ 404 Not Found
    │   └─ Create NotFoundError
    │
    ├─ 401 Unauthorized / 403 Forbidden
    │   └─ Create AuthError
    │
    ├─ 429 Too Many Requests
    │   └─ Create RateLimitError
    │
    └─ Other
        └─ Create APIError

         ▼
Return typed error with:
  - StatusCode
  - Message
  - ErrorKey (optional)
  - Field (optional)
```

### Error Type Hierarchy

```
error (interface)
  │
  └─ APIError
      ├─ NotFoundError (404)
      ├─ AuthError (401/403)
      └─ RateLimitError (429)

Helper methods on APIError:
  - IsAuthError()   → bool
  - IsNotFound()    → bool
  - IsRateLimited() → bool
  - IsRetryable()   → bool
```

---

## Service Pattern

### Interface Definition

```go
type VideoService interface {
  List(ctx context.Context, opts *ListOptions) (*PaginatedResponse[Video], error)
  Get(ctx context.Context, videoID string) (*Video, error)
  Create(ctx context.Context, req *CreateVideoRequest) (*Video, error)
  Update(ctx context.Context, videoID string, req *UpdateVideoRequest) (*Video, error)
  Delete(ctx context.Context, videoID string) error
  Upload(ctx context.Context, videoID string, file io.Reader) error
  // Additional methods...
}
```

### Private Implementation

```go
type videoService struct {
  client    *Client
  libraryID int64
}

func (s *videoService) Get(ctx context.Context, videoID string) (*Video, error) {
  path := fmt.Sprintf("/library/%d/videos/%s", s.libraryID, videoID)
  var video Video
  if err := s.client.doRequest(ctx, "GET", path, nil, &video); err != nil {
    return nil, err
  }
  return &video, nil
}
```

### Service Accessors

```go
func (c *Client) Videos(libraryID int64) VideoService {
  return &videoService{
    client:    c,
    libraryID: libraryID,
  }
}
```

### Benefits of Service Pattern

- **Mockability** - Interface for testing
- **Separation of Concerns** - Service implementation isolated
- **Resource Scoping** - Service parameterized with IDs (libraryID, zoneID)
- **Consistency** - Same pattern across all packages

---

## Type System Design

### Request Types

```go
type CreateVideoRequest struct {
  Title       string  `json:"title"`
  Description *string `json:"description,omitempty"`
}

type UpdateVideoRequest struct {
  Title       *string `json:"title,omitempty"`
  Description *string `json:"description,omitempty"`
}
```

**Design Decision:**
- Use pointers for optional fields
- Only included fields sent to API (omitempty)
- Prevents accidental overwrites

### Response Types

```go
type Video struct {
  VideoID      int64     `json:"videoId"`
  Title        string    `json:"title"`
  Description  string    `json:"description"`
  CreatedAt    time.Time `json:"dateCreated"`
  Views        int64     `json:"views"`
}
```

**Design Decision:**
- Non-pointer fields for guaranteed presence
- Custom time parsing for Bunny.net format
- Direct mapping to JSON response

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

**Design Decision:**
- Generic T for type safety
- Helper methods for navigation
- Follows Bunny.net API format

---

## Data Flow Examples

### Example 1: Create Video

```
Client Code:
  video, err := client.Videos(12345).Create(ctx, &CreateVideoRequest{
    Title: "My Video"
  })

Flow:
  1. videoService.Create() called with request
  2. Request marshaled: {"title":"My Video"}
  3. POST /library/12345/videos with body
  4. Auth header added: Authorization: Bearer {key}
  5. Response received and unmarshaled to Video
  6. Video returned to caller
```

### Example 2: Upload File

```
Client Code:
  file, _ := os.Open("video.mp4")
  err := client.Videos(12345).Upload(ctx, "video-id", file)

Flow:
  1. videoService.Upload() called with io.Reader
  2. POST /library/12345/videos/video-id/source
  3. Stream file directly (no buffering)
  4. Auth header added
  5. Response indicates success/failure
```

### Example 3: Error Handling

```
Client Code:
  video, err := client.Videos(12345).Get(ctx, "invalid")
  if err != nil {
    apiErr := err.(*APIError)
    if apiErr.IsNotFound() {
      // Handle 404
    }
  }

Flow:
  1. GET /library/12345/videos/invalid
  2. Server returns 404 with body:
     {"Message":"Not found"}
  3. Response status checked (404)
  4. Body parsed to extract Message
  5. NotFoundError created with:
     - StatusCode: 404
     - Message: "Not found"
  6. Error returned to caller
```

---

## Storage Region Architecture

### Multi-Region Support

```
User specifies region when creating StorageClient:

  client := NewStorageClient(
    "zone-name",
    "zone-password",
    RegionFalkenstein,  // or ny, la, sg, syd, se, br, jh, uk
  )

Base URL generated:
  https://{region}.bunnycdn.com

Request to:
  https://storage.bunnycdn.com/zone-name/path/to/file
  https://ny.bunnycdn.com/zone-name/path/to/file
  https://la.bunnycdn.com/zone-name/path/to/file
  // etc.
```

### Region Constants

```go
const (
  RegionFalkenstein = "storage"   // Germany (default)
  RegionNewYork     = "ny"
  RegionLosAngeles  = "la"
  RegionSingapore   = "sg"
  RegionSydney      = "syd"
  RegionStockholm   = "se"
  RegionSaoPaulo    = "br"
  RegionJohannesburg = "jh"
  RegionLondon      = "uk"
)
```

---

## Package Independence

Each package designed independently:

```
stream/
  ├─ Standalone from storage
  ├─ Can use NewStreamClient() alone
  └─ Library-scoped API key

storage/
  ├─ Standalone from stream
  ├─ Two ways: NewClient() or NewStorageClient()
  └─ Zone/Global API keys

shield/
  ├─ Standalone
  ├─ Uses NewClient() pattern
  └─ Global API key

scripting/
  ├─ Standalone
  ├─ Uses NewClient() pattern
  └─ Global API key

containers/
  ├─ Standalone
  ├─ Uses NewClient() pattern
  └─ Global API key
```

---

## Concurrency & Context

### Context Propagation

```
┌─ Context with timeout
│    ctx, cancel := context.WithTimeout(...)
│    defer cancel()
│
└─ Passed through entire call stack
     client.Videos().Get(ctx, id)
       └─ videoService.Get(ctx, ...)
           └─ client.doRequest(ctx, ...)
               └─ http.NewRequestWithContext(ctx, ...)
```

### Timeout Enforcement

```
Request → Get → doRequest → http.Do() → Timeout?
                                    │
                                    ├─ Yes → context.Canceled
                                    └─ No  → Response
```

### Cancellation

User can cancel anytime:

```go
cancel()  // Cancels context
// Next operation fails with context.Canceled
```

---

## Testing Architecture

### Mock HTTP Client

```go
type MockHTTPClient struct {
  DoFunc func(*http.Request) (*http.Response, error)
}

Usage:
  mock := &testutil.MockHTTPClient{
    DoFunc: func(req *http.Request) (*http.Response, error) {
      return testutil.NewMockResponse(200, `{...}`), nil
    },
  }

  client := NewClient("key", WithHTTPClient(mock))
```

### Test Isolation

```
Each test:
  1. Create mock HTTP client with custom DoFunc
  2. Create client with mock
  3. Call service method
  4. Assert result
  5. No actual HTTP requests
```

---

## Design Decisions

### Zero Dependencies
- **Why:** Minimal footprint, no version conflicts
- **How:** Use only Go stdlib (net/http, encoding/json, etc.)
- **Trade-off:** No built-in retry logic (user can implement)

### Generics for Pagination
- **Why:** Type-safe pagination across all resources
- **How:** PaginatedResponse[T] with generic T
- **Benefit:** No type assertions needed

### Service Interfaces
- **Why:** Mockable for unit testing
- **How:** Service interface + private implementation
- **Benefit:** Easy to mock, test in isolation

### Functional Options
- **Why:** Flexible client configuration
- **How:** Option func(*Client) pattern
- **Benefit:** Future-proof API

### Pointer Fields for Updates
- **Why:** Distinguish omitted vs. null fields
- **How:** Update{Resource}Request uses *string, *int, etc.
- **Benefit:** Only included fields sent to API

---

## Extension Points

### Adding New Service

1. Create `{service}-service.go` with interface
2. Add types to `types.go`
3. Register in `client.go` as accessor method
4. Implement error handling if needed
5. Add tests

### Custom HTTP Client

```go
client := NewClient("key", WithHTTPClient(customTransport))
```

### Custom User Agent

```go
client := NewClient("key", WithUserAgent("my-app/1.0"))
```

### Custom Base URLs

```go
client := NewClient("key",
  WithStreamBaseURL("https://custom.url"),
  WithStorageBaseURL("https://storage.custom.url"),
)
```

---

## Performance Characteristics

### Memory
- Streaming file uploads/downloads (no buffering)
- Response body read once and closed
- No unnecessary allocations in hot paths

### Network
- HTTP Keep-Alive via http.DefaultClient
- Connection pooling
- One HTTP request per operation (no pipelining)

### Concurrency
- Safe for concurrent use (all methods read context, no shared state mutation)
- Each call independent
- Context supports cancellation across goroutines

