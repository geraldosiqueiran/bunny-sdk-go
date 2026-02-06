# Project Roadmap

## Vision

Become the **standard idiomatic Go SDK** for Bunny.net cloud platform, providing comprehensive API coverage with zero dependencies, excellent developer experience, and production-ready reliability.

---

## Release Timeline

### v0.1 - Alpha (In Development)
**Status:** 70% Complete

Target: Foundation and core APIs operational

**Completed:**
- ✅ Root package with Client types
- ✅ Authentication patterns (Bearer token)
- ✅ Error hierarchy and typed errors
- ✅ Pagination support with generics
- ✅ Stream API (23 methods)
- ✅ Storage Zones API (8 methods)
- ✅ Storage Files API (5 methods, streaming)
- ✅ Shield/WAF API (50+ methods)
- ✅ Edge Scripting API (23 methods)
- ✅ Magic Containers API (40+ methods)
- ✅ Comprehensive test suite (7,300+ lines)
- ✅ Mock HTTP client for testing

**In Progress:**
- 🔄 Documentation completion (codebase summary, standards, architecture)
- 🔄 Additional usage examples
- 🔄 Integration test enhancements

**Not Yet Started:**
- ❌ Performance benchmarks
- ❌ Extended error handling patterns
- ❌ Retryable error documentation

---

### v1.0 - Stable Release
**Status:** Planning

Target: Production-ready, stable API

**Requirements:**
- 95%+ test coverage
- Comprehensive documentation
- All core APIs implemented
- Performance baseline established
- v0.x feedback addressed
- Breaking change: none expected

**Timeline:** Q2 2026

---

## Feature Phases

### Phase 1: Foundation ✅ (Complete)
**Completion:** 100%

**Features:**
- Client types (Client, StreamClient, StorageClient)
- HTTPClient interface for mockability
- Functional options pattern
- Error hierarchy (APIError, NotFoundError, AuthError, RateLimitError)
- Pagination with generics
- Authentication header injection

**Metrics:**
- Core package: 200+ test cases
- Error handling: 100% of status codes covered
- Options: All 4 variants tested

---

### Phase 2: Stream API ✅ (Complete)
**Completion:** 95%

**Methods Implemented:** 23+
- Library CRUD (6)
- Video CRUD (7)
- Video Upload/Streaming (2)
- Captions (4)
- Analytics & Statistics (2)
- Transcription (1)
- OEmbed (1)
- Collections (5)

**Missing (5%):**
- Advanced search/filtering documentation
- Captions batch operations

**Metrics:**
- Test coverage: 95%
- Test cases: 1,600+
- Methods tested: All main operations

---

### Phase 3: Storage API ✅ (Complete)
**Completion:** 100%

**Zone Management (8 methods):**
- ✅ List zones with pagination
- ✅ Get zone details
- ✅ Create zone
- ✅ Update zone configuration
- ✅ Delete zone
- ✅ Check name availability
- ✅ Reset zone password
- ✅ Get zone statistics

**File Operations (5 methods):**
- ✅ Upload file (streaming)
- ✅ Download file (streaming)
- ✅ List directory contents
- ✅ Delete file/directory
- ✅ Get file information

**Regions:** All 9 regions supported

**Metrics:**
- Zone methods: 8/8 (100%)
- File methods: 5/5 (100%)
- Test coverage: 98%
- Test cases: 1,100+

---

### Phase 4: Shield/WAF API ✅ (Complete)
**Completion:** 95%

**Core Services Implemented:**
- ✅ WAF Rules (13 methods)
- ✅ Access Lists (6 methods)
- ✅ Rate Limiting (5 methods)
- ✅ Bot Detection (2 methods)
- ✅ DDoS Protection (1 method)
- ✅ Metrics (7 methods)
- ✅ Event Logs (1 method)
- ✅ Upload Scanning (2 methods)
- ✅ Promo Features (1 method)
- ✅ Zone Management (6 methods)

**Missing (5%):**
- Advanced metrics dashboarding documentation
- Custom rule template examples

