# Project Overview & Product Development Requirements

## Project Overview

**Project Name:** Bunny.net Go SDK

**Module:** `github.com/geraldo/bunny-sdk-go`

**Purpose:** Provide idiomatic, zero-dependency Go SDK for Bunny.net cloud platform APIs, enabling Go developers to seamlessly integrate video streaming, edge storage, security (WAF), edge scripting, and containerized applications.

**Target Audience:** Go developers building applications that leverage Bunny.net services.

---

## Project Goals

1. **Idiomatic Go Design** - Leverage Go conventions, generics, interfaces
2. **Zero Dependencies** - Stdlib only for minimal footprint
3. **Comprehensive API Coverage** - Support all major Bunny.net APIs
4. **Developer Experience** - Easy setup, clear error handling, mockable interfaces
5. **Type Safety** - Leverage Go's type system with generics for pagination, options
6. **Production Ready** - Thorough testing, proper error hierarchies, context support

---

## Feature Set

### 1. Stream API (Video Management)
- **Libraries**: List, create, update, delete video libraries
- **Videos**: Full CRUD operations, upload/streaming, player configuration
- **Captions**: Add, retrieve, delete video captions
- **Analytics**: Video statistics and metrics
- **Transcription**: Automated video transcription support
- **Collections**: Video collections management

**Methods:** 23+ video operations

### 2. Storage API (Edge Storage)
**Zone Management:**
- Zone CRUD operations
- Zone availability checking
- Password reset
- Methods: 8

**File Operations:**
- Upload (streaming, no memory buffering)
- Download (streaming)
- List directory contents
- Delete files/directories
- Get file info
- Methods: 5
- Regions: 9 (Falkenstein, New York, Los Angeles, Singapore, Sydney, Stockholm, Sao Paulo, Johannesburg, London)

### 3. Shield API (WAF & Security)
- **WAF Rules**: 13 methods (CRUD, enable/disable, rule management)
- **Access Lists**: IP/Country-based access control (6 methods)
- **Rate Limiting**: Request rate limiting policies (5 methods)
- **Bot Detection**: Bot detection configuration (2 methods)
- **Upload Scanning**: Malware/file scanning (2 methods)
- **DDoS Protection**: DDoS metrics (1 method)
- **Metrics**: Security metrics and analytics (7 methods)
- **Event Logs**: Security event logging (1 method)
- **Zone Management**: Shield zone configuration (6 methods)
- **Promo Features**: Additional security features (1 method)

**Methods:** 50+

### 4. Edge Scripting API
- **Scripts**: Edge script CRUD and deployment (7 methods)
- **Code**: Script code management (2 methods)
- **Releases**: Script version releases (4 methods)
- **Secrets**: Encrypted secret management (5 methods)
- **Variables**: Environment variables (5 methods)

**Methods:** 23

### 5. Magic Containers API
- **Applications**: Container app management (10 methods)
- **Registries**: Image registry configuration (10 methods)
- **Container Templates**: Template management (5 methods)
- **Endpoints**: Application endpoint configuration (4 methods)
- **Autoscaling**: Auto-scaling policies (2 methods)
- **Regions**: Available deployment regions (varies)
- **Volumes**: Persistent volume management (5 methods)
- **Log Forwarding**: Logging configuration (5 methods)
- **Utilities**: Limits, nodes, pods info (3 services)

**Methods:** 40+

---

## Non-Functional Requirements

### Performance
- Context support on all operations for timeouts, cancellation
- Streaming file uploads/downloads (no memory buffering)
- Minimal allocation strategy for performance-critical paths

### Compatibility
- **Go Version:** 1.23+
- **Generics:** Heavy use of Go 1.18+ generics for type safety
- **Cross-platform:** Linux, macOS, Windows

### Dependencies
- **Zero external dependencies** - stdlib only
- Reduces attack surface, deployment complexity, version conflicts

### Error Handling
- Comprehensive error types: APIError, NotFoundError, AuthError, RateLimitError
- Error helper methods: IsAuthError(), IsNotFound(), IsRetryable()
- Status codes included in all errors

### Testing
- Table-driven tests
- Mock HTTP client support via HTTPClient interface
- 7,300+ lines of test code

### Authentication
- Support for multiple API key types per service
- Stream Library API Key (library-specific)
- Global API Key (account-wide for management APIs)
- Zone passwords for file operations

---

## API Coverage Status

