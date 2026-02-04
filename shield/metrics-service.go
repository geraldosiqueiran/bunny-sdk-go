package shield

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// MetricsService provides methods for accessing Shield metrics.
type MetricsService interface {
	GetOverview(ctx context.Context, opts *DateRangeOptions) (*MetricsOverview, error)
	GetOverviewDetailed(ctx context.Context, opts *MetricsDetailedOptions) (*MetricsOverviewDetailed, error)
	GetWAFRuleMetrics(ctx context.Context, ruleID string, opts *DateRangeOptions) (*WAFRuleMetrics, error)
	GetRateLimitMetrics(ctx context.Context, rateLimitID string, opts *DateRangeOptions) (*RateLimitMetrics, error)
	GetAllRateLimitMetrics(ctx context.Context, opts *DateRangeOptions) (*RateLimitMetricsList, error)
	GetBotDetectionMetrics(ctx context.Context, zoneID string, opts *DateRangeOptions) (*BotDetectionMetrics, error)
	GetUploadScanningMetrics(ctx context.Context, zoneID string, opts *DateRangeOptions) (*UploadScanningMetrics, error)
}

type metricsService struct {
	client httpClient
}

func newMetricsService(client httpClient) MetricsService {
	return &metricsService{client: client}
}

// GetOverview returns overview metrics for all Shield zones.
func (s *metricsService) GetOverview(ctx context.Context, opts *DateRangeOptions) (*MetricsOverview, error) {
	path := "/shield/metrics/overview" + buildDateRangeQuery(opts)
	var resp MetricsOverview
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOverviewDetailed returns detailed overview metrics with breakdown.
func (s *metricsService) GetOverviewDetailed(ctx context.Context, opts *MetricsDetailedOptions) (*MetricsOverviewDetailed, error) {
	path := "/shield/metrics/overview-detailed" + buildMetricsDetailedQuery(opts)
	var resp MetricsOverviewDetailed
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWAFRuleMetrics returns metrics for a specific WAF rule.
func (s *metricsService) GetWAFRuleMetrics(ctx context.Context, ruleID string, opts *DateRangeOptions) (*WAFRuleMetrics, error) {
	path := fmt.Sprintf("/shield/metrics/waf-rule/%s", ruleID) + buildDateRangeQuery(opts)
	var resp WAFRuleMetrics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRateLimitMetrics returns metrics for a specific rate limit rule.
func (s *metricsService) GetRateLimitMetrics(ctx context.Context, rateLimitID string, opts *DateRangeOptions) (*RateLimitMetrics, error) {
	path := fmt.Sprintf("/shield/metrics/rate-limit/%s", rateLimitID) + buildDateRangeQuery(opts)
	var resp RateLimitMetrics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAllRateLimitMetrics returns metrics for all rate limit rules.
func (s *metricsService) GetAllRateLimitMetrics(ctx context.Context, opts *DateRangeOptions) (*RateLimitMetricsList, error) {
	path := "/shield/metrics/rate-limits" + buildDateRangeQuery(opts)
	var resp RateLimitMetricsList
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBotDetectionMetrics returns bot detection metrics for a zone.
func (s *metricsService) GetBotDetectionMetrics(ctx context.Context, zoneID string, opts *DateRangeOptions) (*BotDetectionMetrics, error) {
	path := fmt.Sprintf("/shield/metrics/bot-detection/%s", zoneID) + buildDateRangeQuery(opts)
	var resp BotDetectionMetrics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUploadScanningMetrics returns upload scanning metrics for a zone.
func (s *metricsService) GetUploadScanningMetrics(ctx context.Context, zoneID string, opts *DateRangeOptions) (*UploadScanningMetrics, error) {
	path := fmt.Sprintf("/shield/metrics/upload-scanning/%s", zoneID) + buildDateRangeQuery(opts)
	var resp UploadScanningMetrics
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildDateRangeQuery(opts *DateRangeOptions) string {
	if opts == nil {
		return ""
	}
	params := url.Values{}
	if opts.From != "" {
		params.Set("from", opts.From)
	}
	if opts.To != "" {
		params.Set("to", opts.To)
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}

func buildMetricsDetailedQuery(opts *MetricsDetailedOptions) string {
	if opts == nil {
		return ""
	}
	params := url.Values{}
	if opts.From != "" {
		params.Set("from", opts.From)
	}
	if opts.To != "" {
		params.Set("to", opts.To)
	}
	if opts.ZoneID != "" {
		params.Set("zoneId", opts.ZoneID)
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}
