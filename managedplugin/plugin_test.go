package managedplugin

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
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
	clients, err := NewClients(ctx, PluginSource, []Config{cfg},
		WithDirectory(tmpDir), WithNoSentry())
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

func TestManagedPluginCloudQuery(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	cfg := Config{
		Name:     "awspricing",
		Registry: RegistryCloudQuery,
		Path:     "cloudquery/awspricing",
		Version:  "v3.0.12",
	}
	clients, err := NewClients(ctx, PluginSource, []Config{cfg}, WithDirectory(tmpDir), WithNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	testClient := clients.ClientByName("awspricing")
	if testClient == nil {
		t.Fatal("awspricing client not found")
	}
	if err := clients.Terminate(); err != nil {
		t.Fatal(err)
	}
	localPath := filepath.Join(tmpDir, "plugins", PluginSource.String(), "cloudquery", "awspricing", cfg.Version, "plugin")
	localPath = WithBinarySuffix(localPath)
	cfg = Config{
		Name:     "awspricing",
		Registry: RegistryLocal,
		Path:     localPath,
		Version:  cfg.Version,
	}

	clients, err = NewClients(ctx, PluginSource, []Config{cfg}, WithDirectory(tmpDir), WithNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	testClient = clients.ClientByName("awspricing")
	if testClient == nil {
		t.Fatal("awspricing client not found")
	}
	if err := clients.Terminate(); err != nil {
		t.Fatal(err)
	}
}

func TestManagedPluginCloudQueryDocker(t *testing.T) {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Ping(ctx)
	if err != nil {
		t.Skip("docker not running")
	}
	// note: this test requires an API key and team name to be set in the environment
	if os.Getenv("CLOUDQUERY_TEAM_NAME") == "" {
		t.Skip("CLOUDQUERY_TEAM_NAME not set")
	}
	if os.Getenv("CLOUDQUERY_API_KEY") == "" {
		t.Skip("CLOUDQUERY_API_KEY not set")
	}
	if runtime.GOOS == "windows" {
		// the docker image is not built for Windows, so would require enabling of experimental
		// linux compatibility. We skip this test in CI for now.
		t.Skip("this test is not supported on windows")
	}

	tmpDir := t.TempDir()
	cfg := Config{
		Name:     "test",
		Registry: RegistryCloudQuery,
		Version:  "v1.2.3",
		Path:     "cloudquery/typeform",
	}
	clients, err := NewClients(ctx, PluginSource, []Config{cfg},
		WithDirectory(tmpDir), WithNoSentry(),
		WithAuthToken(os.Getenv("CLOUDQUERY_API_KEY")),
		WithTeamName(os.Getenv("CLOUDQUERY_TEAM_NAME")),
	)
	if err != nil {
		t.Fatal(err)
	}
	testClient := clients.ClientByName("test")
	if testClient == nil {
		t.Fatal("test client not found")
	}
	v, err := testClient.Versions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) < 1 {
		t.Fatal("expected at least 1 version, got 0")
	}
	if err := clients.Terminate(); err != nil {
		t.Fatal(err)
	}
}

func TestManagedPluginDocker(t *testing.T) {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Ping(ctx)
	if err != nil {
		t.Skip("docker not running")
	}
	if runtime.GOOS == "windows" {
		// the docker image is not built for Windows, so would require enabling of experimental
		// linux compatibility. We skip this test in CI for now.
		t.Skip("this test is not supported on windows")
	}

	tmpDir := t.TempDir()
	cfg := Config{
		Name:     "test",
		Registry: RegistryDocker,
		Path:     "ghcr.io/cloudquery/cq-source-test:3.0.3",
	}
	clients, err := NewClients(ctx, PluginSource, []Config{cfg}, WithDirectory(tmpDir), WithNoSentry())
	if err != nil {
		t.Fatal(err)
	}
	testClient := clients.ClientByName("test")
	if testClient == nil {
		t.Fatal("test client not found")
	}
	v, err := testClient.Versions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) < 1 {
		t.Fatal("expected at least 1 version, got 0")
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

func TestValidateLocalExecPath(t *testing.T) {
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
