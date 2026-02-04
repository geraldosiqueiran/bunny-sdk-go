// Package shield provides types and services for the Bunny.net Shield/WAF API.
package shield

// ShieldZone represents a Shield zone in Bunny.net.
type ShieldZone struct {
	ID           string   `json:"Id"`
	Name         string   `json:"Name"`
	HostNames    []string `json:"HostNames,omitempty"`
	EdgeScriptID *string  `json:"EdgeScriptId,omitempty"`
	DateCreated  string   `json:"DateCreated,omitempty"`
}

// ZoneListResponse represents the response from listing Shield zones.
type ZoneListResponse struct {
	Items      []ShieldZone `json:"Items"`
	TotalCount int          `json:"TotalCount"`
}

// CreateZoneRequest represents a request to create a Shield zone.
type CreateZoneRequest struct {
	Name      string   `json:"Name"`
	HostNames []string `json:"HostNames,omitempty"`
}

// UpdateZoneRequest represents a request to update a Shield zone.
type UpdateZoneRequest struct {
	Name      string   `json:"Name,omitempty"`
	HostNames []string `json:"HostNames,omitempty"`
}

// PullZoneMapping represents a Shield zone to Pull zone mapping.
type PullZoneMapping struct {
	ShieldZoneID string `json:"ShieldZoneId"`
	PullZoneID   int64  `json:"PullZoneId"`
}

// PullZoneMappingResponse represents the response from getting zone mappings.
type PullZoneMappingResponse struct {
	Items []PullZoneMapping `json:"Items"`
}

// WAFRule represents a WAF rule (predefined or custom).
type WAFRule struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	RuleType    string `json:"RuleType,omitempty"`
	Category    string `json:"Category,omitempty"`
	IsActive    bool   `json:"IsActive"`
}

// WAFRuleListResponse represents the response from listing WAF rules.
type WAFRuleListResponse struct {
	Items      []WAFRule `json:"Items"`
	TotalCount int       `json:"TotalCount"`
}

// CustomRule represents a custom WAF rule.
type CustomRule struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	Description  string `json:"Description,omitempty"`
	Pattern      string `json:"Pattern,omitempty"`
	Action       string `json:"Action,omitempty"`
	ShieldZoneID string `json:"ShieldZoneId,omitempty"`
	IsActive     bool   `json:"IsActive"`
	DateCreated  string `json:"DateCreated,omitempty"`
}

// CustomRuleListResponse represents the response from listing custom WAF rules.
type CustomRuleListResponse struct {
	Items      []CustomRule `json:"Items"`
	TotalCount int          `json:"TotalCount"`
}

// CreateCustomRuleRequest represents a request to create a custom WAF rule.
type CreateCustomRuleRequest struct {
	Name         string `json:"Name"`
	Description  string `json:"Description,omitempty"`
	Pattern      string `json:"Pattern,omitempty"`
	Action       string `json:"Action,omitempty"`
	ShieldZoneID string `json:"ShieldZoneId,omitempty"`
	IsActive     bool   `json:"IsActive"`
}

// UpdateCustomRuleRequest represents a request to update a custom WAF rule (PATCH).
type UpdateCustomRuleRequest struct {
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
	Pattern     string `json:"Pattern,omitempty"`
	Action      string `json:"Action,omitempty"`
	IsActive    *bool  `json:"IsActive,omitempty"`
}

// ReplaceCustomRuleRequest represents a request to replace a custom WAF rule (PUT).
type ReplaceCustomRuleRequest struct {
	Name         string `json:"Name"`
	Description  string `json:"Description,omitempty"`
	Pattern      string `json:"Pattern,omitempty"`
	Action       string `json:"Action,omitempty"`
	ShieldZoneID string `json:"ShieldZoneId,omitempty"`
	IsActive     bool   `json:"IsActive"`
}

// WAFProfile represents a WAF protection profile.
type WAFProfile struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	IsDefault   bool   `json:"IsDefault"`
}