| API | Status | Completeness | Methods |
|-----|--------|--------------|---------|
| Stream | Implemented | 95% | 23+ |
| Storage Zones | Implemented | 100% | 8 |
| Storage Files | Implemented | 100% | 5 |
| Shield/WAF | Implemented | 95% | 50+ |
| Edge Scripting | Implemented | 100% | 23 |
| Containers | Implemented | 90% | 40+ |

**Not Yet Implemented:**
- Pull Zones API
- DNS Zone API
- Billing API
- Abuse reporting API

---

## Success Criteria

### Functional Success
- All 5 major API packages fully functional
- 100% pass rate on test suite
- All documented examples working
- No breaking changes in v1.x

### Quality Success
- Code coverage: 80%+ across codebase
- Zero external dependencies maintained
- All errors typed and distinguishable
- Comprehensive error messages

### Developer Success
- Clear, concise API documentation
- Runnable examples per package
- Intuitive interface hierarchy (client → services)
- Mock-friendly design for testing

### Production Readiness
- Context support throughout
- Proper timeout/cancellation handling
- Rate limit error detection
- Retryable error identification
- Proper resource cleanup (resp.Body.Close)

---

## Architecture Highlights

### Client Hierarchy
```
NewClient(globalAPIKey)           → Storage/Shield/Scripting/Containers
NewStreamClient(libraryAPIKey)    → Stream API
NewStorageClient(zoneName, pwd)   → File Operations
```

### Request Flow
1. Client method called with context
2. Request marshaled to JSON
3. Auth header injected
4. HTTP request dispatched
5. Response parsed or error returned
6. Error type hierarchy applied

### Error Handling
1. HTTP status code checked
2. Error body parsed (with optional error key/field)
3. Specific error type created (NotFound, Auth, RateLimit, or generic)
4. Helper methods available for checking error type

---

## Development Phases

### Phase 1: Foundation (Complete)
- Root package with client types
- Authentication patterns
- Error hierarchy and handling
- HTTP adapter layer
- Pagination support

### Phase 2: Stream API (Complete)
- Library management
- Video CRUD and metadata
- Upload/streaming support
- Captions and analytics
- Collections support

### Phase 3: Storage API (Complete)
- Zone management (CRUD, availability, passwords)
- File operations (upload, download, list, delete)
- Multi-region support
- Streaming file handling

### Phase 4: Shield/WAF API (Complete)
- Zone configuration
- WAF rule management
- Access list control
- Rate limiting
- Bot detection, DDoS, metrics

### Phase 5: Edge Scripting API (Complete)
- Script management and deployment
- Code, releases, secrets, variables
- Full lifecycle support

### Phase 6: Containers API (Complete)
- Application management
- Registry and container templates
- Endpoints and autoscaling
- Volumes and log forwarding
- Region management

### Phase 7: Polish & Optimization (In Progress)
- Documentation completion
- Additional examples
- Performance profiling
- Extended test coverage

---

## Testing Strategy

### Unit Testing
- Service method testing with mock HTTP client
- Error condition testing for each method
- Pagination testing

### Integration Testing
- `cmd/test/main.go` for live API testing
- All services exercised

### Mock Testing
- HTTPClient interface for full control
- NewMockResponse for response simulation

---

## Documentation Structure

- `README.md` - Quick start, authentication, examples
- `codebase-summary.md` - Architecture, directory structure, patterns
- `code-standards.md` - Naming conventions, patterns, guidelines
- `system-architecture.md` - System design, flows, design decisions
- `project-roadmap.md` - Milestones, planned improvements

---

## Known Limitations & Future Work

### Known Limitations
- No built-in retry logic (applications can implement via context)
- No connection pooling configuration (uses default http.DefaultClient)
- No rate limit backoff handling (errors indicate rate limit)

### Future Improvements
- Optional retry middleware
- Connection pool configuration options
- Automatic rate limit backoff
- Additional API coverage (Pull Zones, DNS, Billing)
- Performance optimization benchmarks

---

## Security & Compliance

### API Key Handling
- API keys not logged or exposed in error messages
- Authentication headers set via dedicated auth function
- No sensitive data in error strings

### Data Protection
- HTTPS only (all endpoints use https://)
- TLS 1.2+ enforced by Go's http package
- No plaintext transmission of credentials

### Error Information
- Error messages safe for user display
- No stack traces or internals exposed
- Field/ErrorKey hints for validation errors

---

## Version & Release

**Current Version:** Implementing v1.0

**Versioning:** Semantic versioning (major.minor.patch)

**Stability:** Beta - API stable but still accepting feedback

---

## Contact & Support

- Repository: `github.com/geraldo/bunny-sdk-go`
- Issues: GitHub issue tracker
- Bunny.net Support: https://support.bunny.net

