package github

import (
	"net/http"
	"strconv"
	"time"
)

func ParseRetryAfter(resp *http.Response) time.Duration {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if sec, err := strconv.Atoi(v); err == nil {
		return time.Duration(sec) * time.Second
	}
	if t, err := time.Parse(http.TimeFormat, v); err == nil {
		return time.Until(t)
	}
	return 0
}