// WAFProfilesResponse represents the response from getting WAF profiles.
type WAFProfilesResponse struct {
	Items []WAFProfile `json:"Items"`
}

// WAFEngineConfigRule represents a rule in WAF engine config.
type WAFEngineConfigRule struct {
	ID       string `json:"Id"`
	IsActive bool   `json:"IsActive"`
	Action   string `json:"Action,omitempty"`
}

// WAFEngineConfig represents WAF engine configuration.
type WAFEngineConfig struct {
	ProfileID    string                `json:"ProfileId,omitempty"`
	IsEnabled    bool                  `json:"IsEnabled"`
	AnalysisMode string                `json:"AnalysisMode,omitempty"`
	Rules        []WAFEngineConfigRule `json:"Rules,omitempty"`
}

// WAFEnums represents available WAF enumeration values.
type WAFEnums struct {
	RuleActions []string `json:"RuleActions,omitempty"`
	RuleTypes   []string `json:"RuleTypes,omitempty"`
	Categories  []string `json:"Categories,omitempty"`
}

// TriggeredRule represents a WAF rule that was triggered and needs review.
type TriggeredRule struct {
	RuleID            string `json:"RuleId"`
	RuleName          string `json:"RuleName"`
	TriggerCount      int64  `json:"TriggerCount"`
	LastTriggered     string `json:"LastTriggered,omitempty"`
	RecommendedAction string `json:"RecommendedAction,omitempty"`
}

// TriggeredRulesResponse represents the response from getting triggered rules.
type TriggeredRulesResponse struct {
	Items      []TriggeredRule `json:"Items"`
	TotalCount int             `json:"TotalCount"`
}

// TriggeredRuleReviewRequest represents a request to submit a triggered rule review.
type TriggeredRuleReviewRequest struct {
	RuleID  string `json:"RuleId"`
	Action  string `json:"Action"`
	Comment string `json:"Comment,omitempty"`
}

// TriggeredRuleReview represents a submitted triggered rule review.
type TriggeredRuleReview struct {
	ReviewID      string `json:"ReviewId"`
	RuleID        string `json:"RuleId"`
	Action        string `json:"Action"`
	Comment       string `json:"Comment,omitempty"`
	DateSubmitted string `json:"DateSubmitted,omitempty"`
}

// AIRecommendation represents an AI recommendation for a triggered rule.
type AIRecommendation struct {
	RuleID         string  `json:"RuleId"`
	RuleName       string  `json:"RuleName"`
	Recommendation string  `json:"Recommendation"`
	Confidence     float64 `json:"Confidence"`
}

// AIRecommendationResponse represents the response from getting AI recommendations.
type AIRecommendationResponse struct {
	Recommendations []AIRecommendation `json:"Recommendations"`
}

// PlanSegment represents a plan segment with rule limits.
type PlanSegment struct {
	PlanID         string `json:"PlanId"`
	AvailableRules int    `json:"AvailableRules"`
	MaxCustomRules int    `json:"MaxCustomRules"`
}

// PlanSegmentationResponse represents the response from getting plan segmentation.
type PlanSegmentationResponse struct {
	Plans []PlanSegment `json:"Plans"`
}

// AccessList contains allowed, blocked, and challenged entries.
type AccessList struct {
	Allowed    []AccessListEntry `json:"Allowed,omitempty"`
	Blocked    []AccessListEntry `json:"Blocked,omitempty"`
	Challenged []AccessListEntry `json:"Challenged,omitempty"`
}

// AccessListEntry represents an entry in an access list.
type AccessListEntry struct {
	Type      string `json:"Type"`
	Value     string `json:"Value"`
	Action    string `json:"Action,omitempty"`
	Comment   string `json:"Comment,omitempty"`
	DateAdded string `json:"DateAdded,omitempty"`
}

// AddAccessListEntryRequest represents a request to add an access list entry.
type AddAccessListEntryRequest struct {
	Type    string `json:"Type"`
	Value   string `json:"Value"`
	Action  string `json:"Action"`
	Comment string `json:"Comment,omitempty"`
}

