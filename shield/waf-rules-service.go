package shield

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// WAFService provides methods for managing WAF rules and configuration.
type WAFService interface {
	// Rules listing
	ListRules(ctx context.Context) (*WAFRuleListResponse, error)
	ListCustomRules(ctx context.Context) (*CustomRuleListResponse, error)

	// Custom rule CRUD
	CreateCustomRule(ctx context.Context, req *CreateCustomRuleRequest) (*CustomRule, error)
	GetCustomRule(ctx context.Context, ruleID string) (*CustomRule, error)
	UpdateCustomRule(ctx context.Context, ruleID string, req *UpdateCustomRuleRequest) (*CustomRule, error)
	ReplaceCustomRule(ctx context.Context, ruleID string, req *ReplaceCustomRuleRequest) (*CustomRule, error)
	DeleteCustomRule(ctx context.Context, ruleID string) error

	// Configuration
	GetProfiles(ctx context.Context) (*WAFProfilesResponse, error)
	GetEngineConfig(ctx context.Context) (*WAFEngineConfig, error)
	GetEnums(ctx context.Context) (*WAFEnums, error)

	// Rule review
	GetTriggeredRules(ctx context.Context) (*TriggeredRulesResponse, error)
	SubmitTriggeredRuleReview(ctx context.Context, req *TriggeredRuleReviewRequest) (*TriggeredRuleReview, error)
	GetAIRecommendation(ctx context.Context, ruleID string) (*AIRecommendationResponse, error)

	// Plan info
	GetPlanSegmentation(ctx context.Context) (*PlanSegmentationResponse, error)
}

type wafService struct {
	client httpClient
}

func newWAFService(client httpClient) WAFService {
	return &wafService{client: client}
}

// ListRules returns all WAF rules (predefined and custom).
func (s *wafService) ListRules(ctx context.Context) (*WAFRuleListResponse, error) {
	var resp WAFRuleListResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/rules", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCustomRules returns all custom WAF rules.
func (s *wafService) ListCustomRules(ctx context.Context) (*CustomRuleListResponse, error) {
	var resp CustomRuleListResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/custom-rules", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCustomRule creates a new custom WAF rule.
func (s *wafService) CreateCustomRule(ctx context.Context, req *CreateCustomRuleRequest) (*CustomRule, error) {
	var rule CustomRule
	if err := s.client.do(ctx, http.MethodPost, "/shield/waf/custom-rule", req, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetCustomRule returns a specific custom WAF rule by ID.
func (s *wafService) GetCustomRule(ctx context.Context, ruleID string) (*CustomRule, error) {
	path := fmt.Sprintf("/shield/waf/custom-rule/%s", ruleID)
	var rule CustomRule
	if err := s.client.do(ctx, http.MethodGet, path, nil, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateCustomRule updates a custom WAF rule (PATCH - partial update).
func (s *wafService) UpdateCustomRule(ctx context.Context, ruleID string, req *UpdateCustomRuleRequest) (*CustomRule, error) {
	path := fmt.Sprintf("/shield/waf/custom-rule/%s", ruleID)
	var rule CustomRule
	if err := s.client.do(ctx, http.MethodPatch, path, req, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// ReplaceCustomRule replaces a custom WAF rule (PUT - full replacement).
func (s *wafService) ReplaceCustomRule(ctx context.Context, ruleID string, req *ReplaceCustomRuleRequest) (*CustomRule, error) {
	path := fmt.Sprintf("/shield/waf/custom-rule/%s", ruleID)
	var rule CustomRule
	if err := s.client.do(ctx, http.MethodPut, path, req, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteCustomRule deletes a custom WAF rule.
func (s *wafService) DeleteCustomRule(ctx context.Context, ruleID string) error {
	path := fmt.Sprintf("/shield/waf/custom-rule/%s", ruleID)
	return s.client.do(ctx, http.MethodDelete, path, nil, nil)
}

// GetProfiles returns available WAF profiles.
func (s *wafService) GetProfiles(ctx context.Context) (*WAFProfilesResponse, error) {
	var resp WAFProfilesResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/profiles", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetEngineConfig returns WAF engine configuration.
func (s *wafService) GetEngineConfig(ctx context.Context) (*WAFEngineConfig, error) {
	var resp WAFEngineConfig
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/engine-config", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetEnums returns WAF enumeration values.
func (s *wafService) GetEnums(ctx context.Context) (*WAFEnums, error) {
	var resp WAFEnums
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/enums", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTriggeredRules returns WAF rules that were triggered and need review.
func (s *wafService) GetTriggeredRules(ctx context.Context) (*TriggeredRulesResponse, error) {
	var resp TriggeredRulesResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/rules/review-triggered", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubmitTriggeredRuleReview submits a review for a triggered WAF rule.
func (s *wafService) SubmitTriggeredRuleReview(ctx context.Context, req *TriggeredRuleReviewRequest) (*TriggeredRuleReview, error) {
	var resp TriggeredRuleReview
	if err := s.client.do(ctx, http.MethodPost, "/shield/waf/rules/review-triggered", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAIRecommendation returns AI-powered recommendations for triggered rules.
func (s *wafService) GetAIRecommendation(ctx context.Context, ruleID string) (*AIRecommendationResponse, error) {
	path := "/shield/waf/rules/review-triggered/ai-recommendation"
	if ruleID != "" {
		params := url.Values{}
		params.Set("ruleId", ruleID)
		path = path + "?" + params.Encode()
	}
	var resp AIRecommendationResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPlanSegmentation returns WAF rule plan segmentation details.
func (s *wafService) GetPlanSegmentation(ctx context.Context) (*PlanSegmentationResponse, error) {
	var resp PlanSegmentationResponse
	if err := s.client.do(ctx, http.MethodGet, "/shield/waf/rules/plan-segmentation", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
