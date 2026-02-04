// Package stream provides types and services for the Bunny.net Stream API.
package stream

import "github.com/geraldo/bunny-sdk-go/internal"

// VideoState represents the processing state of a video.
type VideoState string

const (
	VideoStateCreated    VideoState = "created"
	VideoStateProcessing VideoState = "processing"
	VideoStateFinished   VideoState = "finished"
	VideoStateError      VideoState = "error"
)

// OutputCodec represents video output codec options.
type OutputCodec int

const (
	OutputCodecX264 OutputCodec = 0
	OutputCodecVP9  OutputCodec = 1
	OutputCodecHEVC OutputCodec = 2
	OutputCodecAV1  OutputCodec = 3
)

// Video represents a video in a Bunny Stream library.
type Video struct {
	VideoID              string     `json:"videoId"`
	VideoLibraryID       int64      `json:"videoLibraryId"`
	Title                string     `json:"title"`
	Description          string     `json:"description,omitempty"`
	UploadDate           internal.BunnyTime  `json:"uploadDate"`
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
	DateCreated              internal.BunnyTime `json:"dateCreated"`
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
	UploadDate internal.BunnyTime  `json:"uploadDate"`
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

// CleanupResolutionsOptions specifies options for cleanup operation.
type CleanupResolutionsOptions struct {
	ResolutionsToDelete            string `url:"resolutionsToDelete,omitempty"`
	DeleteNonConfiguredResolutions bool   `url:"deleteNonConfiguredResolutions,omitempty"`
	DeleteOriginal                 bool   `url:"deleteOriginal,omitempty"`
	DeleteMp4Files                 bool   `url:"deleteMp4Files,omitempty"`
	DryRun                         bool   `url:"dryRun,omitempty"`
}

// StatusResponse represents a generic status response.
type StatusResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode"`
}

// HeatmapDataOptions specifies options for heatmap data request.
type HeatmapDataOptions struct {
	Token   string `url:"token,omitempty"`
	Expires int64  `url:"expires,omitempty"`
}

// VideoPlayData represents video playback data with heatmap.
type VideoPlayData struct {
	Video                *Video  `json:"video,omitempty"`
	LibraryName          string  `json:"libraryName,omitempty"`
	CaptionsPath         string  `json:"captionsPath,omitempty"`
	SeekPath             string  `json:"seekPath,omitempty"`
	ThumbnailURL         string  `json:"thumbnailUrl,omitempty"`
	FallbackURL          string  `json:"fallbackUrl,omitempty"`
	VideoPlaylistURL     string  `json:"videoPlaylistUrl,omitempty"`
	OriginalURL          string  `json:"originalUrl,omitempty"`
	PreviewURL           string  `json:"previewUrl,omitempty"`
	Controls             string  `json:"controls,omitempty"`
	EnableDRM            bool    `json:"enableDRM"`
	DRMVersion           int     `json:"drmVersion"`
	PlayerKeyColor       string  `json:"playerKeyColor,omitempty"`
	VastTagURL           string  `json:"vastTagUrl,omitempty"`
	CaptionsFontSize     int     `json:"captionsFontSize"`
	CaptionsFontColor    string  `json:"captionsFontColor,omitempty"`
	CaptionsBackground   string  `json:"captionsBackground,omitempty"`
	UILanguage           string  `json:"uiLanguage,omitempty"`
	AllowEarlyPlay       bool    `json:"allowEarlyPlay"`
	TokenAuthEnabled     bool    `json:"tokenAuthEnabled"`
	EnableMP4Fallback    bool    `json:"enableMP4Fallback"`
	ShowHeatmap          bool    `json:"showHeatmap"`
	FontFamily           string  `json:"fontFamily,omitempty"`
	PlaybackSpeeds       string  `json:"playbackSpeeds,omitempty"`
	RememberPlayerPosition bool  `json:"rememberPlayerPosition"`
}

