package containers

// =============================================================================
// Enums
// =============================================================================

// ApplicationStatus represents the status of an application.
type ApplicationStatus string

const (
	ApplicationStatusUnknown     ApplicationStatus = "Unknown"
	ApplicationStatusActive      ApplicationStatus = "Active"
	ApplicationStatusProgressing ApplicationStatus = "Progressing"
	ApplicationStatusInactive    ApplicationStatus = "Inactive"
	ApplicationStatusFailing     ApplicationStatus = "Failing"
	ApplicationStatusSuspended   ApplicationStatus = "Suspended"
)

// RuntimeType represents the runtime type of an application.
type RuntimeType string

const (
	RuntimeTypeShared   RuntimeType = "Shared"
	RuntimeTypeReserved RuntimeType = "Reserved"
)

// EndpointType represents the type of endpoint.
type EndpointType string

const (
	EndpointTypeCDN      EndpointType = "CDN"
	EndpointTypeAnycast  EndpointType = "Anycast"
	EndpointTypePublicIP EndpointType = "PublicIp"
)

// Protocol represents a network protocol.
type Protocol string

const (
	ProtocolTCP  Protocol = "Tcp"
	ProtocolUDP  Protocol = "Udp"
	ProtocolSCTP Protocol = "Sctp"
)

// ImagePullPolicy represents when to pull container images.
type ImagePullPolicy string

const (
	ImagePullPolicyAlways       ImagePullPolicy = "Always"
	ImagePullPolicyIfNotPresent ImagePullPolicy = "IfNotPresent"
)

// RegistryType represents the type of container registry.
type RegistryType string

const (
	RegistryTypeDockerHub RegistryType = "DockerHub"
	RegistryTypeGitHub    RegistryType = "GitHub"
)

// RegistryStatus represents the status of a registry operation.
type RegistryStatus string

const (
	RegistryStatusSaved                   RegistryStatus = "Saved"
	RegistryStatusSecretsValidationFailed RegistryStatus = "SecretsValidationFailed"
	RegistryStatusUnknownErrorOccured     RegistryStatus = "UnknownErrorOccured"
	RegistryStatusNotFound                RegistryStatus = "NotFound"
	RegistryStatusInvalidInput            RegistryStatus = "InvalidInput"
)

// RegistryDeleteStatus represents the status of a registry delete operation.
type RegistryDeleteStatus string

const (
	RegistryDeleteStatusRemoved  RegistryDeleteStatus = "Removed"
	RegistryDeleteStatusInUse    RegistryDeleteStatus = "InUse"
	RegistryDeleteStatusNotFound RegistryDeleteStatus = "NotFound"
)

// LogForwardingType represents the type of log forwarding.
type LogForwardingType string

const (
	LogForwardingTypeSyslogUDP LogForwardingType = "SyslogUdp"
	LogForwardingTypeSyslogTCP LogForwardingType = "SyslogTcp"
)

// LogForwardingFormat represents the format of log forwarding.
type LogForwardingFormat string

const (
	LogForwardingFormatRfc3164 LogForwardingFormat = "SyslogRfc3164"
	LogForwardingFormatRfc5424 LogForwardingFormat = "SyslogRfc5424"
)

// VolumeStatus represents the status of a volume instance.
type VolumeStatus string

const (
	VolumeStatusAttached  VolumeStatus = "Attached"
	VolumeStatusDetached  VolumeStatus = "Detached"
	VolumeStatusExtending VolumeStatus = "Extending"
	VolumeStatusDeleting  VolumeStatus = "Deleting"
	VolumeStatusCreating  VolumeStatus = "Creating"
	VolumeStatusUnknown   VolumeStatus = "Unknown"
)

// ProvisioningType represents region provisioning type.
type ProvisioningType string

const (
	ProvisioningTypeStatic  ProvisioningType = "Static"
	ProvisioningTypeDynamic ProvisioningType = "Dynamic"
)

// StatisticsGranularity represents the granularity of statistics.
type StatisticsGranularity string

