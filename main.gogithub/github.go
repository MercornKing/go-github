package github

import (
	"net/http"
	"strconv"
	"time"
)

// AbuseRateLimitError occurs when GitHub returns a 403 Forbidden response with the
// "retry-after" header.
type AbuseRateLimitError struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"`

	// RetryAfter is the parsed duration from the Retry-After header.
	// It is nil if the header is missing or could not be parsed.
	RetryAfter *time.Duration
}

// Error implements the error interface.
func (r *AbuseRateLimitError) Error() string {
	return r.Message
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusForbidden && r.Header.Get("Retry-After") != "" {
		err := &AbuseRateLimitError{
			Response:   r,
			Message:    "You have exceeded a secondary rate limit.",
			RetryAfter: parseRetryAfter(r.Header.Get("Retry-After")),
		}
		return err
	}
	return nil
}

// parseRetryAfter parses the Retry-After header, attempting to parse it as an
// integer (seconds) or as an HTTP-date.
func parseRetryAfter(header string) *time.Duration {
	if header == "" {
		return nil
	}

	// 1. Try to parse as integer (seconds)
	if v, err := strconv.ParseInt(header, 10, 64); err == nil {
		d := time.Duration(v) * time.Second
		return &d
	}

	// 2. Try to parse as HTTP-date (RFC 1123 format)
	if t, err := time.Parse(http.TimeFormat, header); err == nil {
		d := time.Until(t)
		if d < 0 {
			d = 0 // If the time has already passed, treat as zero duration
		}
		return &d
	}

	// 3. Graceful degradation on parse failure
	return nil
}