// EncodedResolution represents storage for a single codec/resolution.
type EncodedResolution struct {
	Codec      string `json:"codec,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Size       int64  `json:"size"`
}

// StorageSizeData represents video storage breakdown.
type StorageSizeData struct {
	Encoded       []EncodedResolution `json:"encoded,omitempty"`
	Thumbnails    int64               `json:"thumbnails"`
	Previews      int64               `json:"previews"`
	Originals     int64               `json:"originals"`
	Mp4Fallback   int64               `json:"mp4Fallback"`
	Miscellaneous int64               `json:"miscellaneous"`
	CalculatedAt  string              `json:"calculatedAt,omitempty"`
}

// StorageSizeResponse represents the storage size API response.
type StorageSizeResponse struct {
	Success    bool             `json:"success"`
	Message    string           `json:"message,omitempty"`
	StatusCode int              `json:"statusCode"`
	Data       *StorageSizeData `json:"data,omitempty"`
}

// RepackageOptions specifies options for repackage operation.
type RepackageOptions struct {
	KeepOriginalFiles bool `url:"keepOriginalFiles,omitempty"`
}

// TranscribeRequest represents a request to transcribe a video.
type TranscribeRequest struct {
	TargetLanguages     []string `json:"targetLanguages,omitempty"`
	GenerateTitle       *bool    `json:"generateTitle,omitempty"`
	GenerateDescription *bool    `json:"generateDescription,omitempty"`
	GenerateChapters    *bool    `json:"generateChapters,omitempty"`
	GenerateMoments     *bool    `json:"generateMoments,omitempty"`
	SourceLanguage      string   `json:"sourceLanguage,omitempty"`
}

// TranscribeOptions specifies options for transcribe operation.
type TranscribeOptions struct {
	Force bool `url:"force,omitempty"`
}

// SmartActionsRequest represents a request to trigger smart actions.
type SmartActionsRequest struct {
	GenerateTitle       *bool  `json:"generateTitle,omitempty"`
	GenerateDescription *bool  `json:"generateDescription,omitempty"`
	GenerateChapters    *bool  `json:"generateChapters,omitempty"`
	GenerateMoments     *bool  `json:"generateMoments,omitempty"`
	SourceLanguage      string `json:"sourceLanguage,omitempty"`
}

// ResolutionReference represents a resolution with codec info.
type ResolutionReference struct {
	Resolution string `json:"resolution,omitempty"`
	Codec      string `json:"codec,omitempty"`
}

// StorageObject represents a storage object entry.
type StorageObject struct {
	Path       string `json:"path,omitempty"`
	Size       int64  `json:"size"`
	Resolution string `json:"resolution,omitempty"`
}

// ResolutionsInfoData represents video resolution details.
type ResolutionsInfoData struct {
	VideoID                          string                `json:"videoId,omitempty"`
	VideoLibraryID                   int64                 `json:"videoLibraryId"`
	AvailableResolutions             []string              `json:"availableResolutions,omitempty"`
	ConfiguredResolutions            []string              `json:"configuredResolutions,omitempty"`
	PlaylistResolutions              []ResolutionReference `json:"playlistResolutions,omitempty"`
	StorageResolutions               []ResolutionReference `json:"storageResolutions,omitempty"`
	Mp4Resolutions                   []ResolutionReference `json:"mp4Resolutions,omitempty"`
	StorageObjects                   []StorageObject       `json:"storageObjects,omitempty"`
	OldResolutions                   []StorageObject       `json:"oldResolutions,omitempty"`
	HasBothOldAndNewResolutionFormat bool                  `json:"hasBothOldAndNewResolutionFormat"`
	HasOriginal                      bool                  `json:"hasOriginal"`
}

// ResolutionsInfoResponse represents the resolutions info API response.
type ResolutionsInfoResponse struct {
	Success    bool                 `json:"success"`
	Message    string               `json:"message,omitempty"`
	StatusCode int                  `json:"statusCode"`
	Data       *ResolutionsInfoData `json:"data,omitempty"`
}

// OEmbedOptions specifies options for oEmbed request.
type OEmbedOptions struct {
	URL       string `url:"url,omitempty"`
	MaxWidth  int    `url:"maxWidth,omitempty"`
	MaxHeight int    `url:"maxHeight,omitempty"`
	Token     string `url:"token,omitempty"`
	Expires   int64  `url:"expires,omitempty"`
}

// OEmbedResponse represents an oEmbed response.
type OEmbedResponse struct {
	Version      string `json:"version,omitempty"`
	Title        string `json:"title,omitempty"`
	Type         string `json:"type,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	HTML         string `json:"html,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	ProviderURL  string `json:"provider_url,omitempty"`
}