const (
	StatisticsGranularityDaily  StatisticsGranularity = "Daily"
	StatisticsGranularityHourly StatisticsGranularity = "Hourly"
	StatisticsGranularityMinute StatisticsGranularity = "Minute"
)

// AnycastType represents the type of anycast endpoint.
type AnycastType string

const (
	AnycastTypeIPv4 AnycastType = "IPv4"
)

// =============================================================================
// Common/Pagination Types
// =============================================================================

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	TotalItems int `json:"totalItems"`
}

// ListOptions contains common list query parameters.
type ListOptions struct {
	NextCursor string `url:"nextCursor,omitempty"`
	Limit      int    `url:"limit,omitempty"`
}

// =============================================================================
// Application Types
// =============================================================================

// DisplayEndpoint represents an endpoint displayed in app list.
type DisplayEndpoint struct {
	ID      string       `json:"id,omitempty"`
	Address string       `json:"address,omitempty"`
	Type    EndpointType `json:"type,omitempty"`
}

// ApplicationListItem represents an application in list response.
type ApplicationListItem struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	DisplayEndpoint *DisplayEndpoint  `json:"displayEndpoint,omitempty"`
	Status          ApplicationStatus `json:"status"`
}

// ApplicationListResponse represents the response from listing applications.
type ApplicationListResponse struct {
	Items  []ApplicationListItem `json:"items"`
	Meta   PaginationMeta        `json:"meta"`
	Cursor string                `json:"cursor,omitempty"`
}

// AutoScaling represents autoscaling configuration.
type AutoScaling struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// RegionSettings represents application region settings.
type RegionSettings struct {
	AllowedRegionIds  []string               `json:"allowedRegionIds,omitempty"`
	RequiredRegionIds []string               `json:"requiredRegionIds,omitempty"`
	MaxAllowedRegions int                    `json:"maxAllowedRegions,omitempty"`
	NodeSelectors     map[string]any `json:"nodeSelectors,omitempty"`
	ProvisioningType  ProvisioningType       `json:"provisioningType,omitempty"`
}

// ContainerInstance represents a running container instance.
type ContainerInstance struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty"`
	Region   string `json:"region,omitempty"`
	NodeName string `json:"nodeName,omitempty"`
}

// Application represents a full application detail.
type Application struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Status             ApplicationStatus   `json:"status"`
	RuntimeType        RuntimeType         `json:"runtimeType,omitempty"`
	DisplayEndpoint    *DisplayEndpoint    `json:"displayEndpoint,omitempty"`
	RegionSettings     *RegionSettings     `json:"regionSettings,omitempty"`
	ContainerTemplates []ContainerTemplate `json:"containerTemplates,omitempty"`
	ContainerInstances []ContainerInstance `json:"containerInstances,omitempty"`
	Volumes            []Volume            `json:"volumes,omitempty"`
	AutoScaling        *AutoScaling        `json:"autoScaling,omitempty"`
}

// VolumeRequest represents a volume in create/update requests.
type VolumeRequest struct {
	Name string `json:"name"`
	Size int    `json:"size"` // 1-100 GB
}

// CreateApplicationRequest represents a request to create an application.
type CreateApplicationRequest struct {
	Name                          string                           `json:"name"`
	RuntimeType                   RuntimeType                      `json:"runtimeType"`
	AutoScaling                   AutoScaling                      `json:"autoScaling"`
	RegionSettings                CreateRegionSettingsRequest      `json:"regionSettings"`
	TerminationGracePeriodSeconds int                              `json:"terminationGracePeriodSeconds,omitempty"`
	ContainerTemplates            []CreateContainerTemplateRequest `json:"containerTemplates,omitempty"`
	Volumes                       []VolumeRequest                  `json:"volumes,omitempty"`
}

// CreateRegionSettingsRequest represents region settings for create request.
type CreateRegionSettingsRequest struct {
	AllowedRegionIds  []string `json:"allowedRegionIds,omitempty"`
	RequiredRegionIds []string `json:"requiredRegionIds,omitempty"`
	MaxAllowedRegions int      `json:"maxAllowedRegions,omitempty"`
}

