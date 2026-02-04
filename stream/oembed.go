package stream

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// OEmbedService provides methods for oEmbed video embedding.
type OEmbedService interface {
	Get(ctx context.Context, opts *OEmbedOptions) (*OEmbedResponse, error)
}

type oembedService struct {
	client httpClient
}

func newOEmbedService(client httpClient) OEmbedService {
	return &oembedService{client: client}
}

// Get retrieves oEmbed data for a video.
func (s *oembedService) Get(ctx context.Context, opts *OEmbedOptions) (*OEmbedResponse, error) {
	path := "/OEmbed"
	if opts != nil {
		path = path + "?" + buildOEmbedQuery(opts)
	}

	var resp OEmbedResponse
	if err := s.client.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildOEmbedQuery(opts *OEmbedOptions) string {
	params := url.Values{}
	if opts.URL != "" {
		params.Set("url", opts.URL)
	}
	if opts.MaxWidth > 0 {
		params.Set("maxWidth", strconv.Itoa(opts.MaxWidth))
	}
	if opts.MaxHeight > 0 {
		params.Set("maxHeight", strconv.Itoa(opts.MaxHeight))
	}
	if opts.Token != "" {
		params.Set("token", opts.Token)
	}
	if opts.Expires > 0 {
		params.Set("expires", strconv.FormatInt(opts.Expires, 10))
	}
	return params.Encode()
}
