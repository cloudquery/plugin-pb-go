package managedplugin

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/rs/zerolog"
)

type PluginVersionWarner struct {
	hubClient *cloudquery_api.ClientWithResponses
	logger    zerolog.Logger
}

func NewPluginVersionWarner(logger zerolog.Logger) (*PluginVersionWarner, error) {
	hubClient, err := getHubClient(logger, HubDownloadOptions{}) // Does not use auth token, since API call is public
	if err != nil {
		return nil, err
	}
	return &PluginVersionWarner{hubClient: hubClient, logger: logger}, nil
}

func (p *PluginVersionWarner) getLatestVersion(ctx context.Context, org string, name string, kind string) (*semver.Version, error) {
	if p == nil {
		return nil, fmt.Errorf("plugin version warner is not initialized")
	}
	if kind != PluginSource.String() && kind != PluginDestination.String() && kind != PluginTransformer.String() {
		p.logger.Debug().Str("plugin", name).Str("kind", kind).Msg("invalid kind")
		return nil, fmt.Errorf("invalid kind: %s", kind)
	}
	resp, err := p.hubClient.GetPluginWithResponse(ctx, org, cloudquery_api.PluginKind(kind), name)
	if err != nil {
		p.logger.Debug().Str("plugin", name).Err(err).Msg("failed to get plugin info from hub")
		return nil, err
	}
	if resp.JSON200 == nil {
		p.logger.Debug().Str("plugin", name).Msg("failed to get plugin info from hub, request didn't error but 200 response is nil")
		return nil, fmt.Errorf("failed to get plugin info from hub, request didn't error but 200 response is nil")
	}
	if resp.JSON200.LatestVersion == nil {
		p.logger.Debug().Str("plugin", name).Msg("cannot check if plugin is outdated, latest version is nil")
		return nil, fmt.Errorf("cannot check if plugin is outdated, latest version is nil")
	}
	latestVersion := *resp.JSON200.LatestVersion
	latestSemver, err := semver.NewVersion(latestVersion)
	if err != nil {
		p.logger.Debug().Str("plugin", name).Str("version", latestVersion).Err(err).Msg("failed to parse latest version")
		return nil, err
	}
	return latestSemver, nil
}

// WarnIfOutdated requests the latest version of a plugin from the hub and warns if the client's supplied version is outdated.
// It returns true if nothing went wrong comparing the versions, and the client's version is outdated; false otherwise.
func (p *PluginVersionWarner) WarnIfOutdated(ctx context.Context, org string, name string, kind string, actualVersion string) (bool, error) {
	if p == nil {
		return false, fmt.Errorf("plugin version warner is not initialized")
	}
	if actualVersion == "" {
		return false, nil
	}
	actualVersionSemver, err := semver.NewVersion(actualVersion)
	if err != nil {
		p.logger.Debug().Str("plugin", name).Str("version", actualVersion).Err(err).Msg("failed to parse actual version")
		return false, err
	}
	latestVersionSemver, err := p.getLatestVersion(ctx, org, name, kind)
	if err != nil {
		return false, err
	}
	if actualVersionSemver.LessThan(latestVersionSemver) {
		p.logger.Warn().
			Str("plugin", name).
			Str("using_version", actualVersionSemver.String()).
			Str("latest_version", latestVersionSemver.String()).
			Str("url", fmt.Sprintf("https://hub.cloudquery.io/plugins/%s/%s/%s", kind, org, name)).
			Msg("Plugin is outdated, consider upgrading to the latest version.")
		return true, nil
	}

	return false, nil
}
