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
		cloudquery_api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
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

func isDockerPlugin(ctx context.Context, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions) (bool, error) {
	p, err := c.GetPluginVersionWithResponse(ctx, ops.PluginTeam, cloudquery_api.PluginKind(ops.PluginKind), ops.PluginName, ops.PluginVersion)
	if err != nil {
		return false, fmt.Errorf("failed to get plugin information: %w", err)
	}
	if p.StatusCode() != http.StatusOK {
		return false, fmt.Errorf("failed to get plugin information: %s", p.Status())
	}
	if p.JSON200 == nil {
		return false, fmt.Errorf("failed to get plugin information: response body is empty")
	}
	return p.JSON200.PackageType == "docker", nil
}
