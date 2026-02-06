# Bunny.net Go SDK

Idiomatic Go SDK for [Bunny.net](https://bunny.net) Stream, Storage, Shield/WAF, and Magic Containers APIs. **Zero external dependencies**, pure stdlib, Go 1.23+.

## Installation

```bash
go get github.com/geraldo/bunny-sdk-go
```

## Features

- **Stream API** (23 methods): Videos, libraries, collections, captions, analytics
- **Storage API** (13 methods): Zone management, file uploads/downloads, 9 regions
- **Shield/WAF API** (50+ methods): WAF rules, access lists, rate limiting, bot detection, metrics
- **Edge Scripting API** (23 methods): Scripts, deployments, secrets, variables
- **Magic Containers API** (40+ methods): Applications, registries, volumes, endpoints, autoscaling
- Zero external dependencies (stdlib only)
- Interface-based design for easy testing
- Context support on all operations
- Functional options pattern
- Streaming file uploads/downloads (no memory buffering)
- Typed error hierarchy (NotFoundError, AuthError, RateLimitError)

## Quick Start

### Stream API

```go
client := stream.NewClient(os.Getenv("BUNNY_STREAM_API_KEY"))
videos, err := client.Videos(12345).List(context.Background(), nil)

video, err := client.Videos(12345).Create(context.Background(), &stream.CreateVideoRequest{
    Title: "My Video",
})

f, _ := os.Open("video.mp4")
defer f.Close()
err = client.Videos(12345).Upload(context.Background(), video.VideoID, f)
```

### Storage API

```go
// Zone management
client := storage.NewClient(os.Getenv("BUNNY_API_KEY"))
zones, err := client.Zones().List(context.Background(), nil)

// File operations
fs := storage.NewFileService("zone-name", "password", storage.RegionFalkenstein)
err = fs.Upload(context.Background(), "path/to/file", fileReader, nil)
files, err := fs.List(context.Background(), "directory")
```

### Shield/WAF API

```go
client := shield.NewClient(os.Getenv("BUNNY_API_KEY"))
zones, err := client.Zones().List(context.Background())

rateLimit, err := client.RateLimits().Create(context.Background(), &shield.CreateRateLimitRequest{
    Name:              "API Rate Limit",
    RequestsPerSecond: 100,
    ShieldZoneID:      zones.Items[0].ID,
})

metrics, err := client.Metrics().GetOverview(context.Background(), &shield.DateRangeOptions{
    From: "2026-02-01",
    To:   "2026-02-04",
})
```

### Edge Scripting API

```go
client := scripting.NewClient(os.Getenv("BUNNY_API_KEY"))
scripts, err := client.Scripts().List(context.Background(), nil)

script, err := client.Scripts().Create(context.Background(), &scripting.CreateScriptRequest{
    Name: "my-script",
})
```

### Magic Containers API

```go
client := containers.NewClient(os.Getenv("BUNNY_API_KEY"))
apps, err := client.Applications().List(context.Background(), nil)

app, err := client.Applications().Create(context.Background(), &containers.CreateApplicationRequest{
    Name:        "my-app",
    RuntimeType: containers.RuntimeTypeShared,
})

err = client.Applications().Deploy(context.Background(), app.ID)
```

## Authentication

| Service | API Key Type | Where to Find |
|---------|--------------|---------------|
| Stream API | Stream Library API Key | Library settings → API tab |
| Storage Zones | Global API Key | Account Settings → API |
| File Operations | Storage Zone Password | Zone → FTP & API Access |
| Shield/WAF | Global API Key | Account Settings → API |
| Scripting | Global API Key | Account Settings → API |
| Containers | Global API Key | Account Settings → API |

## Storage Regions

```
RegionFalkenstein  "storage"  Germany (default)
RegionNewYork      "ny"       New York
RegionLosAngeles   "la"       Los Angeles
RegionSingapore    "sg"       Singapore
RegionSydney       "syd"      Sydney
RegionStockholm    "se"       Stockholm
RegionSaoPaulo     "br"       Sao Paulo
RegionJohannesburg "jh"       Johannesburg
RegionLondon       "uk"       London
```

## Error Handling

```go
video, err := client.Videos(123).Get(ctx, "video-id")
if err != nil {
    if apiErr, ok := err.(*stream.APIError); ok {
        if apiErr.IsNotFound() {
            // Handle 404
        } else if apiErr.IsAuthError() {
            // Handle 401/403
        } else if apiErr.IsRetryable() {
            // Implement retry logic
        }
    }
}
```

## Testing

Mock the HTTP client for unit tests:

```go
mock := &testutil.MockHTTPClient{
    DoFunc: func(req *http.Request) (*http.Response, error) {
        return testutil.NewMockResponse(200, `{"videoId":"test"}`), nil
    },
}
client := stream.NewClient("key", stream.WithHTTPClient(mock))
```

## Documentation

- **[Codebase Summary](./docs/codebase-summary.md)** - Directory structure, packages, architecture patterns
- **[Code Standards](./docs/code-standards.md)** - Naming conventions, style guidelines, testing patterns
- **[System Architecture](./docs/system-architecture.md)** - Request flow, error handling, service design
- **[Project Overview & PDR](./docs/project-overview-pdr.md)** - Goals, features, requirements, success criteria
- **[Project Roadmap](./docs/project-roadmap.md)** - Milestones, timeline, planned features

## License

MIT
