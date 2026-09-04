package managedplugin

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func otherArch() string {
	if runtime.GOARCH == "arm64" {
		return "amd64"
	}
	return "arm64"
}

// buildFixture compiles a trivial program for the given architecture so tests
// have a real binary to inspect rather than a hand crafted header.
func buildFixture(t *testing.T, goarch string) string {
	t.Helper()

	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not available")
	}

	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte("module fixture\n\ngo 1.26\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(t.TempDir(), WithBinarySuffix("fixture-"+goarch))
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(), "GOOS="+runtime.GOOS, "GOARCH="+goarch, "CGO_ENABLED=0")
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build fixture for %s: %v\n%s", goarch, err, combined)
	}
	return out
}

func TestValidateBinaryAcceptsHostArch(t *testing.T) {
	if err := validateBinary(buildFixture(t, runtime.GOARCH)); err != nil {
		t.Fatalf("expected host binary to validate, got %v", err)
	}
}

func TestValidateBinaryRejectsOtherArch(t *testing.T) {
	err := validateBinary(buildFixture(t, otherArch()))
	if !errors.Is(err, ErrArchMismatch) {
		t.Fatalf("expected ErrArchMismatch, got %v", err)
	}
}

func TestValidateBinaryRejectsTruncated(t *testing.T) {
	full, err := os.ReadFile(buildFixture(t, runtime.GOARCH))
	if err != nil {
		t.Fatal(err)
	}

	for name, body := range map[string][]byte{
		"empty":     {},
		"truncated": full[:512],
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "plugin")
			if err := os.WriteFile(path, body, 0744); err != nil {
				t.Fatal(err)
			}
			if err := validateBinary(path); err == nil {
				t.Fatal("expected a partially written binary to be rejected")
			} else if os.IsNotExist(err) {
				t.Fatalf("expected a validation error rather than not-exist, got %v", err)
			}
		})
	}
}

func TestValidateBinaryReportsMissing(t *testing.T) {
	err := validateBinary(filepath.Join(t.TempDir(), "absent"))
	if !os.IsNotExist(err) {
		t.Fatalf("expected a not-exist error, got %v", err)
	}
}

func TestPluginCachePathsNamespacesByTarget(t *testing.T) {
	canonical, legacy := pluginCachePaths(".cq", "destination", "cloudquery", "databricks", "v1.6.14")

	target := runtime.GOOS + "_" + runtime.GOARCH
	if !strings.Contains(canonical, target) {
		t.Errorf("expected canonical path %q to be namespaced by %q", canonical, target)
	}
	if strings.Contains(legacy, target) {
		t.Errorf("expected legacy path %q to keep its original layout", legacy)
	}
	if filepath.Dir(canonical) != filepath.Join(filepath.Dir(legacy), target) {
		t.Errorf("expected canonical path to sit under the legacy directory, got %q and %q", canonical, legacy)
	}
}

// TestResolveCachedPluginPrefersUsableBinary covers the reason a wrong
// architecture binary reached exec in the first place: the cache was keyed on
// existence alone, so a host of one architecture happily ran what a host of
// another had left behind.
func TestResolveCachedPluginPrefersUsableBinary(t *testing.T) {
	hostArch := buildFixture(t, runtime.GOARCH)
	wrongArch := buildFixture(t, otherArch())

	copyTo := func(t *testing.T, src, dest string) {
		t.Helper()
		body, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dest, body, 0744); err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range []struct {
		name         string
		canonical    string
		legacy       string
		wantResolved func(canonical, legacy string) string
		wantCached   bool
	}{
		{
			name:         "canonical hit",
			canonical:    hostArch,
			wantResolved: func(canonical, _ string) string { return canonical },
			wantCached:   true,
		},
		{
			name:         "legacy hit with matching arch is reused",
			legacy:       hostArch,
			wantResolved: func(_, legacy string) string { return legacy },
			wantCached:   true,
		},
		{
			name:         "legacy hit with other arch downloads to canonical",
			legacy:       wrongArch,
			wantResolved: func(canonical, _ string) string { return canonical },
			wantCached:   false,
		},
		{
			name:         "no cache downloads to canonical",
			wantResolved: func(canonical, _ string) string { return canonical },
			wantCached:   false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			canonical, legacy := pluginCachePaths(t.TempDir(), "source", "cloudquery", "crowdstrike", "v2.3.9")
			if tc.canonical != "" {
				copyTo(t, tc.canonical, canonical)
			}
			if tc.legacy != "" {
				copyTo(t, tc.legacy, legacy)
			}

			resolved, cached := resolveCachedPlugin(zerolog.Nop(), canonical, legacy)
			if cached != tc.wantCached {
				t.Errorf("got cached=%v, want %v", cached, tc.wantCached)
			}
			if want := tc.wantResolved(canonical, legacy); resolved != want {
				t.Errorf("got resolved=%q, want %q", resolved, want)
			}
		})
	}
}

// TestResolveCachedPluginIsolatesArchitectures asserts that hosts of different
// architectures do not contend for one path, which is what let a shared cq
// directory on Databricks Serverless flip between working and failing.
func TestResolveCachedPluginIsolatesArchitectures(t *testing.T) {
	dir := t.TempDir()
	canonical, legacy := pluginCachePaths(dir, "source", "cloudquery", "crowdstrike", "v2.3.9")

	otherCanonical := WithBinarySuffix(filepath.Join(
		dir, "plugins", "source", "cloudquery", "crowdstrike", "v2.3.9",
		runtime.GOOS+"_"+otherArch(), "plugin",
	))
	if filepath.Dir(canonical) == filepath.Dir(otherCanonical) {
		t.Fatal("expected each architecture to get its own cache directory")
	}

	body, err := os.ReadFile(buildFixture(t, otherArch()))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(otherCanonical), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(otherCanonical, body, 0744); err != nil {
		t.Fatal(err)
	}

	if _, cached := resolveCachedPlugin(zerolog.Nop(), canonical, legacy); cached {
		t.Fatal("expected another architecture's cached binary to be ignored")
	}
	if _, err := os.Stat(otherCanonical); err != nil {
		t.Fatalf("expected another architecture's cache to be left intact, got %v", err)
	}
}
