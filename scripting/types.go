// Package scripting provides types and services for the Bunny.net Edge Scripting API.
package scripting

// ScriptType represents the type of edge script.
type ScriptType string

const (
	ScriptTypeDNS        ScriptType = "DNS"
	ScriptTypeCDN        ScriptType = "CDN"
	ScriptTypeMiddleware ScriptType = "Middleware"
)

// ReleaseStatus represents the status of a release.
type ReleaseStatus string

const (
	ReleaseStatusArchived ReleaseStatus = "Archived"
	ReleaseStatusLive     ReleaseStatus = "Live"
)

// EdgeScript represents an edge script in Bunny.net.
type EdgeScript struct {
	ID                  int64                  `json:"Id"`
	Name                *string                `json:"Name,omitempty"`
	LastModified        string                 `json:"LastModified"`
	ScriptType          ScriptType             `json:"ScriptType"`
	CurrentReleaseID    int64                  `json:"CurrentReleaseId"`
	EdgeScriptVariables []EdgeScriptVariable   `json:"EdgeScriptVariables,omitempty"`
	Deleted             bool                   `json:"Deleted"`
	LinkedPullZones     []LinkedPullZone       `json:"LinkedPullZones,omitempty"`
	Integration         *SourceCodeIntegration `json:"Integration,omitempty"`
	DefaultHostname     *string                `json:"DefaultHostname,omitempty"`
	SystemHostname      *string                `json:"SystemHostname,omitempty"`
	DeploymentKey       *string                `json:"DeploymentKey,omitempty"`
	RepositoryID        *int64                 `json:"RepositoryId,omitempty"`
	IntegrationID       *int64                 `json:"IntegrationId,omitempty"`
	MonthlyCost         float64                `json:"MonthlyCost"`
	MonthlyRequestCount int64                  `json:"MonthlyRequestCount"`
	MonthlyCpuTime      int64                  `json:"MonthlyCpuTime"`
}

// LinkedPullZone represents a pull zone linked to an edge script.
type LinkedPullZone struct {
	ID   int64  `json:"Id"`
	Name string `json:"Name,omitempty"`
}

// SourceCodeIntegration represents source code integration configuration.
type SourceCodeIntegration struct {
	ID             int64   `json:"Id"`
	IntegrationType string `json:"IntegrationType,omitempty"`
	RepositoryURL  *string `json:"RepositoryUrl,omitempty"`
	Branch         *string `json:"Branch,omitempty"`
}

// EdgeScriptRelease represents a release of an edge script.
type EdgeScriptRelease struct {
	ID            int64          `json:"Id"`
	Deleted       bool           `json:"Deleted"`
	Code          *string        `json:"Code,omitempty"`
	UUID          *string        `json:"Uuid,omitempty"`
	Note          *string        `json:"Note,omitempty"`
	Author        *string        `json:"Author,omitempty"`
	AuthorEmail   *string        `json:"AuthorEmail,omitempty"`
	CommitSha     *string        `json:"CommitSha,omitempty"`
	Status        *ReleaseStatus `json:"Status,omitempty"`
	DateReleased  string         `json:"DateReleased"`
	DatePublished string         `json:"DatePublished"`
}

// EdgeScriptSecret represents a secret for an edge script.
type EdgeScriptSecret struct {
	ID           int64   `json:"Id"`
	Name         *string `json:"Name,omitempty"`
	LastModified string  `json:"LastModified"`
}

// EdgeScriptVariable represents a variable for an edge script.
type EdgeScriptVariable struct {
	ID           int64   `json:"Id"`
	Name         *string `json:"Name,omitempty"`
	Required     bool    `json:"Required"`
	DefaultValue *string `json:"DefaultValue,omitempty"`
}

// EdgeScriptCode represents the code of an edge script.
type EdgeScriptCode struct {
	Code         *string `json:"Code,omitempty"`
	LastModified string  `json:"LastModified"`
}

