package shield

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// EventLogsService provides methods for accessing security event logs.
type EventLogsService interface {
	List(ctx context.Context, opts *EventLogListOptions) (*EventLogListResponse, error)
}

type eventLogsService struct {
	client httpClient
}

func newEventLogsService(client httpClient) EventLogsService {
	return &eventLogsService{client: client}
}

// List returns security event logs with optional filtering.
func (s *eventLogsService) List(ctx context.Context, opts *EventLogListOptions) (*EventLogListResponse, error) {
	path := "/shield/event-logs" + buildEventLogQuery(opts)
	var resp EventLogListResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildEventLogQuery(opts *EventLogListOptions) string {
	if opts == nil {
		return ""
	}
	params := url.Values{}
	if opts.ZoneID != "" {
		params.Set("zoneId", opts.ZoneID)
	}
	if opts.From != "" {
		params.Set("from", opts.From)
	}
	if opts.To != "" {
		params.Set("to", opts.To)
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}
