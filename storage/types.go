// Package storage provides types and services for the Bunny.net Storage API.
package storage

import "github.com/geraldo/bunny-sdk-go/internal"

// Zone represents a storage zone in Bunny.net.
type Zone struct {
	ID                  int64     `json:"Id"`
	Name                string    `json:"Name"`
	Password            string    `json:"Password,omitempty"`
	ReadOnlyPassword    string    `json:"ReadOnlyPassword,omitempty"`
	Region              string    `json:"Region"`
	ReplicationRegions  []string  `json:"ReplicationRegions,omitempty"`
	StorageUsed         int64     `json:"StorageUsed"`
	FilesStored         int64     `json:"FilesStored"`
	DateModified        internal.BunnyTime `json:"DateModified"`
	Deleted             bool      `json:"Deleted"`
	PullZones           []int64   `json:"PullZones,omitempty"`
	OriginURL           string    `json:"OriginUrl,omitempty"`
	Custom404FilePath   string    `json:"Custom404FilePath,omitempty"`
	Rewrite404To200     bool      `json:"Rewrite404To200"`
}

// File represents a file or directory in a storage zone.
type File struct {
	GUID            string    `json:"Guid"`
	StorageZoneName string    `json:"StorageZoneName"`
	Path            string    `json:"Path"`
	ObjectName      string    `json:"ObjectName"`
	Length          int64     `json:"Length"`
	LastChanged     internal.BunnyTime `json:"LastChanged"`
	IsDirectory     bool      `json:"IsDirectory"`
	DateCreated     internal.BunnyTime `json:"DateCreated"`
	ServerID        int       `json:"ServerId"`
	StorageZoneID   int64     `json:"StorageZoneId"`
	UserID          string    `json:"UserId,omitempty"`
}

// ZoneListResponse represents the paginated response from listing zones.
type ZoneListResponse struct {
	Items       []Zone `json:"Items"`
	TotalItems  int    `json:"TotalItems"`
	CurrentPage int    `json:"CurrentPage"`
	PageSize    int    `json:"PageSize"`
}

// HasMore returns true if there are more pages to fetch.
func (r *ZoneListResponse) HasMore() bool {
	return (r.CurrentPage+1)*r.PageSize < r.TotalItems
}

// CreateZoneRequest represents a request to create a storage zone.
type CreateZoneRequest struct {
	Name               string   `json:"Name"`
	Region             string   `json:"Region,omitempty"`
	ReplicationRegions []string `json:"ReplicationRegions,omitempty"`
	OriginURL          string   `json:"OriginUrl,omitempty"`
}

// UpdateZoneRequest represents a request to update a storage zone.
type UpdateZoneRequest struct {
	ReplicationRegions []string `json:"ReplicationRegions,omitempty"`
	OriginURL          string   `json:"OriginUrl,omitempty"`
	Custom404FilePath  string   `json:"Custom404FilePath,omitempty"`
	Rewrite404To200    *bool    `json:"Rewrite404To200,omitempty"`
}

// ZoneListOptions specifies options for listing storage zones.
type ZoneListOptions struct {
	Page           int    `url:"page,omitempty"`
	PerPage        int    `url:"perPage,omitempty"`
	IncludeDeleted bool   `url:"includeDeleted,omitempty"`
	Search         string `url:"search,omitempty"`
}

// UploadOptions specifies options for file uploads.
type UploadOptions struct {
	Checksum    string // SHA256 uppercase hex
	ContentType string // MIME type
}

// AvailabilityResponse represents the response from checking zone name availability.
type AvailabilityResponse struct {
	Available bool   `json:"Available"`
	Name      string `json:"Name"`
}

// ResetPasswordResponse represents the response from resetting a zone password.
type ResetPasswordResponse struct {
	ID       int64  `json:"Id"`
	Password string `json:"Password,omitempty"`
	Success  bool   `json:"Success"`
}

// ResetReadOnlyPasswordResponse represents the response from resetting a zone's read-only password.
type ResetReadOnlyPasswordResponse struct {
	ID               int64  `json:"Id"`
	ReadOnlyPassword string `json:"ReadOnlyPassword,omitempty"`
	Success          bool   `json:"Success"`
}

// Region represents storage region codes.
type Region string

const (
	RegionFalkenstein  Region = "de"  // Falkenstein, Germany (default)
	RegionNewYork      Region = "ny"  // New York
	RegionLosAngeles   Region = "la"  // Los Angeles
	RegionSingapore    Region = "sg"  // Singapore
	RegionSydney       Region = "syd" // Sydney
	RegionStockholm    Region = "se"  // Stockholm
	RegionSaoPaulo     Region = "br"  // Sao Paulo
	RegionJohannesburg Region = "jh"  // Johannesburg
	RegionLondon       Region = "uk"  // London
)

// RegionBaseURL returns the base URL for a given region.
func RegionBaseURL(region Region) string {
	switch region {
	case RegionFalkenstein:
		return "https://storage.bunnycdn.com"
	case RegionNewYork:
		return "https://ny.storage.bunnycdn.com"
	case RegionLosAngeles:
		return "https://la.storage.bunnycdn.com"
	case RegionSingapore:
		return "https://sg.storage.bunnycdn.com"
	case RegionSydney:
		return "https://syd.storage.bunnycdn.com"
	case RegionStockholm:
		return "https://se.storage.bunnycdn.com"
	case RegionSaoPaulo:
		return "https://br.storage.bunnycdn.com"
	case RegionJohannesburg:
		return "https://jh.storage.bunnycdn.com"
	case RegionLondon:
		return "https://uk.storage.bunnycdn.com"
	default:
		return "https://storage.bunnycdn.com"
	}
}