// UpdateApplicationRequest represents a request to update an application (PUT).
type UpdateApplicationRequest struct {
	Name                          string                           `json:"name"`
	RuntimeType                   RuntimeType                      `json:"runtimeType"`
	AutoScaling                   AutoScaling                      `json:"autoScaling"`
	RegionSettings                CreateRegionSettingsRequest      `json:"regionSettings"`
	TerminationGracePeriodSeconds int                              `json:"terminationGracePeriodSeconds,omitempty"`
	ContainerTemplates            []CreateContainerTemplateRequest `json:"containerTemplates,omitempty"`
	Volumes                       []VolumeRequest                  `json:"volumes,omitempty"`
}

// PatchApplicationRequest represents a request to partially update an application.
type PatchApplicationRequest struct {
	Name               string                           `json:"name,omitempty"`
	RuntimeType        RuntimeType                      `json:"runtimeType,omitempty"`
	AutoScaling        *AutoScaling                     `json:"autoScaling,omitempty"`
	RegionSettings     *CreateRegionSettingsRequest     `json:"regionSettings,omitempty"`
	ContainerTemplates []CreateContainerTemplateRequest `json:"containerTemplates,omitempty"`
	Volumes            []VolumeRequest                  `json:"volumes,omitempty"`
}

// ApplicationIDResponse represents a response with just an application ID.
type ApplicationIDResponse struct {
	ID string `json:"id"`
}

// MetricValue represents a metric with value and status.
type MetricValue struct {
	Value  float64 `json:"value"`
	Status string  `json:"status,omitempty"`
}

// RegionOverview represents region info in overview.
type RegionOverview struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Instances int     `json:"instances,omitempty"`
	Latency   float64 `json:"latency,omitempty"`
}

// ApplicationOverview represents application overview metrics.
type ApplicationOverview struct {
	TargetLatency       *MetricValue      `json:"targetLatency,omitempty"`
	CurrentLatency      *MetricValue      `json:"currentLatency,omitempty"`
	ActiveRegions       *MetricValue      `json:"activeRegions,omitempty"`
	ActiveInstances     *MetricValue      `json:"activeInstances,omitempty"`
	DesiredInstances    int               `json:"desiredInstances,omitempty"`
	Status              ApplicationStatus `json:"status,omitempty"`
	AverageCPU          *MetricValue      `json:"averageCPU,omitempty"`
	AverageRAM          *MetricValue      `json:"averageRAM,omitempty"`
	AverageVolumesUsage *MetricValue      `json:"averageVolumesUsage,omitempty"`
	Regions             []RegionOverview  `json:"regions,omitempty"`
	AverageLatency      float64           `json:"averageLatency,omitempty"`
	TotalVolumeSizeInGb int               `json:"totalVolumeSizeInGb,omitempty"`
	MonthlyCost         float64           `json:"monthlyCost,omitempty"`
	LatencyChart        map[string]any    `json:"latencyChart,omitempty"`
}

// StatisticsOptions represents query parameters for statistics.
type StatisticsOptions struct {
	FromDate    string                `url:"fromDate"`
	ToDate      string                `url:"toDate,omitempty"`
	Granularity StatisticsGranularity `url:"granularity"`
}

// ApplicationStatistics represents application statistics response.
type ApplicationStatistics struct {
	TargetLatencyChart        map[string]any `json:"targetLatencyChart,omitempty"`
	ActiveRegionsChart        map[string]any `json:"activeRegionsChart,omitempty"`
	LatencyChart              map[string]any `json:"latencyChart,omitempty"`
	CPUUsageChart             map[string]any `json:"cpuUsageChart,omitempty"`
	RAMUsageChart             map[string]any `json:"ramUsageChart,omitempty"`
	TrafficChart              map[string]any `json:"trafficChart,omitempty"`
	InstancesChart            map[string]any `json:"instancesChart,omitempty"`
	VolumesUsageChart         map[string]any `json:"volumesUsageChart,omitempty"`
	VolumesCapacityChart      map[string]any `json:"volumesCapacityChart,omitempty"`
	VolumesSplitUsageChart    map[string]any `json:"volumesSplitUsageChart,omitempty"`
	VolumesSplitCapacityChart map[string]any `json:"volumesSplitCapacityChart,omitempty"`
}

