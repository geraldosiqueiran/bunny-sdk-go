// Package stream provides types and services for the Bunny.net Stream API.
package stream

import "time"

// VideoState represents the processing state of a video.
type VideoState string

const (
	VideoStateCreated    VideoState = "created"
	VideoStateProcessing VideoState = "processing"
	VideoStateFinished   VideoState = "finished"
	VideoStateError      VideoState = "error"
)

// Video represents a video in a Bunny Stream library.
type Video struct {
	VideoID              string     `json:"videoId"`
	VideoLibraryID       int64      `json:"videoLibraryId"`
	Title                string     `json:"title"`
	Description          string     `json:"description,omitempty"`
	UploadDate           time.Time  `json:"uploadDate"`
	Views                int64      `json:"views"`
	Duration             int        `json:"duration"`
	Width                int        `json:"width,omitempty"`
	Height               int        `json:"height,omitempty"`
	State                VideoState `json:"state"`
	Framerate            float64    `json:"framerate,omitempty"`
	VideoCodec           string     `json:"videoCodec,omitempty"`
	AudioCodec           string     `json:"audioCodec,omitempty"`
	CollectionID         string     `json:"collectionId,omitempty"`
	ThumbnailFileName    string     `json:"thumbnailFileName,omitempty"`
	PreviewImageUrls     []string   `json:"previewImageUrls,omitempty"`
	Moments              []Moment   `json:"moments,omitempty"`
	Chapters             []Chapter  `json:"chapters,omitempty"`
	Captions             []Caption  `json:"captions,omitempty"`
	MetaTags             []MetaTag  `json:"metaTags,omitempty"`
	TranscodingMessages  []string   `json:"transcodingMessages,omitempty"`
}

// Moment represents a labeled moment in a video timeline.
type Moment struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Timestamp int    `json:"timestamp"` // milliseconds
}

// Chapter represents a chapter in a video.
type Chapter struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Start int    `json:"start"` // milliseconds
	End   int    `json:"end"`   // milliseconds
}

// Caption represents a caption track for a video.
type Caption struct {
	SrcLang string `json:"srclang"`
	Label   string `json:"label"`
}

// MetaTag represents a custom metadata tag on a video.
type MetaTag struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

// Library represents a video library in Bunny Stream.
type Library struct {
	LibraryID                int64     `json:"libraryId"`
	Name                     string    `json:"name"`
	DateCreated              time.Time `json:"dateCreated"`
	StorageUsed              int64     `json:"storageUsed"`
	StorageLimitGB           int64     `json:"storageLimitGB,omitempty"`
	VideoCacheExpirationDays int       `json:"videoCacheExpirationDays"`
	VideoCount               int       `json:"videoCount,omitempty"`
	Collections              int       `json:"collections,omitempty"`
	Region                   string    `json:"region,omitempty"`
	ReplicationRegions       []string  `json:"replicationRegions,omitempty"`
}

// Collection represents a collection of videos within a library.
type Collection struct {
	VideoLibraryID   int64    `json:"videoLibraryId"`
	GUID             string   `json:"guid"`
	Name             string   `json:"name"`
	VideoCount       int      `json:"videoCount"`
	TotalSize        int64    `json:"totalSize"`
	PreviewVideoIDs  []string `json:"previewVideoIds,omitempty"`
	PreviewImageUrls []string `json:"previewImageUrls,omitempty"`
}

// LibraryStatistics represents statistics for a video library.
type LibraryStatistics struct {
	LibraryID      int64               `json:"libraryId"`
	TotalViews     int64               `json:"totalViews"`
	TotalWatchTime int64               `json:"totalWatchTime"`
	VideoCount     int                 `json:"videoCount"`
	Bandwidth      int64               `json:"bandwidth"`
	TopVideos      []TopVideo          `json:"topVideos,omitempty"`
	ViewsByCountry []CountryViews      `json:"viewsByCountry,omitempty"`
	ViewsByDevice  map[string]int64    `json:"viewsByDevice,omitempty"`
}

// TopVideo represents a video in the top videos list.
type TopVideo struct {
	VideoID string `json:"videoId"`
	Title   string `json:"title"`
	Views   int64  `json:"views"`
}

// CountryViews represents views from a specific country.
type CountryViews struct {
	Country string `json:"country"`
	Views   int64  `json:"views"`
}

// VideoStatistics represents statistics for a single video.
type VideoStatistics struct {
	VideoID          string           `json:"videoId"`
	Views            int64            `json:"views"`
	EngagementRate   float64          `json:"engagementRate"`
	AverageWatchTime int              `json:"averageWatchTime"`
	TotalWatchTime   int64            `json:"totalWatchTime"`
	UniqueViewers    int64            `json:"uniqueViewers"`
	DeviceTypes      map[string]int64 `json:"deviceTypes,omitempty"`
	Countries        []CountryViews   `json:"countries,omitempty"`
}

