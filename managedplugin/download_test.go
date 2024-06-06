package managedplugin

import (
	"context"
	"path"
	"testing"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/rs/zerolog"
)

func TestDownloadPluginFromGithubIntegration(t *testing.T) {
	tmp := t.TempDir()
	cases := []struct {
		name    string
		org     string
		plugin  string
		version string
		wantErr bool
		typ     PluginType
	}{
		{name: "monorepo source", org: "cloudquery", plugin: "hackernews", version: "v1.1.4", typ: PluginSource},
		{name: "many repo source", org: "cloudquery", plugin: "simple-analytics", version: "v1.0.0", typ: PluginSource},
		{name: "monorepo destination", org: "cloudquery", plugin: "postgresql", version: "v2.0.7", typ: PluginDestination},
		{name: "community source", org: "hermanschaaf", plugin: "simple-analytics", version: "v1.0.0", typ: PluginSource},
		{name: "invalid community source", org: "cloudquery", plugin: "invalid-plugin", version: "v0.0.x", wantErr: true, typ: PluginSource},
		{name: "invalid monorepo source", org: "not-cloudquery", plugin: "invalid-plugin", version: "v0.0.x", wantErr: true, typ: PluginSource},
	}
	logger := zerolog.Logger{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := DownloadPluginFromGithub(context.Background(), logger, path.Join(tmp, tc.name), tc.org, tc.plugin, tc.version, tc.typ, DownloaderOptions{})
			if (err != nil) != tc.wantErr {
				t.Errorf("DownloadPluginFromGithub() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestDownloadPluginFromCloudQueryHub(t *testing.T) {
	tmp := t.TempDir()
	cases := []struct {
		testName string
		team     string
		plugin   string
		version  string
		wantErr  bool
		typ      PluginType
	}{
		{testName: "should download test plugin from cloudquery registry", team: "cloudquery", plugin: "aws", version: "v22.18.0", typ: PluginSource},
	}
	c, err := cloudquery_api.NewClientWithResponses(APIBaseURL())
	if err != nil {
		t.Fatalf("failed to create Hub API client: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			err := DownloadPluginFromHub(context.Background(), c, HubDownloadOptions{
				LocalPath:     path.Join(tmp, tc.testName),
				AuthToken:     "",
				TeamName:      "",
				PluginTeam:    tc.team,
				PluginKind:    tc.typ.String(),
				PluginName:    tc.plugin,
				PluginVersion: tc.version,
			},
				DownloaderOptions{},
			)
			if (err != nil) != tc.wantErr {
				t.Errorf("TestDownloadPluginFromCloudQueryIntegration() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
