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
	"time"

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

		{name: "stalled download", err: fmt.Errorf("%w: no data from host: i/o timeout", errDownloadStalled), want: true},
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
	fastRetries(t)

	body := []byte("cloudquery-plugin-binary-payload")

	var (
		attempts  int
		hijackErr error
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body[:5])
			w.(http.Flusher).Flush()
			// Close the connection mid-body so the client sees a truncated response.
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				hijackErr = err
				return
			}
			conn.Close()
			return
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	checksum, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.NoError(t, hijackErr)
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
	fastRetries(t)

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
	fastRetries(t)

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

// fastRetries removes the real backoff so the retry tests do not sleep through it.
func fastRetries(t *testing.T) {
	t.Helper()

	delay, maxDelay := downloadRetryDelay, downloadRetryMaxDelay
	downloadRetryDelay, downloadRetryMaxDelay = time.Millisecond, time.Millisecond
	t.Cleanup(func() {
		downloadRetryDelay, downloadRetryMaxDelay = delay, maxDelay
	})
}

// TestDownloadFileRedactsSignedTokenFromTransportErrors covers the leak that
// url.Error reopens: it prints its URL verbatim, so wrapping one puts the signed
// token back into the message once per attempt.
func TestDownloadFileRedactsSignedTokenFromTransportErrors(t *testing.T) {
	fastRetries(t)

	const token = "SUPERSECRETTOKEN"

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	server.Config.ConnState = func(c net.Conn, state http.ConnState) {
		if state == http.StateActive {
			c.Close()
		}
	}
	server.Start()
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err := downloadFile(context.Background(), localPath, server.URL+"/asset?verify="+token, DownloaderOptions{NoProgress: true})
	require.Error(t, err)
	require.NotContains(t, err.Error(), token)
	require.NotContains(t, err.Error(), "verify=")
}

func TestDownloadFileRetriesConnectionRefused(t *testing.T) {
	fastRetries(t)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err = downloadFile(context.Background(), localPath, "http://"+addr+"/asset", DownloaderOptions{NoProgress: true})
	require.Error(t, err)

	// The refusal message is platform specific, so count the attempts that the
	// retry aggregate numbers instead of matching on it.
	require.Contains(t, err.Error(), fmt.Sprintf("#%d:", downloadRetryAttempts))
}

func TestIsRetryableDownloadErrorWindowsSocketMessages(t *testing.T) {
	cases := []struct {
		name string
		msg  string
	}{
		{name: "WSAECONNRESET", msg: "wsarecv: An existing connection was forcibly closed by the remote host."},
		{name: "WSAECONNREFUSED", msg: "connectex: No connection could be made because the target machine actively refused it."},
		{name: "WSAETIMEDOUT", msg: "connectex: A connection attempt failed because the connected party did not properly respond after a period of time."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.True(t, isRetryableDownloadError(fmt.Errorf("failed to get url: %w", errors.New(tc.msg))))
		})
	}
}

func TestIsRetryableDownloadErrorDNS(t *testing.T) {
	require.True(t, isRetryableDownloadError(fmt.Errorf("dial: %w", &net.DNSError{Err: "no such host", Name: "assets.cloudquery.io", IsNotFound: true})))
}

// TestTransportErrorFromRealRequestIsRetryable pins the gap that left
// getURLLocation aborting on the first attempt: a transport failure surfaces as a
// wrapped *url.Error, which its old identity comparison could never match.
func TestTransportErrorFromRealRequestIsRetryable(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	server.Config.ConnState = func(c net.Conn, state http.ConnState) {
		if state == http.StateActive {
			c.Close()
		}
	}
	server.Start()
	t.Cleanup(server.Close)

	resp, err := http.Get(server.URL + "/asset?verify=SUPERSECRETTOKEN")
	if resp != nil {
		resp.Body.Close()
	}
	require.Error(t, err)

	wrapped := fmt.Errorf("failed to get url %s: %w", redactURLQuery(server.URL), redactURLError(err))
	require.True(t, isRetryableDownloadError(wrapped))
	require.NotContains(t, wrapped.Error(), "SUPERSECRETTOKEN")
}
