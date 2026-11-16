package managedplugin

import "testing"

func TestIsCloudQueryRegistryImage(t *testing.T) {
	tests := []struct {
		imageName string
		want      bool
	}{
		{imageName: "docker.cloudquery.io/org/plugin:v1", want: true},
		{imageName: "DOCKER.CLOUDQUERY.IO/org/plugin:v1", want: false},
		{imageName: "docker.cloudquery.io.evil.com/x:v1", want: false},
		{imageName: "docker.cloudquery.iox/x", want: false},
		{imageName: "evil.com/docker.cloudquery.io/x", want: false},
		{imageName: "docker.cloudquery.io:8443/org/plugin:v1", want: false},
		{imageName: "cloudquery/plugin:v1", want: false},
		{imageName: "plugin:v1", want: false},
		{imageName: "docker.cloudquery.io/INVALID:v1", want: false},
	}
	for _, tc := range tests {
		t.Run(tc.imageName, func(t *testing.T) {
			if got := isCloudQueryRegistryImage(tc.imageName); got != tc.want {
				t.Fatalf("isCloudQueryRegistryImage(%q) = %v, want %v", tc.imageName, got, tc.want)
			}
		})
	}
}