// =============================================================================
// Container Registry Types
// =============================================================================

// ContainerRegistry represents a container registry.
type ContainerRegistry struct {
	ID                   int64  `json:"id"`
	AccountID            string `json:"accountId,omitempty"`
	UserID               string `json:"userId,omitempty"`
	NamespaceID          string `json:"namespaceId,omitempty"`
	DisplayName          string `json:"displayName"`
	HostName             string `json:"hostName,omitempty"`
	UserName             string `json:"userName,omitempty"`
	FirstPasswordSymbols string `json:"firstPasswordSymbols,omitempty"`
	LastPasswordSymbols  string `json:"lastPasswordSymbols,omitempty"`
	CreatedAt            string `json:"createdAt,omitempty"`
	IsPublic             bool   `json:"isPublic,omitempty"`
	LastUpdatedAt        string `json:"lastUpdatedAt,omitempty"`
}

// ContainerRegistryListResponse represents the response from listing registries.
type ContainerRegistryListResponse struct {
	Items  []ContainerRegistry `json:"items"`
	Meta   PaginationMeta      `json:"meta"`
	Cursor string              `json:"cursor,omitempty"`
}

// PasswordCredentials represents registry credentials.
type PasswordCredentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// CreateRegistryRequest represents a request to create a registry.
type CreateRegistryRequest struct {
	DisplayName         string               `json:"displayName"`
	Type                RegistryType         `json:"type,omitempty"`
	PasswordCredentials *PasswordCredentials `json:"passwordCredentials,omitempty"`
}

// UpdateRegistryRequest represents a request to update a registry.
type UpdateRegistryRequest struct {
	DisplayName         string               `json:"displayName"`
	Type                RegistryType         `json:"type,omitempty"`
	PasswordCredentials *PasswordCredentials `json:"passwordCredentials,omitempty"`
}

// RegistryOperationResponse represents a registry operation response.
type RegistryOperationResponse struct {
	ID     int64          `json:"id,omitempty"`
	Error  string         `json:"error,omitempty"`
	Status RegistryStatus `json:"status,omitempty"`
}

// RegistryDeleteResponse represents a registry delete response.
type RegistryDeleteResponse struct {
	Status       RegistryDeleteStatus `json:"status"`
	Applications []string             `json:"applications,omitempty"`
}

// ListImagesRequest represents a request to list container images.
type ListImagesRequest struct {
	RegistryID string `json:"registryId"`
}

// ContainerImage represents a container image.
type ContainerImage struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace,omitempty"`
}

// ListTagsRequest represents a request to list image tags.
type ListTagsRequest struct {
	RegistryID     string `json:"registryId"`
	ImageName      string `json:"imageName"`
	ImageNamespace string `json:"imageNamespace"`
}

// ImageTag represents an image tag.
type ImageTag struct {
	Name string `json:"name"`
}

// GetDigestRequest represents a request to get image digest.
type GetDigestRequest struct {
	RegistryID     string `json:"registryId"`
	ImageName      string `json:"imageName"`
	ImageNamespace string `json:"imageNamespace"`
	Tag            string `json:"tag"`
}

// ImageDigest represents an image digest response.
type ImageDigest struct {
	ImageNamespace string `json:"imageNamespace,omitempty"`
	Image          string `json:"image,omitempty"`
	Tag            string `json:"tag,omitempty"`
	Digest         string `json:"digest,omitempty"`
}

