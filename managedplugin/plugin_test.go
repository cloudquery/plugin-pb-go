package managedplugin

import (
	"context"
	"path/filepath"
	"testing"
)

func TestManagedPluginGitHub(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	cfg := Config{
		Name:     "hackernews",
		Registry: RegistryGithub,
		Path:     "cloudquery/hackernews",
		Version:  "v1.1.4",
	}
	clients, err := NewClients(ctx, PluginSource, []Config{cfg}, WithDirectory(tmpDir), WithNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	hnClient := clients.ClientByName("hackernews")
	if hnClient == nil {
		t.Fatal("hackernews client not found")
	}
	if err := clients.Terminate(); err != nil {
		t.Fatal(err)
	}
	localPath := filepath.Join(tmpDir, "plugins", PluginSource.String(), "cloudquery", "hackernews", cfg.Version, "plugin")
	localPath = WithBinarySuffix(localPath)
	cfg = Config{
		Name:     "hackernews",
		Registry: RegistryLocal,
		Path:     localPath,
		Version:  "v1.1.4",
	}

	clients, err = NewClients(ctx, PluginSource, []Config{cfg}, WithDirectory(tmpDir), WithNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	hnClient = clients.ClientByName("hackernews")
	if hnClient == nil {
		t.Fatal("hackernews client not found")
	}
	if err := clients.Terminate(); err != nil {
		t.Fatal(err)
	}
}
