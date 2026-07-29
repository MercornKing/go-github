package github

import (
	"net/http"
	"testing"
	"time"
)

func TestCheckResponse_RetryAfterInteger(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusForbidden,
		Header:     make(http.Header),
	}
	resp.Header.Set("Retry-After", "60")

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	abuseErr, ok := err.(*AbuseRateLimitError)
	if !ok {
		t.Fatalf("Expected AbuseRateLimitError, got %T", err)
	}

	if abuseErr.RetryAfter == nil {
		t.Fatal("Expected RetryAfter to be set")
	}

	if *abuseErr.RetryAfter != 60*time.Second {
		t.Errorf("Expected RetryAfter to be 60s, got %v", *abuseErr.RetryAfter)
	}
}

func TestCheckResponse_RetryAfterHTTPDate(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusForbidden,
		Header:     make(http.Header),
	}
	
	futureTime := time.Now().Add(5 * time.Minute)
	resp.Header.Set("Retry-After", futureTime.Format(http.TimeFormat))

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	abuseErr, ok := err.(*AbuseRateLimitError)
	if !ok {
		t.Fatalf("Expected AbuseRateLimitError, got %T", err)
	}

	if abuseErr.RetryAfter == nil {
		t.Fatal("Expected RetryAfter to be set")
	}

	// Because of parsing precision, we check if it's close to 5 minutes
	diff := *abuseErr.RetryAfter - 5*time.Minute
	if diff < -2*time.Second || diff > 2*time.Second {
		t.Errorf("Expected RetryAfter to be ~5m, got %v", *abuseErr.RetryAfter)
	}
}

func TestCheckResponse_RetryAfterInvalid(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusForbidden,
		Header:     make(http.Header),
	}
	resp.Header.Set("Retry-After", "invalid-date-format")

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	abuseErr, ok := err.(*AbuseRateLimitError)
	if !ok {
		t.Fatalf("Expected AbuseRateLimitError, got %T", err)
	}

	if abuseErr.RetryAfter != nil {
		t.Errorf("Expected RetryAfter to be nil for invalid header, got %v", *abuseErr.RetryAfter)
	}
}
