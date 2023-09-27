package managedplugin

import (
	"context"
	"path"
	"testing"
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
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := DownloadPluginFromGithub(context.Background(), path.Join(tmp, tc.name), tc.org, tc.plugin, tc.version, tc.typ)
			if (err != nil) != tc.wantErr {
				t.Errorf("DownloadPluginFromGithub() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestDownloadPluginFromCloudQueryIntegration(t *testing.T) {
	tmp := t.TempDir()
	cases := []struct {
		testName string
		team     string
		plugin   string
		version  string
		wantErr  bool
		typ      PluginType
	}{
		{testName: "should download test plugin from cloudquery registry", team: "cloudquery", plugin: "test", version: "v3.1.11", typ: PluginSource},
	}
	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			err := DownloadPluginFromHub(context.Background(), path.Join(tmp, tc.testName), tc.team, tc.plugin, tc.version, tc.typ)
			if (err != nil) != tc.wantErr {
				t.Errorf("TestDownloadPluginFromCloudQueryIntegration() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