// AccessListEntryUpdate represents an update to an access list entry.
type AccessListEntryUpdate struct {
	Type    string `json:"Type"`
	Value   string `json:"Value"`
	Action  string `json:"Action,omitempty"`
	Comment string `json:"Comment,omitempty"`
}

// UpdateAccessListEntriesRequest represents a request to update access list entries.
type UpdateAccessListEntriesRequest struct {
	Updates []AccessListEntryUpdate `json:"Updates"`
}

// AccessListEntryIdentifier identifies an access list entry for deletion.
type AccessListEntryIdentifier struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

// DeleteAccessListEntriesRequest represents a request to delete access list entries.
type DeleteAccessListEntriesRequest struct {
	Entries []AccessListEntryIdentifier `json:"Entries"`
}

// AccessListEnums represents available access list types and actions.
type AccessListEnums struct {
	Types   []string `json:"Types,omitempty"`
	Actions []string `json:"Actions,omitempty"`
}

// AccessListConfig represents access list configuration.
type AccessListConfig struct {
	DefaultAction string `json:"DefaultAction,omitempty"`
	IsEnabled     bool   `json:"IsEnabled"`
	LogUnmatched  bool   `json:"LogUnmatched"`
	DateUpdated   string `json:"DateUpdated,omitempty"`
}

// UpdateAccessListConfigRequest represents a request to update access list config.
type UpdateAccessListConfigRequest struct {
	DefaultAction string `json:"DefaultAction,omitempty"`
	IsEnabled     *bool  `json:"IsEnabled,omitempty"`
	LogUnmatched  *bool  `json:"LogUnmatched,omitempty"`
}

// RateLimit represents a rate limit rule.
type RateLimit struct {
	ID                string `json:"Id"`
	Name              string `json:"Name"`
	Path              string `json:"Path,omitempty"`
	RequestsPerSecond int    `json:"RequestsPerSecond,omitempty"`
	RequestsPerMinute int    `json:"RequestsPerMinute,omitempty"`
	Action            string `json:"Action,omitempty"`
	ShieldZoneID      string `json:"ShieldZoneId,omitempty"`
	IsActive          bool   `json:"IsActive"`
	DateCreated       string `json:"DateCreated,omitempty"`
}

// RateLimitListResponse represents the response from listing rate limits.
type RateLimitListResponse struct {
	Items      []RateLimit `json:"Items"`
	TotalCount int         `json:"TotalCount"`
}

// CreateRateLimitRequest represents a request to create a rate limit rule.
type CreateRateLimitRequest struct {
	Name              string `json:"Name"`
	Path              string `json:"Path,omitempty"`
	RequestsPerSecond int    `json:"RequestsPerSecond,omitempty"`
	RequestsPerMinute int    `json:"RequestsPerMinute,omitempty"`
	Action            string `json:"Action,omitempty"`
	ShieldZoneID      string `json:"ShieldZoneId,omitempty"`
	IsActive          bool   `json:"IsActive"`
}

// UpdateRateLimitRequest represents a request to update a rate limit rule.
type UpdateRateLimitRequest struct {
	Name              string `json:"Name,omitempty"`
	Path              string `json:"Path,omitempty"`
	RequestsPerSecond *int   `json:"RequestsPerSecond,omitempty"`
	RequestsPerMinute *int   `json:"RequestsPerMinute,omitempty"`
	Action            string `json:"Action,omitempty"`
	IsActive          *bool  `json:"IsActive,omitempty"`
}

// BotDetectionSettings represents bot detection settings for a zone.
type BotDetectionSettings struct {
	IsEnabled      bool     `json:"IsEnabled"`
	DetectionLevel string   `json:"DetectionLevel,omitempty"`
	Action         string   `json:"Action,omitempty"`
	AllowedBots    []string `json:"AllowedBots,omitempty"`
	BlockedBots    []string `json:"BlockedBots,omitempty"`
}

