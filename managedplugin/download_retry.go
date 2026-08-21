package managedplugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
)

var (
	errNotFound  = errors.New("not found")
	errShortRead = errors.New("truncated response body")
)

type httpStatusError struct {
	statusCode int
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("statusCode %d", e.statusCode)
}

func isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout, http.StatusTooManyRequests:
		return true
	}
	return statusCode >= http.StatusInternalServerError
}

// Go does not export the HTTP/2 stream and connection error types used by the
// net/http transport, so the mid-body resets we get from the asset CDN can only
// be matched on their message.
var transientTransportMessages = []string{
	"stream error",
	"server sent goaway",
	"connection reset by peer",
	"broken pipe",
	"unexpected eof",
	"use of closed network connection",
	"server closed idle connection",
	"transport connection broken",
	"i/o timeout",
}

func isRetryableDownloadError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, errNotFound) {
		return false
	}

	var statusErr *httpStatusError
	if errors.As(err, &statusErr) {
		return isRetryableStatusCode(statusErr.statusCode)
	}

	switch {
	case errors.Is(err, errShortRead),
		errors.Is(err, io.ErrUnexpectedEOF),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.ECONNRESET),
		errors.Is(err, syscall.EPIPE),
		errors.Is(err, syscall.ETIMEDOUT):
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	msg := strings.ToLower(err.Error())
	for _, transient := range transientTransportMessages {
		if strings.Contains(msg, transient) {
			return true
		}
	}
	return false
}

// redactURLQuery strips the query string so that the signed download token never
// reaches stdout or a log aggregator.
func redactURLQuery(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

// IsTransientDownloadError reports whether err is a transient plugin download
// failure - a network or server-side problem rather than a bad plugin reference.
// Callers use it to keep advice about plugin resolution off errors that have
// nothing to do with it.
func IsTransientDownloadError(err error) bool {
	return isRetryableDownloadError(err)
}