**Metrics:**
- Total services: 10
- Total methods: 50+
- Test coverage: 94%
- Test cases: 1,500+

---

### Phase 5: Edge Scripting API ✅ (Complete)
**Completion:** 100%

**Features:**
- ✅ Script CRUD (7 methods)
- ✅ Code management (2 methods)
- ✅ Release management (4 methods)
- ✅ Secrets management (5 methods)
- ✅ Environment variables (5 methods)

**Metrics:**
- Methods: 23/23 (100%)
- Test coverage: 100%
- Test cases: 990+

---

### Phase 6: Magic Containers API ✅ (Complete)
**Completion:** 90%

**Services Implemented:**
- ✅ Applications (10 methods)
- ✅ Registries (10 methods)
- ✅ Container Templates (5 methods)
- ✅ Endpoints (4 methods)
- ✅ Autoscaling (2 methods)
- ✅ Regions (varies)
- ✅ Volumes (5 methods)
- ✅ Log Forwarding (5 methods)
- ✅ Utilities - Limits (1)
- ✅ Utilities - Nodes (1)
- ✅ Utilities - Pods (1)

**Missing (10%):**
- Pod log streaming documentation
- Advanced deployment templates

**Metrics:**
- Services: 12 total
- Methods: 40+ implemented
- Test coverage: 88%
- Test cases: 1,300+

---

### Phase 7: Documentation & Polish 🔄 (In Progress)
**Target Completion:** Feb 2026

**Deliverables:**
- 📝 codebase-summary.md (Architecture, directory structure)
- 📝 code-standards.md (Naming, patterns, guidelines)
- 📝 system-architecture.md (Design, flows, decisions)
- 📝 project-overview-pdr.md (Goals, requirements, success criteria)
- 📝 project-roadmap.md (Milestones, timeline)
- 📝 README.md updates (Quick start, authentication table)
- 🔄 Example programs for each package
- 🔄 Troubleshooting guide
- 🔄 API coverage matrix

---

## Planned Features (Post-v1.0)

### Additional APIs
- **Pull Zones API** - CDN pull zone management
  - **Effort:** Medium
  - **Timeline:** Q3 2026
  - **Methods:** ~15 estimated

- **DNS API** - DNS zone and record management
  - **Effort:** Medium
  - **Timeline:** Q3 2026
  - **Methods:** ~20 estimated

- **Billing API** - Invoice and usage tracking
  - **Effort:** Low
  - **Timeline:** Q4 2026
  - **Methods:** ~8 estimated

- **Abuse API** - Abuse reporting and management
  - **Effort:** Low
  - **Timeline:** Q4 2026
  - **Methods:** ~5 estimated

### Performance & Reliability
- **Automatic Retry Logic**
  - ✋ Not implemented (applications can implement)
  - Planned with exponential backoff
  - Rate limit aware (backoff on 429)

- **Rate Limit Handling**
  - Current: Returns RateLimitError on 429
  - Planned: Optional auto-retry with backoff
  - Parse rate limit headers for better decisions

- **Connection Pooling Configuration**
  - Current: Uses http.DefaultClient
  - Planned: Configurable via Options
  - Control max idle connections, timeouts

- **Request Metrics**
  - Track request/response times
  - Error rate monitoring
  - Per-endpoint statistics

### Developer Experience
- **CLI Tool**
  - bunny-sdk-go-cli for quick API testing
  - Authentication helpers
  - Response pretty-printing

- **Code Generators**
  - Generate typed clients from OpenAPI
  - Auto-generate mock HTTP clients

- **Middleware Support**
  - Logging middleware
  - Metrics collection
  - Request/response tracing

### Testing Enhancements
- **Integration Test Suite**
  - Live API testing (opt-in with API key)
  - All services end-to-end
  - Performance baselines

- **Benchmark Suite**
  - Memory allocations
  - Throughput measurements
  - Garbage collection impact

---

## Milestone Targets

### Stability Milestones

