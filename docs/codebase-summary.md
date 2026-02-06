# Bunny.net Go SDK - Codebase Summary

## Overview

This Go SDK provides idiomatic, zero-dependency access to Bunny.net cloud platform APIs. Implements 5 major API packages with 70+ service methods supporting video streaming, edge storage, security (WAF), edge scripting, and containerized applications.

**Total Implementation**: ~4,500 lines of Go code | ~7,300 lines of tests | **Zero external dependencies** (stdlib only)

---

## Directory Structure

### Root Package `/`
Core client types and shared utilities

**Key Files:**
- `bunny.go` (238 lines) - Main client types (Client, StreamClient, StorageClient), request dispatching, URL building
- `auth.go` (13 lines) - Access key authentication header
- `errors.go` (81 lines) - APIError hierarchy (NotFoundError, AuthError, RateLimitError)
- `http.go` (10 lines) - HTTPClient interface for mockability
- `options.go` (32 lines) - Functional options (WithHTTPClient, WithUserAgent, WithBaseURL variants)
- `pagination.go` (30 lines) - PaginatedResponse[T], ListOptions structures
- `bunny_test.go` (10,542 lines) - Comprehensive test coverage

**Responsibilities:**
- Entry point for all API clients
- Request marshaling, error handling, pagination
- Authentication header injection
- HTTP client adapter layer

---

### Stream API Package `/stream`
Video library management, video CRUD, upload, captions, analytics, transcription

**Key Files:**
- `client.go` (118 lines) - StreamClient, dual-endpoint routing (api.bunny.net + video.bunnycdn.com)
- `video.go` (494 lines) - VideoService: 23 methods (List, Get, Create, Update, Delete, Upload, GetPlayer, AddCaption, GetCaption, DeleteCaption, GetStatistics, TranscribeVideo, etc.)
- `library.go` (121 lines) - LibraryService: 6 methods (List, Get, Create, Update, Delete, GetLanguages)
- `collection.go` (110 lines) - CollectionService: 5 methods (List, Get, Create, Update, Delete)
- `oembed.go` (55 lines) - OEmbedService: 1 method (Get)
- `types.go` (439 lines) - Video, Library, Collection, Statistics, Caption, Transcription types
- `client_test.go` (1,631 lines) - Full service method testing

**Design Pattern:** Service interface + implementation per resource

**Authentication:** Stream Library API Key (library-specific)

---

### Storage API Package `/storage`
Zone management (CRUD, availability, password reset), file operations (upload, download, list, delete)

**Key Files:**
- `client.go` (139 lines) - StorageClient, Zone and File service factories
- `zone.go` (137 lines) - ZoneService: 8 methods (List, Get, Create, Update, Delete, CheckNameAvailability, ResetPassword)
- `file.go` (274 lines) - FileService: 5 methods (Upload, Download, List, Delete, GetInfo)
- `types.go` (141 lines) - Zone, File types; 9 Region constants (de, ny, la, sg, syd, se, br, jh, uk)
- `client_test.go` (1,152 lines) - Comprehensive testing

**Design Pattern:** Scoped services (Zone and File services parameterized with credentials)

**Authentication:**
- Zone Management: Global API Key
- File Operations: Zone Password + Region endpoint

---

### Shield/WAF API Package `/shield`
Web application firewall, access control, rate limiting, bot detection, metrics, DDoS, promos

**Key Files:**
- `client.go` (192 lines) - Main client, 10 service accessors
- `zone-service.go` (82 lines) - ZoneService: 6 methods
- `waf-rules-service.go` (175 lines) - WAFService: 13 methods (rules CRUD, enable/disable)
- `access-list-service.go` (78 lines) - AccessListService: 6 methods (IP/Country access rules)
- `rate-limit-service.go` (68 lines) - RateLimitService: 5 methods
- `bot-detection-service.go` (42 lines) - BotDetectionService: 2 methods
- `upload-scanning-service.go` (42 lines) - UploadScanningService: 2 methods
- `ddos-service.go` (28 lines) - DDoSService: 1 method
- `promo-service.go` (28 lines) - PromoService: 1 method
- `metrics-service.go` (134 lines) - MetricsService: 7 methods (analytics queries)
- `event-logs-service.go` (57 lines) - EventLogsService: 1 method
- `types.go` (523 lines) - Comprehensive WAF/security types
- `errors.go` (35 lines) - Package-specific error handling
- `client_test.go` (1,586 lines) - Full test coverage

**Authentication:** Global API Key

---

### Edge Scripting API Package `/scripting`
Serverless edge script management, deployments, secrets, environment variables

**Key Files:**
- `client.go` (167 lines) - Main client, 5 service accessors
- `script-service.go` (133 lines) - ScriptService: 7 methods (CRUD, publish)
- `code-service.go` (38 lines) - CodeService: 2 methods (Get, Update)
- `release-service.go` (71 lines) - ReleaseService: 4 methods (list, get, publish, rollback)
- `secret-service.go` (72 lines) - SecretService: 5 methods (CRUD)
- `variable-service.go` (72 lines) - VariableService: 5 methods (CRUD)
- `types.go` (208 lines) - EdgeScript, Release, Secret, Variable types
- `errors.go` (35 lines) - Package-specific errors
- `client_test.go` (995 lines) - Comprehensive testing

**Authentication:** Global API Key

---

### Magic Containers API Package `/containers`
Container orchestration: applications, registries, volumes, endpoints, autoscaling, log forwarding

