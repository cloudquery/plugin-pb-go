package managedplugin

import (
	"context"
	"fmt"
	"net/http"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/rs/zerolog"
)

func getHubClient(logger zerolog.Logger, ops HubDownloadOptions) (*cloudquery_api.ClientWithResponses, error) {
	c, err := cloudquery_api.NewClientWithResponses(APIBaseURL(),
		cloudquery_api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			if ops.AuthToken != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ops.AuthToken))
			}
			logger.Debug().Interface("ops", ops).Msg(fmt.Sprintf("Requesting %s %s", req.Method, req.URL))
			return nil
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to create Hub API client: %w", err)
	}
	return c, nil
}

// validateDockerPlugin checks if the plugin has PackageType=docker and ops.PluginVersion exists
func validateDockerPlugin(ctx context.Context, logger zerolog.Logger, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions) (bool, error) {
	var errFailed = fmt.Sprintf("failed to get %s plugin (name: %s/%s@%s) information", cloudquery_api.PluginKind(ops.PluginKind), ops.PluginTeam, ops.PluginName, ops.PluginVersion)

	p, err := c.GetPluginVersionWithResponse(ctx, ops.PluginTeam, cloudquery_api.PluginKind(ops.PluginKind), ops.PluginName, ops.PluginVersion)
	if err != nil {
		return false, fmt.Errorf(errFailed+": %w", err)
	}
	if p.StatusCode() == http.StatusNotFound {
		// See if the plugin exists, but not the version.
		pvw, err := NewPluginVersionWarner(logger, ops.AuthToken)
		if err != nil {
			return false, fmt.Errorf(errFailed+": %w", err)
		}

		ver, err := pvw.getLatestVersion(ctx, ops.PluginTeam, ops.PluginName, ops.PluginKind)
		if err != nil {
			return false, fmt.Errorf(errFailed+": %w", err)
		}

		if ver != nil {
			return false, fmt.Errorf("version %s does not exist, consider using the latest version at %s", ops.PluginVersion,
				fmt.Sprintf("https://hub.cloudquery.io/plugins/%s/%s/%s/v%s", ops.PluginKind, ops.PluginTeam, ops.PluginName, ver.String()))
		}
	}
	if p.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("failed to get %s plugin (name: %s/%s@%s) information: %s", cloudquery_api.PluginKind(ops.PluginKind), ops.PluginTeam, ops.PluginName, ops.PluginVersion, p.Status())
	}
	if p.JSON200 == nil {
		return false, fmt.Errorf("failed to get %s plugin (name: %s/%s@%s) information: response body is empty", cloudquery_api.PluginKind(ops.PluginKind), ops.PluginTeam, ops.PluginName, ops.PluginVersion)
	}
	return p.JSON200.PackageType == "docker", nil
}
