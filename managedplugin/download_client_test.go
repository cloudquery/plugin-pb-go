package managedplugin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// fastStall shrinks the transport timeouts so the tests do not wait the real 30s.
func fastStall(t *testing.T, d time.Duration) {
	t.Helper()

	prevHeader, prevIdle, prevClient := downloadResponseHeaderTimeout, downloadIdleReadTimeout, downloadClient
	downloadResponseHeaderTimeout, downloadIdleReadTimeout = d, d
	downloadClient = newDownloadClient()
	t.Cleanup(func() {
		downloadResponseHeaderTimeout, downloadIdleReadTimeout = prevHeader, prevIdle
		downloadClient = prevClient
	})
}

func TestDownloadFileRetriesStalledHeaders(t *testing.T) {
	fastRetries(t)
	fastStall(t, 150*time.Millisecond)

	body := []byte("payload")
	release := make(chan struct{})

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			<-release
			return
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	checksum, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.NoError(t, err)
	require.GreaterOrEqual(t, attempts.Load(), int32(2))
	require.Equal(t, sha256Hex(body), checksum)
}

func TestDownloadFileRetriesStalledBody(t *testing.T) {
	fastRetries(t)
	fastStall(t, 150*time.Millisecond)

	body := []byte("cloudquery-plugin-binary-payload")
	release := make(chan struct{})

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body[:5])
			w.(http.Flusher).Flush()
			<-release
			return
		}
		_, _ = w.Write(body)
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	checksum, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.NoError(t, err)
	require.GreaterOrEqual(t, attempts.Load(), int32(2))

	written, err := os.ReadFile(localPath)
	require.NoError(t, err)
	require.Equal(t, body, written)
	require.Equal(t, sha256Hex(body), checksum)
}

func TestDownloadFileSlowButProgressingSucceeds(t *testing.T) {
	fastRetries(t)
	fastStall(t, 300*time.Millisecond)

	body := []byte("slow-but-moving-plugin-payload")

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(http.StatusOK)
		for _, chunk := range [][]byte{body[:6], body[6:12], body[12:18], body[18:24], body[24:]} {
			_, _ = w.Write(chunk)
			w.(http.Flusher).Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}))
	t.Cleanup(server.Close)

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	checksum, err := downloadFile(context.Background(), localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.NoError(t, err)
	require.EqualValues(t, 1, attempts.Load(), "a transfer that keeps making progress must not trip the idle deadline")
	require.Equal(t, sha256Hex(body), checksum)
}

func TestDownloadFileCallerCancelIsNotRetried(t *testing.T) {
	fastRetries(t)

	release := make(chan struct{})

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Length", "32")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("early"))
		w.(http.Flusher).Flush()
		<-release
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err := downloadFile(ctx, localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	require.False(t, IsTransientDownloadError(err), "the caller's own cancellation must not be retried")
}

func TestDownloadFileCallerDeadlineIsNotRetried(t *testing.T) {
	fastRetries(t)
	fastStall(t, 10*time.Second)

	release := make(chan struct{})

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		<-release
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	localPath := filepath.Join(t.TempDir(), "plugin.zip")
	_, err := downloadFile(ctx, localPath, server.URL, DownloaderOptions{NoProgress: true})
	require.Error(t, err)
	require.NotErrorIs(t, err, errDownloadStalled)
	require.False(t, IsTransientDownloadError(err), "the caller's own deadline must not be retried")
}
