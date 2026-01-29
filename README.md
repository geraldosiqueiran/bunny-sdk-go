# Bunny.net Go SDK

Idiomatic Go SDK for [Bunny.net](https://bunny.net) Stream and Storage APIs.

## Installation

```bash
go get github.com/geraldo/bunny-sdk-go
```

## Features

- **Stream API**: Videos, Libraries, Collections management
- **Storage API**: Zone management and file operations
- Interface-based design for easy mocking
- Context support on all operations
- Functional options pattern
- Streaming file uploads/downloads (no memory buffering)
- Typed errors with status codes

## Quick Start

### Stream API

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/geraldo/bunny-sdk-go/stream"
)

func main() {
    // Use Stream Library API Key (found in library settings)
    client := stream.NewClient(os.Getenv("BUNNY_STREAM_API_KEY"))
    ctx := context.Background()

    // List videos in a library
    videos, err := client.Videos(12345).List(ctx, nil)
    if err != nil {
        panic(err)
    }

    for _, v := range videos.Items {
        fmt.Printf("Video: %s (%s)\n", v.Title, v.VideoID)
    }

    // Create a video
    video, err := client.Videos(12345).Create(ctx, &stream.CreateVideoRequest{
        Title: "My New Video",
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created: %s\n", video.VideoID)

    // Upload video file
    f, _ := os.Open("video.mp4")
    defer f.Close()

    err = client.Videos(12345).Upload(ctx, video.VideoID, f)
    if err != nil {
        panic(err)
    }
}
```

### Storage API - Zone Management

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/geraldo/bunny-sdk-go/storage"
)

func main() {
    // Use Global API Key for zone management
    client := storage.NewClient(os.Getenv("BUNNY_API_KEY"))
    ctx := context.Background()

    // List storage zones
    zones, err := client.Zones().List(ctx, nil)
    if err != nil {
        panic(err)
    }

    for _, z := range zones.Items {
        fmt.Printf("Zone: %s (Region: %s)\n", z.Name, z.Region)
    }

    // Create a storage zone
    zone, err := client.Zones().Create(ctx, &storage.CreateZoneRequest{
        Name:   "my-new-zone",
        Region: "de",
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created zone: %s (Password: %s)\n", zone.Name, zone.Password)
}
```

### Storage API - File Operations

```go
package main

import (
    "context"
    "fmt"
    "io"
    "os"

    "github.com/geraldo/bunny-sdk-go/storage"
)

func main() {
    // Use Storage Zone Password for file operations
    fs := storage.NewFileService(
        "my-zone",                    // zone name
        os.Getenv("BUNNY_ZONE_PASS"), // zone password
        storage.RegionFalkenstein,    // region
    )
    ctx := context.Background()

    // Upload a file
    f, _ := os.Open("document.pdf")
    defer f.Close()

    err := fs.Upload(ctx, "documents/report.pdf", f, nil)
    if err != nil {
        panic(err)
    }

    // List directory
    files, err := fs.List(ctx, "documents")
    if err != nil {
        panic(err)
    }

    for _, file := range files {
        fmt.Printf("%s (%d bytes)\n", file.ObjectName, file.Length)
    }

    // Download a file
    reader, err := fs.Download(ctx, "documents/report.pdf")
    if err != nil {
        panic(err)
    }
    defer reader.Close()

    out, _ := os.Create("downloaded.pdf")
    defer out.Close()
    io.Copy(out, reader)
}
```

## Authentication

Bunny.net uses different API keys for different services:

| Service | API Key Type | Where to Find |
|---------|--------------|---------------|
| Stream API | Stream Library API Key | Stream library → API tab |
| Storage Zone Management | Global API Key | Account Settings → API |
| File Operations | Storage Zone Password | Storage zone → FTP & API Access |

## Available Regions

For Storage file operations:

| Region | Constant | Location |
|--------|----------|----------|
| `de` | `RegionFalkenstein` | Germany (default) |
| `ny` | `RegionNewYork` | New York |
| `la` | `RegionLosAngeles` | Los Angeles |
| `sg` | `RegionSingapore` | Singapore |
| `syd` | `RegionSydney` | Sydney |
| `se` | `RegionStockholm` | Stockholm |
| `br` | `RegionSaoPaulo` | Sao Paulo |
| `jh` | `RegionJohannesburg` | Johannesburg |
| `uk` | `RegionLondon` | London |

## Error Handling

All errors include the HTTP status code:

```go
video, err := client.Videos(123).Get(ctx, "nonexistent")
if err != nil {
    // Check error type
    if apiErr, ok := err.(*stream.APIError); ok {
        if apiErr.StatusCode == 404 {
            fmt.Println("Video not found")
        }
    }
}
```

## Testing

The SDK uses interfaces for HTTP clients, making it easy to mock:

```go
mock := &testutil.MockHTTPClient{
    DoFunc: func(req *http.Request) (*http.Response, error) {
        return testutil.NewMockResponse(200, `{"videoId":"test"}`), nil
    },
}

client := stream.NewClient("key", stream.WithHTTPClient(mock))
```

## License

MIT