// UpdateBotDetectionRequest represents a request to update bot detection settings.
type UpdateBotDetectionRequest struct {
	IsEnabled      *bool    `json:"IsEnabled,omitempty"`
	DetectionLevel string   `json:"DetectionLevel,omitempty"`
	Action         string   `json:"Action,omitempty"`
	AllowedBots    []string `json:"AllowedBots,omitempty"`
	BlockedBots    []string `json:"BlockedBots,omitempty"`
}

// UploadScanningConfig represents upload scanning configuration.
type UploadScanningConfig struct {
	IsEnabled          bool     `json:"IsEnabled"`
	ScanLevel          string   `json:"ScanLevel,omitempty"`
	QuarantineInfected bool     `json:"QuarantineInfected"`
	NotifyOnDetection  bool     `json:"NotifyOnDetection"`
	AllowedFileTypes   []string `json:"AllowedFileTypes,omitempty"`
	MaxFileSize        int64    `json:"MaxFileSize,omitempty"`
}

// UpdateUploadScanningRequest represents a request to update upload scanning config.
type UpdateUploadScanningRequest struct {
	IsEnabled          *bool    `json:"IsEnabled,omitempty"`
	ScanLevel          string   `json:"ScanLevel,omitempty"`
	QuarantineInfected *bool    `json:"QuarantineInfected,omitempty"`
	NotifyOnDetection  *bool    `json:"NotifyOnDetection,omitempty"`
	AllowedFileTypes   []string `json:"AllowedFileTypes,omitempty"`
	MaxFileSize        *int64   `json:"MaxFileSize,omitempty"`
}

// DateRange represents a date range in metrics responses.
type DateRange struct {
	From string `json:"From"`
	To   string `json:"To"`
}

// DateRangeOptions specifies date range filtering for metrics.
type DateRangeOptions struct {
	From string
	To   string
}

// MetricsDetailedOptions specifies options for detailed metrics.
type MetricsDetailedOptions struct {
	From   string
	To     string
	ZoneID string
}

// MetricsOverview represents overview metrics for Shield.
type MetricsOverview struct {
	TotalRequests      int64     `json:"TotalRequests"`
	BlockedRequests    int64     `json:"BlockedRequests"`
	AllowedRequests    int64     `json:"AllowedRequests"`
	BotDetectionBlocks int64     `json:"BotDetectionBlocks"`
	RateLimitBlocks    int64     `json:"RateLimitBlocks"`
	AccessListBlocks   int64     `json:"AccessListBlocks"`
	DateRange          DateRange `json:"DateRange"`
}

// MetricsBreakdown represents a breakdown of blocked requests.
type MetricsBreakdown struct {
	BotDetection int64 `json:"BotDetection"`
	RateLimit    int64 `json:"RateLimit"`
	AccessList   int64 `json:"AccessList"`
}

// ZoneMetrics represents metrics for a specific zone.
type ZoneMetrics struct {
	ZoneID          string           `json:"ZoneId"`
	ZoneName        string           `json:"ZoneName"`
	TotalRequests   int64            `json:"TotalRequests"`
	BlockedRequests int64            `json:"BlockedRequests"`
	AllowedRequests int64            `json:"AllowedRequests"`
	Breakdown       MetricsBreakdown `json:"Breakdown"`
}

// MetricsOverviewDetailed represents detailed overview metrics.
type MetricsOverviewDetailed struct {
	Zones     []ZoneMetrics `json:"Zones"`
	DateRange DateRange     `json:"DateRange"`
}

// TopURLEntry represents a top URL in WAF rule metrics.
type TopURLEntry struct {
	URL          string `json:"Url"`
	TriggerCount int64  `json:"TriggerCount"`
}

// WAFRuleMetrics represents metrics for a specific WAF rule.
type WAFRuleMetrics struct {
	RuleID       string        `json:"RuleId"`
	RuleName     string        `json:"RuleName"`
	TriggerCount int64         `json:"TriggerCount"`
	BlockedCount int64         `json:"BlockedCount"`
	AllowedCount int64         `json:"AllowedCount"`
	TopURLs      []TopURLEntry `json:"TopUrls,omitempty"`
}