// GetConfigSuggestionsRequest represents a request for config suggestions.
type GetConfigSuggestionsRequest struct {
	RegistryID     string `json:"registryId"`
	ImageName      string `json:"imageName"`
	ImageNamespace string `json:"imageNamespace"`
	Tag            string `json:"tag"`
}

// EndpointSuggestion represents an endpoint suggestion.
type EndpointSuggestion struct {
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// EnvironmentVariableSuggestion represents an env var suggestion.
type EnvironmentVariableSuggestion struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// ConfigSuggestions represents configuration suggestions.
type ConfigSuggestions struct {
	EndpointSuggestions             []EndpointSuggestion            `json:"endpointSuggestions,omitempty"`
	EnvironmentVariablesSuggestions []EnvironmentVariableSuggestion `json:"environmentVariablesSuggestions,omitempty"`
	AppName                         string                          `json:"appName,omitempty"`
	Description                     string                          `json:"description,omitempty"`
	Instructions                    string                          `json:"instructions,omitempty"`
	RegistryURL                     string                          `json:"registryUrl,omitempty"`
}

// SearchPublicImagesRequest represents a request to search public images.
type SearchPublicImagesRequest struct {
	RegistryID string `json:"registryId"`
	Prefix     string `json:"prefix"`
	Size       int    `json:"size,omitempty"`
	Page       int    `json:"page,omitempty"`
}

// =============================================================================
// Container Template Types
// =============================================================================

// EnvironmentVariable represents an environment variable.
type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EntryPoint represents container entry point configuration.
type EntryPoint struct {
	Command          []string `json:"command,omitempty"`
	Arguments        []string `json:"arguments,omitempty"`
	WorkingDirectory string   `json:"workingDirectory,omitempty"`
}

// ProbeHTTPGet represents HTTP probe configuration.
type ProbeHTTPGet struct {
	Path   string `json:"path,omitempty"`
	Port   int    `json:"port,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}

// ProbeTCPSocket represents TCP probe configuration.
type ProbeTCPSocket struct {
	Port int `json:"port,omitempty"`
}

// ProbeExec represents exec probe configuration.
type ProbeExec struct {
	Command []string `json:"command,omitempty"`
}

// Probe represents a container probe.
type Probe struct {
	HTTPGet             *ProbeHTTPGet   `json:"httpGet,omitempty"`
	TCPSocket           *ProbeTCPSocket `json:"tcpSocket,omitempty"`
	Exec                *ProbeExec      `json:"exec,omitempty"`
	InitialDelaySeconds int             `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int             `json:"periodSeconds,omitempty"`
	TimeoutSeconds      int             `json:"timeoutSeconds,omitempty"`
	SuccessThreshold    int             `json:"successThreshold,omitempty"`
	FailureThreshold    int             `json:"failureThreshold,omitempty"`
}

// Probes represents all container probes.
type Probes struct {
	Startup   *Probe `json:"startup,omitempty"`
	Readiness *Probe `json:"readiness,omitempty"`
	Liveness  *Probe `json:"liveness,omitempty"`
}

// VolumeMount represents a volume mount configuration.
type VolumeMount struct {
	VolumeName string `json:"volumeName,omitempty"`
	MountPath  string `json:"mountPath,omitempty"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
}

// VolumeMountRequest represents a volume mount in requests.
type VolumeMountRequest struct {
	VolumeName string `json:"volumeName"`
	MountPath  string `json:"mountPath"`
	ReadOnly   bool   `json:"readOnly,omitempty"`
}

// ContainerEndpoint represents an endpoint attached to container.
type ContainerEndpoint struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Type        string `json:"type,omitempty"`
}

// ContainerTemplate represents a container template.
type ContainerTemplate struct {
	ID                   string                `json:"id,omitempty"`
	Name                 string                `json:"name"`
	PackageID            string                `json:"packageId,omitempty"`
	Image                string                `json:"image,omitempty"`
	ImageName            string                `json:"imageName"`
	ImageNamespace       string                `json:"imageNamespace"`
	ImageTag             string                `json:"imageTag"`
	ImageRegistryID      string                `json:"imageRegistryId"`
	ImageDigest          string                `json:"imageDigest,omitempty"`
	ImagePullPolicy      ImagePullPolicy       `json:"imagePullPolicy,omitempty"`
	EntryPoint           *EntryPoint           `json:"entryPoint,omitempty"`
	Probes               *Probes               `json:"probes,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
	Endpoints            []ContainerEndpoint   `json:"endpoints,omitempty"`
	VolumeMounts         []VolumeMount         `json:"volumeMounts,omitempty"`
}

// CreateContainerTemplateRequest represents a request to create a container template.
type CreateContainerTemplateRequest struct {
	Name                 string                `json:"name"`
	Image                string                `json:"image,omitempty"`
	ImageName            string                `json:"imageName"`
	ImageNamespace       string                `json:"imageNamespace"`
	ImageTag             string                `json:"imageTag"`
	ImageRegistryID      string                `json:"imageRegistryId"`
	ImageDigest          string                `json:"imageDigest,omitempty"`
	ImagePullPolicy      ImagePullPolicy       `json:"imagePullPolicy,omitempty"`
	EntryPoint           *EntryPoint           `json:"entryPoint,omitempty"`
	Probes               *Probes               `json:"probes,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
	Endpoints            []EndpointRequest     `json:"endpoints,omitempty"`
	VolumeMounts         []VolumeMountRequest  `json:"volumeMounts,omitempty"`
}

// PatchContainerTemplateRequest represents a request to patch a container template.
type PatchContainerTemplateRequest struct {
	Name                 string                `json:"name,omitempty"`
	Image                string                `json:"image,omitempty"`
	ImageName            string                `json:"imageName,omitempty"`
	ImageNamespace       string                `json:"imageNamespace,omitempty"`
	ImageTag             string                `json:"imageTag,omitempty"`
	ImageDigest          string                `json:"imageDigest,omitempty"`
	ImageRegistryID      string                `json:"imageRegistryId,omitempty"`
	ImagePullPolicy      ImagePullPolicy       `json:"imagePullPolicy,omitempty"`
	EntryPoint           *EntryPoint           `json:"entryPoint,omitempty"`
	Probes               *Probes               `json:"probes,omitempty"`
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables,omitempty"`
	Endpoints            []EndpointRequest     `json:"endpoints,omitempty"`
	VolumeMounts         []VolumeMountRequest  `json:"volumeMounts,omitempty"`
}

// SetEnvironmentVariablesRequest represents env vars as key-value map.
type SetEnvironmentVariablesRequest map[string]string

// =============================================================================
// Endpoint Types
// =============================================================================

// PortMapping represents a port mapping configuration.
type PortMapping struct {
	ContainerPort int        `json:"containerPort"`
	ExposedPort   int        `json:"exposedPort,omitempty"`
	Protocols     []Protocol `json:"protocols,omitempty"`
}

// StickySessions represents sticky session configuration.
type StickySessions struct {
	Enabled        bool     `json:"enabled,omitempty"`
	SessionHeaders []string `json:"sessionHeaders,omitempty"`
	CookieName     string   `json:"cookieName,omitempty"`
}

// Endpoint represents an application endpoint.
type Endpoint struct {
	ID                  string          `json:"id,omitempty"`
	DisplayName         string          `json:"displayName,omitempty"`
	PublicHost          string          `json:"publicHost,omitempty"`
	Type                EndpointType    `json:"type,omitempty"`
	IsSslEnabled        bool            `json:"isSslEnabled,omitempty"`
	PullZoneID          string          `json:"pullZoneId,omitempty"`
	PortMappings        []PortMapping   `json:"portMappings,omitempty"`
	ContainerName       string          `json:"containerName,omitempty"`
	ContainerID         string          `json:"containerId,omitempty"`
	StickySessions      *StickySessions `json:"stickySessions,omitempty"`
	InternalIPAddresses []string        `json:"internalIpAddresses,omitempty"`
	PublicIPAddresses   []string        `json:"publicIpAddresses,omitempty"`
}

// EndpointListResponse represents the response from listing endpoints.
type EndpointListResponse struct {
	Items  []Endpoint     `json:"items"`
	Meta   PaginationMeta `json:"meta"`
	Cursor string         `json:"cursor,omitempty"`
}

// CDNEndpointConfig represents CDN endpoint configuration.
type CDNEndpointConfig struct {
	IsSslEnabled   bool            `json:"isSslEnabled,omitempty"`
	StickySessions *StickySessions `json:"stickySessions,omitempty"`
	PullZoneID     int64           `json:"pullZoneId,omitempty"`
	PortMappings   []PortMapping   `json:"portMappings,omitempty"`
}

// AnycastEndpointConfig represents Anycast endpoint configuration.
type AnycastEndpointConfig struct {
	Type         AnycastType   `json:"type,omitempty"`
	PortMappings []PortMapping `json:"portMappings,omitempty"`
}

// EndpointRequest represents a request to create/update an endpoint.
type EndpointRequest struct {
	DisplayName string                 `json:"displayName"`
	CDN         *CDNEndpointConfig     `json:"cdn,omitempty"`
	Anycast     *AnycastEndpointConfig `json:"anycast,omitempty"`
}

// EndpointIDResponse represents a response with endpoint ID.
type EndpointIDResponse struct {
	ID string `json:"id"`
}

// =============================================================================
// Region Types
// =============================================================================

// Region represents a deployment region.
type Region struct {
	ID                string `json:"id"`
	Name              string `json:"name,omitempty"`
	Group             string `json:"group,omitempty"`
	HasAnycastSupport bool   `json:"hasAnycastSupport,omitempty"`
	HasCapacity       bool   `json:"hasCapacity,omitempty"`
}

// RegionListResponse represents the response from listing regions.
type RegionListResponse struct {
	Items  []Region       `json:"items"`
	Meta   PaginationMeta `json:"meta"`
	Cursor string         `json:"cursor,omitempty"`
}

// OptimalRegionResponse represents the optimal region response.
type OptimalRegionResponse struct {
	Region *Region `json:"region,omitempty"`
}

// UpdateRegionSettingsRequest represents a request to update region settings.
type UpdateRegionSettingsRequest struct {
	AllowedRegionIds  []string               `json:"allowedRegionIds,omitempty"`
	RequiredRegionIds []string               `json:"requiredRegionIds,omitempty"`
	MaxAllowedRegions int                    `json:"maxAllowedRegions,omitempty"`
	NodeSelectors     map[string]any `json:"nodeSelectors,omitempty"`
}

// =============================================================================
// Volume Types
// =============================================================================

// VolumeInstance represents an instance of a volume.
type VolumeInstance struct {
	ID                 string       `json:"id,omitempty"`
	AttachedPods       []string     `json:"attachedPods,omitempty"`
	AttachedContainers []string     `json:"attachedContainers,omitempty"`
	Region             string       `json:"region,omitempty"`
	Status             VolumeStatus `json:"status,omitempty"`
	Size               int          `json:"size,omitempty"`
	Usage              float64      `json:"usage,omitempty"`
}

// Volume represents a volume.
type Volume struct {
	ID                     string           `json:"id,omitempty"`
	Name                   string           `json:"name"`
	Size                   int              `json:"size"`
	TotalUsage             float64          `json:"totalUsage,omitempty"`
	TotalInstancesCount    int              `json:"totalInstancesCount,omitempty"`
	AttachedInstancesCount int              `json:"attachedInstancesCount,omitempty"`
	ContainersCount        int              `json:"containersCount,omitempty"`
	VolumeInstances        []VolumeInstance `json:"volumeInstances,omitempty"`
}

// VolumeSummary represents volume list summary.
type VolumeSummary struct {
	TotalPods       int `json:"totalPods,omitempty"`
	TotalContainers int `json:"totalContainers,omitempty"`
	TotalStorage    int `json:"totalStorage,omitempty"`
}

// VolumeListResponse represents the response from listing volumes.
type VolumeListResponse struct {
	Items   []Volume       `json:"items"`
	Meta    PaginationMeta `json:"meta"`
	Cursor  string         `json:"cursor,omitempty"`
	Summary *VolumeSummary `json:"summary,omitempty"`
}

// UpdateVolumeRequest represents a request to update a volume.
type UpdateVolumeRequest struct {
	Name string `json:"name,omitempty"`
	Size int    `json:"size,omitempty"`
}

// VolumeUpdateResponse represents volume update response.
type VolumeUpdateResponse struct {
	Name string `json:"name,omitempty"`
	Size int    `json:"size,omitempty"`
}

// VolumeNameResponse represents response with volume name.
type VolumeNameResponse struct {
	Name string `json:"name,omitempty"`
}

// VolumeInstanceIDResponse represents response with volume instance ID.
type VolumeInstanceIDResponse struct {
	ID string `json:"id,omitempty"`
}

// VolumeInstanceIDsResponse represents response with volume instance IDs.
type VolumeInstanceIDsResponse struct {
	IDs []string `json:"ids,omitempty"`
}

// =============================================================================
// Log Forwarding Types
// =============================================================================

// LogForwardingConfig represents a log forwarding configuration.
type LogForwardingConfig struct {
	ID        string              `json:"id,omitempty"`
	App       string              `json:"app,omitempty"`
	ProductID string              `json:"productId,omitempty"`
	Type      LogForwardingType   `json:"type"`
	Endpoint  string              `json:"endpoint"`
	Port      int                 `json:"port"`
	CreatedAt string              `json:"createdAt,omitempty"`
	Token     string              `json:"token,omitempty"`
	Format    LogForwardingFormat `json:"format"`
	Enabled   bool                `json:"enabled"`
}

// LogForwardingListResponse represents the response from listing log forwarding configs.
type LogForwardingListResponse struct {
	Items []LogForwardingConfig `json:"items"`
}

// CreateLogForwardingRequest represents a request to create log forwarding.
type CreateLogForwardingRequest struct {
	App      string              `json:"app"`
	Type     LogForwardingType   `json:"type"`
	Endpoint string              `json:"endpoint"`
	Port     int                 `json:"port"`
	Token    string              `json:"token,omitempty"`
	Format   LogForwardingFormat `json:"format"`
	Enabled  bool                `json:"enabled"`
}

// UpdateLogForwardingRequest represents a request to update log forwarding.
type UpdateLogForwardingRequest struct {
	App      string              `json:"app"`
	Type     LogForwardingType   `json:"type"`
	Endpoint string              `json:"endpoint"`
	Port     int                 `json:"port"`
	Token    string              `json:"token,omitempty"`
	Format   LogForwardingFormat `json:"format"`
	Enabled  bool                `json:"enabled"`
}

// =============================================================================
// Nodes, Pods, Limits Types
// =============================================================================

// NodeListResponse represents the response from listing nodes.
type NodeListResponse struct {
	Items  []string       `json:"items"`
	Meta   PaginationMeta `json:"meta"`
	Cursor string         `json:"cursor,omitempty"`
}

// UserLimits represents user limits.
type UserLimits struct {
	MaxNumberOfApplications            int `json:"maxNumberOfApplications,omitempty"`
	ExistingNumberOfApplications       int `json:"existingNumberOfApplications,omitempty"`
	MaxNumberOfRegionsPerApplication   int `json:"maxNumberOfRegionsPerApplication,omitempty"`
	MaxNumberOfInstancesPerRegion      int `json:"maxNumberOfInstancesPerRegion,omitempty"`
	MaxNumberOfInstancesPerApplication int `json:"maxNumberOfInstancesPerApplication,omitempty"`
	MaxNumberOfVolumesPerApplication   int `json:"maxNumberOfVolumesPerApplication,omitempty"`
	MaxVolumeSize                      int `json:"maxVolumeSize,omitempty"`
}