// HeatmapData represents video engagement heatmap data.
type HeatmapData struct {
	VideoID     string         `json:"videoId"`
	HeatmapData []HeatmapPoint `json:"heatmapData"`
}

// HeatmapPoint represents a single point in the video heatmap.
type HeatmapPoint struct {
	Timestamp  int     `json:"timestamp"` // milliseconds
	WatchTime  int     `json:"watchTime"` // milliseconds
	Percentage float64 `json:"percentage"`
}

// PlaybackInfo represents playback information for a video.
type PlaybackInfo struct {
	VideoID       string         `json:"videoId"`
	PlaybackURL   string         `json:"playbackUrl"`
	HLSURL        string         `json:"hlsUrl"`
	DashURL       string         `json:"dashUrl"`
	CaptionTracks []CaptionTrack `json:"captionTracks,omitempty"`
}

// CaptionTrack represents a caption track with URL.
type CaptionTrack struct {
	SrcLang string `json:"srclang"`
	Label   string `json:"label"`
	URL     string `json:"url"`
}

// CreateVideoRequest represents a request to create a new video.
type CreateVideoRequest struct {
	Title         string `json:"title"`
	CollectionID  string `json:"collectionId,omitempty"`
	ThumbnailTime int    `json:"thumbnailTime,omitempty"` // milliseconds
}

// UpdateVideoRequest represents a request to update video metadata.
type UpdateVideoRequest struct {
	Title        string `json:"title,omitempty"`
	CollectionID string `json:"collectionId,omitempty"`
}

// FetchVideoRequest represents a request to fetch a video from URL.
type FetchVideoRequest struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

// FetchVideoResponse represents the response from fetching a video.
type FetchVideoResponse struct {
	VideoID    string     `json:"videoId"`
	State      VideoState `json:"state"`
	UploadDate time.Time  `json:"uploadDate"`
}

// ReencodeRequest represents a request to re-encode a video.
type ReencodeRequest struct {
	Resolutions []Resolution `json:"resolutions,omitempty"`
}

// Resolution represents a video resolution configuration.
type Resolution struct {
	Resolution int `json:"resolution"` // height in pixels (e.g., 720, 1080)
	Bitrate    int `json:"bitrate"`    // kbps
}

// AddCaptionRequest represents a request to add a caption to a video.
type AddCaptionRequest struct {
	SrcLang      string `json:"srclang"`      // ISO 639-1 language code
	Label        string `json:"label"`        // Display name
	CaptionsFile string `json:"captionsFile"` // Base64-encoded VTT/SRT content
}

// SetThumbnailRequest represents a request to set video thumbnail.
type SetThumbnailRequest struct {
	ThumbnailTime int `json:"thumbnailTime"` // milliseconds
}

// SetThumbnailResponse represents the response from setting a thumbnail.
type SetThumbnailResponse struct {
	VideoID           string `json:"videoId"`
	ThumbnailFileName string `json:"thumbnailFileName"`
	ThumbnailURL      string `json:"thumbnailUrl"`
}

// CreateLibraryRequest represents a request to create a library.
type CreateLibraryRequest struct {
	Name                     string `json:"name"`
	VideoCacheExpirationDays int    `json:"videoCacheExpirationDays,omitempty"`
}

// UpdateLibraryRequest represents a request to update a library.
type UpdateLibraryRequest struct {
	Name                     string `json:"name,omitempty"`
	VideoCacheExpirationDays *int   `json:"videoCacheExpirationDays,omitempty"`
}

// CreateCollectionRequest represents a request to create a collection.
type CreateCollectionRequest struct {
	Name string `json:"name"`
}

// UpdateCollectionRequest represents a request to update a collection.
type UpdateCollectionRequest struct {
	Name string `json:"name"`
}

// VideoListOptions specifies options for listing videos.
type VideoListOptions struct {
	Page              int    `url:"page,omitempty"`
	ItemsPerPage      int    `url:"itemsPerPage,omitempty"`
	Search            string `url:"search,omitempty"`
	Collection        string `url:"collection,omitempty"`
	OrderBy           string `url:"orderBy,omitempty"`
	IncludeThumbnails bool   `url:"includeThumbnails,omitempty"`
}

// CollectionListOptions specifies options for listing collections.
type CollectionListOptions struct {
	Page              int    `url:"page,omitempty"`
	ItemsPerPage      int    `url:"itemsPerPage,omitempty"`
	OrderBy           string `url:"orderBy,omitempty"`
	IncludeThumbnails bool   `url:"includeThumbnails,omitempty"`
}

// LibraryListOptions specifies options for listing libraries.
type LibraryListOptions struct {
	Page         int `url:"page,omitempty"`
	ItemsPerPage int `url:"itemsPerPage,omitempty"`
}

// StatisticsOptions specifies date range options for statistics.
type StatisticsOptions struct {
	DateFrom string `url:"dateFrom,omitempty"` // ISO 8601 date
	DateTo   string `url:"dateTo,omitempty"`   // ISO 8601 date
}