// TopIPEntry represents a top IP in rate limit metrics.
type TopIPEntry struct {
	IP         string `json:"IP"`
	BlockCount int64  `json:"BlockCount"`
}

// RateLimitMetrics represents metrics for a specific rate limit rule.
type RateLimitMetrics struct {
	RuleID          string       `json:"RuleId"`
	RuleName        string       `json:"RuleName"`
	Path            string       `json:"Path,omitempty"`
	BlockedRequests int64        `json:"BlockedRequests"`
	TopBlockedIPs   []TopIPEntry `json:"TopBlockedIPs,omitempty"`
}

// RateLimitMetricsSummary represents summary metrics for a rate limit.
type RateLimitMetricsSummary struct {
	RuleID          string `json:"RuleId"`
	RuleName        string `json:"RuleName"`
	BlockedRequests int64  `json:"BlockedRequests"`
}

// RateLimitMetricsList represents the response from getting all rate limit metrics.
type RateLimitMetricsList struct {
	Items      []RateLimitMetricsSummary `json:"Items"`
	TotalCount int                       `json:"TotalCount"`
}

// BotTypeEntry represents a bot type in metrics.
type BotTypeEntry struct {
	BotType string `json:"BotType"`
	Count   int64  `json:"Count"`
}

// BotDetectionMetrics represents bot detection metrics for a zone.
type BotDetectionMetrics struct {
	ZoneID              string         `json:"ZoneId"`
	ZoneName            string         `json:"ZoneName"`
	TotalBotRequests    int64          `json:"TotalBotRequests"`
	BlockedBots         int64          `json:"BlockedBots"`
	ChallengedBots      int64          `json:"ChallengedBots"`
	DetectionLevel      string         `json:"DetectionLevel,omitempty"`
	TopDetectedBotTypes []BotTypeEntry `json:"TopDetectedBotTypes,omitempty"`
}

// UploadScanningMetrics represents upload scanning metrics for a zone.
type UploadScanningMetrics struct {
	ZoneID            string `json:"ZoneId"`
	ZoneName          string `json:"ZoneName"`
	TotalScannedFiles int64  `json:"TotalScannedFiles"`
	CleanFiles        int64  `json:"CleanFiles"`
	MaliciousFiles    int64  `json:"MaliciousFiles"`
	QuarantinedFiles  int64  `json:"QuarantinedFiles"`
	ScanErrors        int64  `json:"ScanErrors"`
}

// EventLogListOptions specifies options for listing event logs.
type EventLogListOptions struct {
	ZoneID string
	From   string
	To     string
	Limit  int
	Offset int
}

// EventLog represents a security event log entry.
type EventLog struct {
	ID         string `json:"Id"`
	ZoneID     string `json:"ZoneId"`
	Timestamp  string `json:"Timestamp"`
	EventType  string `json:"EventType"`
	RuleID     string `json:"RuleId,omitempty"`
	RuleName   string `json:"RuleName,omitempty"`
	SourceIP   string `json:"SourceIP"`
	Path       string `json:"Path"`
	Method     string `json:"Method"`
	Action     string `json:"Action"`
	StatusCode int    `json:"StatusCode"`
	UserAgent  string `json:"UserAgent,omitempty"`
}

// EventLogListResponse represents the response from listing event logs.
type EventLogListResponse struct {
	Items      []EventLog `json:"Items"`
	TotalCount int        `json:"TotalCount"`
}

// DDoSProfile represents a DDoS protection profile.
type DDoSProfile struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

// DDoSEnums represents DDoS protection enumeration values.
type DDoSEnums struct {
	Profiles []DDoSProfile `json:"Profiles,omitempty"`
	Triggers []string      `json:"Triggers,omitempty"`
}

// Promo represents a promotional offer.
type Promo struct {
	ID          string `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description,omitempty"`
	ValidUntil  string `json:"ValidUntil,omitempty"`
}

// PromoInfo represents promotional information.
type PromoInfo struct {
	CurrentPromos []Promo `json:"CurrentPromos,omitempty"`
}
