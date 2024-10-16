package managedplugin

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginVersionWarnerUnknownPluginFails(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop())
	require.NoError(t, err)
	warned, err := versionWarner.WarnIfOutdated(context.Background(), "unknown", "unknown", "source", "1.0.0")
	assert.Error(t, err)
	assert.False(t, warned)
}

// Note: this is an integration test that requires Internet access and the hub to be running
func TestPluginLatestVersionDoesNotWarn(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop())
	require.NoError(t, err)
	latestVersion, err := versionWarner.getLatestVersion(context.Background(), "cloudquery", "aws", "source")
	assert.NoError(t, err)
	hasWarned, err := versionWarner.WarnIfOutdated(context.Background(), "cloudquery", "aws", "source", latestVersion.String())
	assert.NoError(t, err)
	assert.False(t, hasWarned)
}

// Note: this is an integration test that requires Internet access and the hub to be running
// CloudQuery's aws source plugin must exist in the hub, and be over version v1.0.0
func TestPluginLatestVersionWarns(t *testing.T) {
	versionWarner, err := NewPluginVersionWarner(zerolog.Nop())
	require.NoError(t, err)
	hasWarned, err := versionWarner.WarnIfOutdated(context.Background(), "cloudquery", "aws", "source", "v1.0.0")
	assert.NoError(t, err)
	assert.True(t, hasWarned)
}
