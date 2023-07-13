package managedplugin

import (
	"context"
	"os"
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

func TestIsDirectory(t *testing.T) {
	tempDir := t.TempDir()
	directoryBool, err := isDirectory(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	// directory should be `true`
	if !directoryBool {
		t.Fatal("expected directory")
	}
	tempFile, err := os.Create(tempDir + "testFile")
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()
	isFileBool, err := isDirectory(tempDir + "testFile")
	if err != nil {
		t.Fatal(err)
	}
	if isFileBool {
		t.Fatal("expected file")
	}
}

func TestValidatevLocalExecPath(t *testing.T) {
	tempDir := t.TempDir()
	// passing a directory should result in an error
	err := validateLocalExecPath(tempDir)
	if err == nil {
		t.Fatal(err)
	}

	tempFile, err := os.Create(tempDir + "testFile")
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()
	err = validateLocalExecPath(tempDir + "testFile")
	// passing a valid non directory path should result in no error
	if err != nil {
		t.Fatal(err)
	}
	// passing a file path that doesn't exist should result in an error
	err = validateLocalExecPath(tempDir + "randomfile")
	if err == nil {
		t.Fatal(err)
	}
}