**Key Files:**
- `client.go` (234 lines) - Main client, 12 service accessors
- `application-service.go` (156 lines) - ApplicationService: 10 methods (CRUD, deploy, get overview)
- `registry-service.go` (144 lines) - RegistryService: 10 methods (image registry management)
- `container-template-service.go` (87 lines) - ContainerTemplateService: 5 methods
- `endpoint-service.go` (68 lines) - EndpointService: 4 methods
- `autoscaling-service.go` (45 lines) - AutoscalingService: 2 methods
- `region-service.go` (93 lines) - RegionService + RegionSettingsService combined
- `volume-service.go` (94 lines) - VolumeService: 5 methods
- `log-forwarding-service.go` (79 lines) - LogForwardingService: 5 methods
- `misc-service.go` (82 lines) - LimitsService, NodeService, PodService utilities
- `types.go` (819 lines) - Application, Registry, Endpoint, Volume, Region, and 20+ related types
- `errors.go` (35 lines) - Package-specific errors
- `*_test.go` (1,338 lines combined) - Comprehensive testing

**Authentication:** Global API Key

---

### Internal Package `/internal`
Shared utilities for request handling, JSON parsing, testing helpers

**Key Files:**
- `request.go` (52 lines) - NewRequest, DecodeResponse[T], ParseErrorResponse
- `bunny-flexible-time-parser.go` (49 lines) - Custom time marshaler for Bunny.net's flexible timestamp format
- `testutil/mock.go` (36 lines) - MockHTTPClient, NewMockResponse for testing

---

### Testing

**Patterns Used:**
- Table-driven tests with `testing.T` subtests
- Mock HTTP client via HTTPClient interface
- Integration test command at `cmd/test/main.go` (184 lines)

**Coverage:**
- Root package: ~1,000+ test cases
- Each service package: 200-600+ test cases per major service
- 7,300+ total test lines across codebase

---

## Package Responsibilities

| Package | Purpose | Methods | Key Type |
|---------|---------|---------|----------|
| `bunny` | Root client, routing, auth | Client, StreamClient, StorageClient | APIError |
| `stream` | Video library management | 23 video + 12 library/collection | Video, Library |
| `storage` | Zone & file ops | 8 zone + 5 file | Zone, File |
| `shield` | WAF, security, metrics | 50+ across 10 services | WAFRule, RateLimit |
| `scripting` | Edge scripts, deployment | 23 across 5 services | EdgeScript |
| `containers` | App orchestration | 40+ across 12 services | Application |
| `internal` | HTTP, parsing, testing | Request, Response, Mock | — |

---

## File Organization Conventions

### Service Files
- `{resource}-service.go`: Implements ServiceInterface
- `types.go`: All data types for package
- `client.go`: Client initialization, service accessors
- `errors.go`: Package-specific error types (if needed)
- `*_test.go`: Test files named after tested file

### Naming
- **Functions**: PascalCase (List, Get, Create, Update, Delete)
- **Types**: PascalCase (Video, Library, Zone)
- **Fields**: PascalCase with JSON tags using camelCase
- **Files**: kebab-case for multi-word files (e.g., `waf-rules-service.go`)

### Type Design
- **Pointer fields for optional fields** in PATCH/Update requests
- **Embedded structs** for shared type hierarchies (e.g., BaseVideo)
- **Custom JSON time parsing** for Bunny.net's flexible timestamp format

---

## Architecture Patterns

### 1. Service Interface Pattern
Each resource exposed as interface for mockability:

```go
type VideoService interface {
  List(ctx context.Context, opts *ListOptions) (*PaginatedResponse[Video], error)
  Get(ctx context.Context, videoID string) (*Video, error)
  Create(ctx context.Context, req *CreateVideoRequest) (*Video, error)
  // ...
}
```

### 2. Functional Options
All clients use functional options for configuration:

```go
client := NewClient("key",
  WithHTTPClient(customClient),
  WithUserAgent("my-app/1.0"),
)
```

### 3. Adapter Pattern
Bridges public interfaces to internal HTTP client:

```go
type clientAdapter struct {
  client HTTPClient
  // ...
}
```

### 4. Generic Pagination
PaginatedResponse[T] with HasMore() method

### 5. Error Hierarchy
- APIError (base) → NotFoundError, AuthError, RateLimitError
- Methods like IsAuthError(), IsRetryable() for error handling

### 6. Resource Scoping
Services parameterized with IDs (libraryID, zoneID, etc.) to scope operations

---

## Dependencies

**Zero external dependencies** - stdlib only:
- `net/http` - HTTP client
- `encoding/json` - JSON marshaling
- `context` - Context support
- `io` - Stream operations
- `fmt`, `bytes`, `strconv`, `strings` - Utilities

---

## Authentication Model

| Service | API Key Type | Location |
|---------|--------------|----------|
| Stream | Library API Key | Library settings → API tab |
| Storage Zones | Global API Key | Account Settings → API |
| Storage Files | Zone Password | Zone → FTP & API Access |
| Shield | Global API Key | Account Settings → API |
| Scripting | Global API Key | Account Settings → API |
| Containers | Global API Key | Account Settings → API |

---

## Development Workflow

### Adding a New Service

1. Create `{service}-service.go` with ServiceInterface
2. Add types to `types.go`
3. Register service in `client.go`
4. Add tests to `*_test.go`
5. Update this summary

### Testing New Code

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./stream -v
```

### Error Handling

Always check error types:

```go
if apiErr, ok := err.(*bunny.APIError); ok {
  if apiErr.IsAuthError() { /* handle auth */ }
  if apiErr.IsNotFound() { /* handle 404 */ }
  if apiErr.IsRetryable() { /* retry with backoff */ }
}
```
