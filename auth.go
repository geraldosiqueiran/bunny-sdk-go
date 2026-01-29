package bunny

import "net/http"

const (
	// HeaderAccessKey is the HTTP header used for API authentication.
	HeaderAccessKey = "AccessKey"
)

// setAuthHeader sets the authentication header on the request.
func setAuthHeader(req *http.Request, apiKey string) {
	req.Header.Set(HeaderAccessKey, apiKey)
}
