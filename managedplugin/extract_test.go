package managedplugin

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTestArchive(t *testing.T, dir, entry string, contents []byte) string {
	t.Helper()

	archivePath := filepath.Join(dir, "plugin.zip")
	f, err := os.Create(archivePath)
	require.NoError(t, err)

	w := zip.NewWriter(f)
	entryWriter, err := w.Create(entry)
	require.NoError(t, err)
	_, err = entryWriter.Write(contents)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	require.NoError(t, f.Close())

	return archivePath
}

func TestExtractPluginBinary(t *testing.T) {
	dir := t.TempDir()
	binary := []byte("plugin-binary")
	archivePath := writeTestArchive(t, dir, "plugin-aws-v1.0.0-linux-amd64", binary)
	localPath := filepath.Join(dir, "aws")

	require.NoError(t, extractPluginBinary(archivePath, "plugin-aws-v1.0.0-linux-amd64", localPath))

	got, err := os.ReadFile(localPath)
	require.NoError(t, err)
	require.Equal(t, binary, got)

	info, err := os.Stat(localPath)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0744), info.Mode().Perm())
}

// TestExtractPluginBinaryLeavesNoPartialFile guards the caching path: a failed
// extraction that left bytes at localPath would make the next run skip the
// download and exec a truncated binary.
func TestExtractPluginBinaryLeavesNoPartialFile(t *testing.T) {
	dir := t.TempDir()
	archivePath := writeTestArchive(t, dir, "plugin-aws-v1.0.0-linux-amd64", []byte("plugin-binary"))
	localPath := filepath.Join(dir, "aws")

	err := extractPluginBinary(archivePath, "plugin-aws-v9.9.9-linux-amd64", localPath)
	require.Error(t, err)

	_, statErr := os.Stat(localPath)
	require.ErrorIs(t, statErr, os.ErrNotExist)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, entry := range entries {
		require.NotContains(t, entry.Name(), ".tmp", "the temporary file must not survive a failed extraction")
	}
}
