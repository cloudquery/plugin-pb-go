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
	errNotFound        = errors.New("not found")
	errShortRead       = errors.New("truncated response body")
	errDownloadStalled = errors.New("download stalled")
)

// Overridable so tests do not pay the real backoff.
var (
	downloadRetryAttempts = uint(RetryAttempts)
	downloadRetryDelay    = RetryWaitTime
	downloadRetryMaxDelay = MaxRetryWaitTime
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
// net/http transport, and on Windows the syscall.E* constants are synthetic
// values that never match a real WSA socket error, so these failures can only be
// matched on their message.
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
	"connection refused",
	"no such host",
	"network is unreachable",
	"no route to host",
	// Windows WSAECONNRESET, WSAECONNREFUSED and WSAETIMEDOUT respectively.
	"forcibly closed by the remote host",
	"actively refused it",
	"did not properly respond after a period of time",
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
	case errors.Is(err, errDownloadStalled),
		errors.Is(err, errShortRead),
		errors.Is(err, io.ErrUnexpectedEOF),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.ECONNRESET),
		errors.Is(err, syscall.ECONNREFUSED),
		errors.Is(err, syscall.EPIPE),
		errors.Is(err, syscall.ETIMEDOUT),
		errors.Is(err, syscall.EHOSTUNREACH),
		errors.Is(err, syscall.ENETUNREACH):
		return true
	}

	// The asset host is fixed, so a resolution failure against it is a resolver
	// problem rather than a bad name.
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
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

// downloadTimeoutError re-labels a timeout produced by our own transport, which
// otherwise reaches the classifier indistinguishable from the caller's deadline:
// both satisfy errors.Is(err, context.DeadlineExceeded), and only the caller's is
// terminal. The underlying error is formatted with %v so that shared
// context.DeadlineExceeded does not travel on in the chain. Returns nil when the
// caller's own context ended the request, or when err is not a timeout at all.
func downloadTimeoutError(ctx context.Context, urlForLog string, err error) error {
	if ctx.Err() != nil {
		return nil
	}
	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		return nil
	}
	return fmt.Errorf("%w: no data from %s: %v", errDownloadStalled, urlForLog, redactURLError(err))
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

// redactURLError rewrites the URL that *url.Error prints verbatim. Wrapping such
// an error re-exposes the signed token that redactURLQuery removed from the
// surrounding message.
func redactURLError(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		urlErr.URL = redactURLQuery(urlErr.URL)
	}
	return err
}

// IsTransientDownloadError reports whether err is a transient plugin download
// failure - a network or server-side problem rather than a bad plugin reference.
// Callers use it to keep advice about plugin resolution off errors that have
// nothing to do with it.
func IsTransientDownloadError(err error) bool {
	return isRetryableDownloadError(err)
}
