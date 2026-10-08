package managedplugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginVersionWarnerUnknownPluginFails(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)
	warned, err := versionWarner.WarnIfOutdated(context.Background(), "unknown", "unknown", "source", "1.0.0")
	assert.Error(t, err)
	assert.False(t, warned)
}

func TestPluginVersionWarnerInvalidOrgOrNameFails(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)
	for _, tc := range []struct{ org, name string }{
		{org: "", name: ".cq"},
		{org: "cloudquery", name: ""},
		{org: "", name: ""},
		{org: "cloudquery/aws", name: "plugin"},
	} {
		warned, err := versionWarner.WarnIfOutdated(context.Background(), tc.org, tc.name, "source", "1.0.0")
		assert.Error(t, err, "%q/%q", tc.org, tc.name)
		assert.False(t, warned)
	}
}

// Note: this is an integration test that requires Internet access and the hub to be running
func TestPluginLatestVersionDoesNotWarn(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)
	latestVersion, err := versionWarner.LatestVersion(context.Background(), "cloudquery", "aws", "source")
	assert.NoError(t, err)
	hasWarned, err := versionWarner.WarnIfOutdated(context.Background(), "cloudquery", "aws", "source", latestVersion.String())
	assert.NoError(t, err)
	assert.False(t, hasWarned)
}

// Note: this is an integration test that requires Internet access and the hub to be running
// CloudQuery's aws source plugin must exist in the hub, and be over version v1.0.0
func TestPluginLatestVersionWarns(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)
	hasWarned, err := versionWarner.WarnIfOutdated(context.Background(), "cloudquery", "aws", "source", "v1.0.0")
	assert.NoError(t, err)
	assert.True(t, hasWarned)
}

func TestPluginVersionWarnerRequestsLatestVersionOnce(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		assert.Equal(t, "/plugins/cloudquery/source/aws", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"latest_version":"v32.1.0"}`))
	}))
	t.Cleanup(server.Close)
	t.Setenv("CLOUDQUERY_API_URL", server.URL)

	versionWarner, err := NewPluginVersionWarner(zerolog.Nop(), "")
	require.NoError(t, err)

	warned, err := versionWarner.WarnIfOutdated(context.Background(), "cloudquery", "aws", "source", "v31.0.0")
	require.NoError(t, err)
	assert.True(t, warned)

	latestVersion, err := versionWarner.LatestVersion(context.Background(), "cloudquery", "aws", "source")
	require.NoError(t, err)
	assert.Equal(t, "32.1.0", latestVersion.String())
	assert.Equal(t, int32(1), requests.Load())
}
