package managedplugin

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// corruptZip builds a valid zip whose entry has a damaged deflate stream, so
// archive.Open succeeds and io.Copy fails part way through the entry.
func corruptZip(t *testing.T, entry string) []byte {
	t.Helper()

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	entryWriter, err := w.Create(entry)
	require.NoError(t, err)
	_, err = entryWriter.Write(bytes.Repeat([]byte("cloudquery-plugin-payload"), 4096))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	raw := buf.Bytes()
	// Damage the middle of the compressed stream, well past the local file header.
	for i := len(raw) / 2; i < len(raw)/2+64; i++ {
		raw[i] ^= 0xff
	}
	return raw
}

// TestDownloadPluginFromHubRemovesPartialBinary covers the caching trap: the hub
// path never re-downloads once a file exists at LocalPath, so an extraction that
// fails after the file is created would poison every later run. The download
// checksum cannot catch this - it is verified before extraction begins, and is
// skipped entirely when the hub returns no checksum.
func TestDownloadPluginFromHubRemovesPartialBinary(t *testing.T) {
	const (
		pluginName    = "envzero"
		pluginVersion = "v2.1.0"
	)
	archive := corruptZip(t, fmt.Sprintf("plugin-%s-%s-%s-%s", pluginName, pluginVersion, runtime.GOOS, runtime.GOARCH))

	mux := http.NewServeMux()
	mux.HandleFunc("/asset.zip", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archive)
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// An empty checksum is only warned about, so the archive reaches
		// extraction unverified - the realistic way a bad zip gets that far.
		_ = json.NewEncoder(w).Encode(cloudquery_api.PluginAsset{
			Location: server.URL + "/asset.zip",
		})
	})

	apiClient, err := cloudquery_api.NewClientWithResponses(server.URL)
	require.NoError(t, err)

	localPath := filepath.Join(t.TempDir(), "plugin")
	err = doDownloadPluginFromHub(context.Background(), zerolog.Nop(), apiClient, HubDownloadOptions{
		LocalPath:     localPath,
		PluginTeam:    "cloudquery",
		PluginKind:    PluginSource.String(),
		PluginName:    pluginName,
		PluginVersion: pluginVersion,
	}, DownloaderOptions{NoProgress: true})

	require.Error(t, err)
	_, statErr := os.Stat(localPath)
	require.ErrorIs(t, statErr, os.ErrNotExist, "a partial binary would be served as a cached plugin forever")
}
