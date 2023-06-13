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
	}{
		{name: "monorepo source", org: "cloudquery", plugin: "hackernews", version: "v1.1.4"},
		{name: "many repo source", org: "cloudquery", plugin: "simple-analytics", version: "v1.0.0"},
		{name: "monorepo destination", org: "cloudquery", plugin: "postgresql", version: "v2.0.7"},
		{name: "community source", org: "hermanschaaf", plugin: "simple-analytics", version: "v1.0.0"},
		{name: "invalid community source", org: "cloudquery", plugin: "invalid-plugin", version: "v0.0.x", wantErr: true},
		{name: "invalid monorepo source", org: "not-cloudquery", plugin: "invalid-plugin", version: "v0.0.x", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := DownloadPluginFromGithub(context.Background(), path.Join(tmp, tc.name), tc.org, tc.plugin, tc.version)
			if (err != nil) != tc.wantErr {
				t.Errorf("DownloadPluginFromGithub() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
