package managedplugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginVersionWarnerUnknownPluginFails(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)
	latestVersion, err := versionWarner.LatestVersion(context.Background(), "unknown", "unknown", "source")
	assert.Error(t, err)
	assert.Nil(t, latestVersion)
}

func TestPluginVersionWarnerLatestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plugins/cloudquery/source/aws", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"latest_version":"v32.1.0"}`))
	}))
	t.Cleanup(server.Close)
	t.Setenv("CLOUDQUERY_API_URL", server.URL)

	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)

	latestVersion, err := versionWarner.LatestVersion(context.Background(), "cloudquery", "aws", "source")
	require.NoError(t, err)
	assert.Equal(t, "32.1.0", latestVersion.String())
}

func TestPluginVersionWarnerInvalidKindFails(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)
	_, err = versionWarner.LatestVersion(context.Background(), "cloudquery", "aws", "invalid")
	assert.Error(t, err)
}