| Milestone | Target Date | Criteria |
|-----------|-------------|----------|
| **Alpha** | ✅ Complete | Core APIs working, basic docs |
| **Beta** | Feb 28, 2026 | Full documentation, examples |
| **RC1** | Mar 15, 2026 | Bug fixes, performance tuning |
| **v1.0** | Apr 30, 2026 | Stable API, production ready |

### Quality Milestones

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| **Test Coverage** | 90% | 85% | 🟡 In Progress |
| **API Completeness** | 100% | 95% | 🟡 In Progress |
| **Documentation** | 100% | 70% | 🟡 In Progress |
| **Zero Dependencies** | Maintained | ✅ Yes | 🟢 Complete |

---

## Known Issues & Improvements

### Current Limitations

1. **No Built-in Retry Logic**
   - Status: By design
   - Reason: Let applications handle retry strategy
   - Workaround: Check IsRetryable() and implement backoff
   - Timeline: Optional middleware in v1.1

2. **No Rate Limit Backoff**
   - Status: By design
   - Reason: RateLimitError returned immediately
   - Workaround: Check IsRateLimited() and implement delay
   - Timeline: Optional feature in v1.1

3. **No Connection Pool Configuration**
   - Status: Uses defaults
   - Reason: http.DefaultClient works for most cases
   - Workaround: Pass custom http.Client via WithHTTPClient
   - Timeline: Configurable options in v1.1

### Documentation Gaps

| Gap | Impact | Timeline |
|-----|--------|----------|
| Comprehensive examples | Medium | Feb 2026 |
| Troubleshooting guide | Medium | Feb 2026 |
| API coverage matrix | Low | Feb 2026 |
| Performance guidelines | Low | Mar 2026 |

---

## Success Metrics

### Adoption
- [ ] 100+ GitHub stars
- [ ] 10+ external projects using SDK
- [ ] Featured in Bunny.net documentation

### Quality
- [ ] 90%+ test coverage
- [ ] 0 critical bugs in v1.0
- [ ] <1 day response time for issues

### Performance
- [ ] <100ms latency for list operations
- [ ] <1s latency for upload/download
- [ ] <10MB memory for typical workload

### Documentation
- [ ] All methods documented
- [ ] Example for each service
- [ ] Clear error handling guide

---

## Community & Feedback

### v0.x Feedback Priorities
1. API usability (confirm interface design)
2. Error handling ergonomics
3. Performance requirements
4. Missing features

### How to Contribute
- Report bugs via GitHub issues
- Suggest features with use cases
- Submit PRs (with tests)
- Improve documentation

---

## Long-Term Vision (2027+)

### Ecosystem
- Official Bunny.net SDK
- Support in other languages (Python, Node.js, Rust)
- Integration examples with popular frameworks

### Features
- GraphQL client option
- Streaming WebSocket support
- Advanced analytics integration

### Developer Tools
- VS Code extension
- IntelliJ IDE plugin
- API explorer web UI

---

## Maintenance Plan

### Update Frequency
- **Critical bugs:** Immediate hotfix
- **Regular maintenance:** Monthly releases
- **Feature releases:** Quarterly
- **Major versions:** As needed for major features

### Support Timeline
- **v1.x:** 2+ years active support
- **v0.x:** Support until v1.0 released

### Deprecation Policy
- Announce 2 releases in advance
- Provide migration guide
- Support old version 1 release after deprecation

---

## Next Steps

### Immediate (This Sprint)
- [ ] Complete documentation (this week)
- [ ] Add comprehensive examples
- [ ] Create API coverage matrix
- [ ] Collect feedback on v0.1

### Short-term (Next Sprint)
- [ ] Resolve documentation gaps
- [ ] Add integration test suite
- [ ] Performance benchmarking
- [ ] Release v0.2 beta

### Medium-term (Q2 2026)
- [ ] Finalize v1.0 API
- [ ] Address all v0.x feedback
- [ ] Comprehensive guide updates
- [ ] Release v1.0 stable

---

## Related Documents
- Project Overview & PDR: `project-overview-pdr.md`
- Code Standards: `code-standards.md`
- System Architecture: `system-architecture.md`
- Codebase Summary: `codebase-summary.md`