// ScriptStatistics represents statistics for an edge script.
type ScriptStatistics struct {
	TotalRequestsServed        int64              `json:"TotalRequestsServed"`
	TotalCpuUsed               float64            `json:"TotalCpuUsed"`
	TotalMonthlyCost           float64            `json:"TotalMonthlyCost"`
	AverageCpuTimePerExecution float64            `json:"AverageCpuTimePerExecution"`
	RequestsServedChart        map[string]int64   `json:"RequestsServedChart,omitempty"`
	AverageCpuTimeChart        map[string]float64 `json:"AverageCpuTimeChart,omitempty"`
	TotalCpuTimeChart          map[string]float64 `json:"TotalCpuTimeChart,omitempty"`
}

// ScriptListResponse represents the response from listing edge scripts.
type ScriptListResponse struct {
	Items        []EdgeScript `json:"Items"`
	CurrentPage  int          `json:"CurrentPage"`
	TotalItems   int          `json:"TotalItems"`
	HasMoreItems bool         `json:"HasMoreItems"`
}

// ReleaseListResponse represents the response from listing releases.
type ReleaseListResponse struct {
	Items        []EdgeScriptRelease `json:"Items"`
	CurrentPage  int                 `json:"CurrentPage"`
	TotalItems   int                 `json:"TotalItems"`
	HasMoreItems bool                `json:"HasMoreItems"`
}

// SecretListResponse represents the response from listing secrets.
type SecretListResponse struct {
	Secrets []EdgeScriptSecret `json:"Secrets"`
}

// ScriptListOptions specifies options for listing scripts.
type ScriptListOptions struct {
	Page                   int
	PerPage                int
	Search                 string
	IncludeLinkedPullzones bool
	IntegrationID          *int64
}

// StatisticsOptions specifies options for getting statistics.
type StatisticsOptions struct {
	DateFrom   string
	DateTo     string
	LoadLatest bool
	Hourly     bool
}

// ReleaseListOptions specifies options for listing releases.
type ReleaseListOptions struct {
	Page    int
	PerPage int
}

// CreateScriptRequest represents a request to create an edge script.
type CreateScriptRequest struct {
	Name                 string                 `json:"Name,omitempty"`
	Code                 string                 `json:"Code,omitempty"`
	ScriptType           ScriptType             `json:"ScriptType,omitempty"`
	CreateLinkedPullZone bool                   `json:"CreateLinkedPullZone"`
	LinkedPullZoneName   string                 `json:"LinkedPullZoneName,omitempty"`
	Integration          *SourceCodeIntegration `json:"Integration,omitempty"`
}

// UpdateScriptRequest represents a request to update an edge script.
type UpdateScriptRequest struct {
	Name       string     `json:"Name,omitempty"`
	ScriptType ScriptType `json:"ScriptType,omitempty"`
}

// UpdateCodeRequest represents a request to set script code.
type UpdateCodeRequest struct {
	Code string `json:"Code,omitempty"`
}

// PublishReleaseRequest represents a request to publish a release.
type PublishReleaseRequest struct {
	Note string `json:"Note,omitempty"`
}

// AddSecretRequest represents a request to add a secret.
type AddSecretRequest struct {
	Name   string `json:"Name"`
	Secret string `json:"Secret,omitempty"`
}

// UpdateSecretRequest represents a request to update a secret.
type UpdateSecretRequest struct {
	Secret string `json:"Secret,omitempty"`
}

// UpsertSecretRequest represents a request to upsert a secret.
type UpsertSecretRequest struct {
	Name   string `json:"Name"`
	Secret string `json:"Secret,omitempty"`
}

// AddVariableRequest represents a request to add a variable.
type AddVariableRequest struct {
	Name         string `json:"Name"`
	Required     bool   `json:"Required"`
	DefaultValue string `json:"DefaultValue,omitempty"`
}

// UpdateVariableRequest represents a request to update a variable.
type UpdateVariableRequest struct {
	DefaultValue string `json:"DefaultValue,omitempty"`
	Required     *bool  `json:"Required,omitempty"`
}

// UpsertVariableRequest represents a request to upsert a variable.
type UpsertVariableRequest struct {
	Name         string `json:"Name"`
	Required     *bool  `json:"Required,omitempty"`
	DefaultValue string `json:"DefaultValue,omitempty"`
}
