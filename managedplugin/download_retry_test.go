package managedplugin

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsRetryableDownloadError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{
			name: "http2 mid-body internal error stream reset",
			err:  fmt.Errorf("failed to copy body to file plugin.zip: %w", errors.New("stream error: stream ID 1; INTERNAL_ERROR; received from peer")),
			want: true,
		},
		{
			name: "http2 refused stream",
			err:  errors.New("stream error: stream ID 3; REFUSED_STREAM"),
			want: true,
		},
		{
			name: "http2 goaway",
			err:  errors.New("http2: server sent GOAWAY and closed the connection"),
			want: true,
		},
		{name: "connection reset", err: fmt.Errorf("read tcp: %w", syscall.ECONNRESET), want: true},
		{name: "broken pipe", err: fmt.Errorf("write tcp: %w", syscall.EPIPE), want: true},
		{name: "unexpected EOF", err: fmt.Errorf("failed to copy body to file: %w", io.ErrUnexpectedEOF), want: true},
		{name: "bare EOF", err: io.EOF, want: true},
		{name: "short read", err: fmt.Errorf("%w: got 10 bytes, want 20", errShortRead), want: true},
		{name: "net timeout", err: &net.OpError{Op: "read", Err: os.ErrDeadlineExceeded}, want: true},

		{name: "expired signed URL 403", err: &httpStatusError{statusCode: http.StatusForbidden}, want: false},
		{name: "unauthorized 401", err: &httpStatusError{statusCode: http.StatusUnauthorized}, want: false},
		{name: "not found sentinel", err: errNotFound, want: false},
		{name: "not found 404 status", err: &httpStatusError{statusCode: http.StatusNotFound}, want: false},
		{name: "bad request 400", err: &httpStatusError{statusCode: http.StatusBadRequest}, want: false},

		{name: "request timeout 408", err: &httpStatusError{statusCode: http.StatusRequestTimeout}, want: true},
		{name: "too many requests 429", err: &httpStatusError{statusCode: http.StatusTooManyRequests}, want: true},
		{name: "bad gateway 502", err: &httpStatusError{statusCode: http.StatusBadGateway}, want: true},
		{name: "service unavailable 503", err: &httpStatusError{statusCode: http.StatusServiceUnavailable}, want: true},

		{name: "context canceled", err: fmt.Errorf("get url: %w", context.Canceled), want: false},
		{name: "context deadline exceeded", err: fmt.Errorf("get url: %w", context.DeadlineExceeded), want: false},
		{name: "checksum mismatch is permanent", err: errors.New("checksum mismatch: expected abc, got def"), want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, isRetryableDownloadError(tc.err))
		})
	}
}

func TestRedactURLQuery(t *testing.T) {
	require.Equal(t,
		"https://assets.cloudquery.io/cq-cloud-releases/cloudquery/source/envzero/v2.1.0/linux_amd64",
		redactURLQuery("https://assets.cloudquery.io/cq-cloud-releases/cloudquery/source/envzero/v2.1.0/linux_amd64?verify=1787270563-PVHLM7Vfma5I0Mzx1YZXt4hPqXydx5WrqxPQf7B80I8%3D"),
	)
}

// TestDownloadFileRetriesTruncatedBody reproduces the production failure: the first
// attempt writes part of the body and then the connection drops mid-copy. The retry
// must start the file from scratch rather than append to the partial bytes.
func TestDownloadFileRetriesTruncatedBody(t *testing.T) {
	body := []byte("cloudquery-plugin-binary-payload")

	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body[:5])
			w.(http.Flusher).Flush()
			// Close the connection mid-body so the client sees a truncated response.
			conn, _, err := w.(http.Hijacker).Hijack()
			require.NoError(t, err)
			conn.Close()
			return
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	checksum, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.NoError(t, err)
	require.Equal(t, 2, attempts)

	written, err := os.ReadFile(localPath)
	require.NoError(t, err)
	require.Equal(t, body, written, "retry must not append to the partial first attempt")
	require.Equal(t, sha256Hex(body), checksum)
}

func TestDownloadFileDoesNotRetryExpiredSignedURL(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.Error(t, err)
	require.Equal(t, 1, attempts, "an expired signed URL must fail fast")
	require.Contains(t, err.Error(), "statusCode 403")
}

func TestDownloadFileDoesNotRetryNotFound(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.ErrorIs(t, err, errNotFound)
	require.Equal(t, 1, attempts)
}

func TestDownloadFileRetriesServerError(t *testing.T) {
	body := []byte("payload")

	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	checksum, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.NoError(t, err)
	require.Equal(t, 3, attempts)
	require.Equal(t, sha256Hex(body), checksum)
}

func TestDownloadFileGivesUpAfterRetryAttempts(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.Error(t, err)
	require.Equal(t, RetryAttempts, attempts)
	require.Contains(t, err.Error(), "failed downloading URL")
	require.NotContains(t, err.Error(), "verify=", "the signed token must not reach the error message")
}

func sha256Hex(b []byte) string {
	s := sha256.New()
	s.Write(b)
	return fmt.Sprintf("%x", s.Sum(nil))
}
