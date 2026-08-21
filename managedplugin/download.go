package managedplugin

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/avast/retry-go/v5"
	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"
)

const (
	DefaultDownloadDir = ".cq"
	RetryAttempts      = 5
	RetryWaitTime      = 1 * time.Second
	MaxRetryWaitTime   = 8 * time.Second
)

func APIBaseURL() string {
	const (
		envAPIURL  = "CLOUDQUERY_API_URL"
		apiBaseURL = "https://api.cloudquery.io"
	)

	if v := os.Getenv(envAPIURL); v != "" {
		return v
	}
	return apiBaseURL
}

// getURLLocation return the URL of the plugin
// this does a few HEAD requests because we had a few breaking changes to where
// we store the plugins on GitHub
// TODO: we can improve this by just embedding all plugins and last version that exist in different places then
// the latest
func getURLLocation(ctx context.Context, org string, name string, version string, typ PluginType) (string, error) {
	urls := []string{
		// TODO: add this back when we move to the new plugin system
		// fmt.Sprintf("https://github.com/%s/cq-plugin-%s/releases/download/%s/cq-%s_%s_%s.zip", org, name, version, name, runtime.GOOS, runtime.GOARCH),
		fmt.Sprintf("https://github.com/%s/cq-source-%s/releases/download/%s/cq-source-%s_%s_%s.zip", org, name, version, name, runtime.GOOS, runtime.GOARCH),
	}
	if org == "cloudquery" {
		// TODO: add this back when we move to the new plugin system
		// urls = append(urls, fmt.Sprintf("https://github.com/cloudquery/cloudquery/releases/download/plugins-%s-%s/%s_%s_%s.zip", name, version, name, runtime.GOOS, runtime.GOARCH))
		urls = append(urls, fmt.Sprintf("https://github.com/cloudquery/cloudquery/releases/download/plugins-source-%s-%s/%s_%s_%s.zip", name, version, name, runtime.GOOS, runtime.GOARCH))
	}
	if typ == PluginDestination {
		urls = []string{
			// TODO: add this back when we move to the new plugin system
			// fmt.Sprintf("https://github.com/%s/cq-plugin-%s/releases/download/%s/cq-%s_%s_%s.zip", org, name, version, name, runtime.GOOS, runtime.GOARCH),
			fmt.Sprintf("https://github.com/%s/cq-destination-%s/releases/download/%s/cq-destination-%s_%s_%s.zip", org, name, version, name, runtime.GOOS, runtime.GOARCH),
		}
		if org == "cloudquery" {
			// TODO: add this back when we move to the new plugin system
			// urls = append(urls, fmt.Sprintf("https://github.com/cloudquery/cloudquery/releases/download/plugins-%s-%s/%s_%s_%s.zip", name, version, name, runtime.GOOS, runtime.GOARCH))
			urls = append(urls, fmt.Sprintf("https://github.com/cloudquery/cloudquery/releases/download/plugins-destination-%s-%s/%s_%s_%s.zip", name, version, name, runtime.GOOS, runtime.GOARCH))
		}
	}

	var (
		err404 = errors.New("404")
		err401 = errors.New("401")
		err429 = errors.New("429")
	)

	options := []retry.Option{
		retry.RetryIf(func(err error) bool {
			return err == err401 || err == err429
		}),
		retry.Context(ctx),
		retry.Attempts(RetryAttempts),
		retry.Delay(RetryWaitTime),
		retry.LastErrorOnly(true),
	}
	retrier := retry.New(options...)
	for _, downloadURL := range urls {
		err := retrier.Do(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
			if err != nil {
				return fmt.Errorf("failed create request %s: %w", downloadURL, err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to get url %s: %w", downloadURL, err)
			}
			resp.Body.Close()
			// Check server response
			switch resp.StatusCode {
			case http.StatusOK:
				return nil
			case http.StatusNotFound:
				return err404
			case http.StatusUnauthorized:
				fmt.Printf("Failed downloading %s with status code %d. Retrying\n", downloadURL, resp.StatusCode)
				return err401
			case http.StatusTooManyRequests:
				fmt.Printf("Failed downloading %s with status code %d. Retrying\n", downloadURL, resp.StatusCode)
				return err429
			default:
				fmt.Printf("Failed downloading %s with status code %d\n", downloadURL, resp.StatusCode)
				return fmt.Errorf("statusCode %d", resp.StatusCode)
			}
		})
		if err == err404 {
			continue
		}
		return downloadURL, err
	}

	return "", fmt.Errorf("failed to find plugin %s/%s version %s", org, name, version)
}

type HubDownloadOptions struct {
	AuthToken     string
	TeamName      string
	LocalPath     string
	PluginTeam    string
	PluginKind    string
	PluginName    string
	PluginVersion string
}
type DownloaderOptions struct {
	NoProgress bool
}

func DownloadPluginFromHub(ctx context.Context, logger zerolog.Logger, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions, dops DownloaderOptions) (AssetSource, error) {
	if _, err := os.Stat(ops.LocalPath); err == nil {
		return AssetSourceCached, nil
	}
	return AssetSourceRemote, doDownloadPluginFromHub(ctx, logger, c, ops, dops)
}

func doDownloadPluginFromHub(ctx context.Context, logger zerolog.Logger, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions, dops DownloaderOptions) error {
	downloadDir := filepath.Dir(ops.LocalPath)
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	pluginAsset, statusCode, err := downloadPluginAssetFromHub(ctx, c, ops)
	if err != nil {
		return fmt.Errorf("failed to get plugin metadata from hub: %w", err)
	}

	switch statusCode {
	case http.StatusOK:
		// we allow this status code
	case http.StatusUnauthorized:
		return errors.New("unauthorized. Try logging in via `cloudquery login`")
	case http.StatusNotFound:
		var errRetryWithLogin = fmt.Errorf("failed to download plugin %v %v/%v@%v: plugin version not found. If you're trying to use a private plugin you'll need to run `cloudquery login` first", ops.PluginKind, ops.PluginTeam, ops.PluginName, ops.PluginVersion)

		// See if the plugin exists, but not the version.
		pvw, err := NewPluginVersionWarner(logger, ops.AuthToken)
		if err != nil {
			return errRetryWithLogin
		}

		ver, err := pvw.getLatestVersion(ctx, ops.PluginTeam, ops.PluginName, ops.PluginKind)
		if err != nil {
			return errRetryWithLogin
		}

		if ver != nil {
			return fmt.Errorf("version %s does not exist, consider using the latest version at %s", ops.PluginVersion,
				fmt.Sprintf("https://www.cloudquery.io/hub/plugins/%s/%s/%s/v%s", ops.PluginKind, ops.PluginTeam, ops.PluginName, ver.String()))
		}

		return errRetryWithLogin
	case http.StatusTooManyRequests:
		return errors.New("too many download requests. Try logging in via `cloudquery login` to increase rate limits")
	default:
		return fmt.Errorf("failed to download plugin %v %v/%v@%v: unexpected status code %v", ops.PluginKind, ops.PluginTeam, ops.PluginName, ops.PluginVersion, statusCode)
	}
	if pluginAsset == nil {
		return fmt.Errorf("failed to get plugin metadata from hub for %v %v/%v@%v: missing json response", ops.PluginKind, ops.PluginTeam, ops.PluginName, ops.PluginVersion)
	}
	location := pluginAsset.Location
	if len(location) == 0 {
		return errors.New("failed to get plugin metadata from hub: empty location from response")
	}
	pluginZipPath := ops.LocalPath + ".zip"
	defer os.Remove(pluginZipPath)

	writtenChecksum, err := downloadFile(ctx, pluginZipPath, location, dops)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	if pluginAsset.Checksum == "" {
		fmt.Printf("Warning - checksum not verified: %s\n", writtenChecksum)
	} else if writtenChecksum != pluginAsset.Checksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", pluginAsset.Checksum, writtenChecksum)
	}

	pathInArchive := fmt.Sprintf("plugin-%s-%s-%s-%s", ops.PluginName, ops.PluginVersion, runtime.GOOS, runtime.GOARCH)
	return extractPluginBinary(pluginZipPath, pathInArchive, ops.LocalPath)
}

// extractPluginBinary writes the binary to a temporary file and renames it into
// place, so a failure part way through never leaves a truncated binary that the
// next run treats as a cached plugin.
func extractPluginBinary(archivePath, pathInArchive, localPath string) error {
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

	fileInArchive, err := archive.Open(pathInArchive)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive %s: %w", pathInArchive, err)
	}
	defer fileInArchive.Close()

	out, err := os.CreateTemp(filepath.Dir(localPath), filepath.Base(localPath)+".tmp")
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	tmpPath := out.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(out, fileInArchive); err != nil {
		out.Close()
		return fmt.Errorf("failed to copy body to file: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	if err := os.Chmod(tmpPath, 0744); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", localPath, err)
	}
	if err := os.Rename(tmpPath, localPath); err != nil {
		return fmt.Errorf("failed to move plugin binary to %s: %w", localPath, err)
	}
	return nil
}

func downloadPluginAssetFromHub(ctx context.Context, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions) (*cloudquery_api.PluginAsset, int, error) {
	target := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	aj := "application/json"

	switch {
	case ops.TeamName != "":
		resp, err := c.DownloadPluginAssetByTeamWithResponse(
			ctx,
			ops.TeamName,
			ops.PluginTeam,
			cloudquery_api.PluginKind(ops.PluginKind),
			ops.PluginName,
			ops.PluginVersion,
			target,
			&cloudquery_api.DownloadPluginAssetByTeamParams{Accept: &aj},
		)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to request with team: %w", err)
		}
		return resp.JSON200, resp.StatusCode(), nil
	default:
		resp, err := c.DownloadPluginAssetWithResponse(
			ctx,
			ops.PluginTeam,
			cloudquery_api.PluginKind(ops.PluginKind),
			ops.PluginName,
			ops.PluginVersion,
			target,
			&cloudquery_api.DownloadPluginAssetParams{Accept: &aj},
		)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to request: %w", err)
		}
		return resp.JSON200, resp.StatusCode(), nil
	}
}

func DownloadPluginFromGithub(ctx context.Context, logger zerolog.Logger, localPath string, org string, name string, version string, typ PluginType, dops DownloaderOptions) (AssetSource, error) {
	if _, err := os.Stat(localPath); err == nil {
		return AssetSourceCached, nil
	}
	return AssetSourceRemote, doDownloadPluginFromGithub(ctx, logger, localPath, org, name, version, typ, dops)
}

func doDownloadPluginFromGithub(ctx context.Context, logger zerolog.Logger, localPath string, org string, name string, version string, typ PluginType, dops DownloaderOptions) error {
	downloadDir := filepath.Dir(localPath)
	pluginZipPath := localPath + ".zip"

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	downloadURL, err := getURLLocation(ctx, org, name, version, typ)
	if err != nil {
		return fmt.Errorf("failed to get plugin url: %w", err)
	}
	logger.Debug().Msg(fmt.Sprintf("Downloading %s", redactURLQuery(downloadURL)))
	defer os.Remove(pluginZipPath)

	if _, err := downloadFile(ctx, pluginZipPath, downloadURL, dops); err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	var pathInArchive string
	switch {
	case strings.HasPrefix(downloadURL, "https://github.com/cloudquery/cloudquery/releases/download/plugins-plugin"):
		pathInArchive = fmt.Sprintf("plugins/plugin/%s", name)
	case strings.HasPrefix(downloadURL, "https://github.com/cloudquery/cloudquery/releases/download/plugins-source"):
		pathInArchive = fmt.Sprintf("plugins/source/%s", name)
	case strings.HasPrefix(downloadURL, "https://github.com/cloudquery/cloudquery/releases/download/plugins-destination"):
		pathInArchive = fmt.Sprintf("plugins/destination/%s", name)
	case strings.HasPrefix(downloadURL, fmt.Sprintf("https://github.com/%s/cq-plugin", org)):
		pathInArchive = fmt.Sprintf("cq-plugin-%s", name)
	case strings.HasPrefix(downloadURL, fmt.Sprintf("https://github.com/%s/cq-source", org)):
		pathInArchive = fmt.Sprintf("cq-source-%s", name)
	case strings.HasPrefix(downloadURL, fmt.Sprintf("https://github.com/%s/cq-destination", org)):
		pathInArchive = fmt.Sprintf("cq-destination-%s", name)
	default:
		return fmt.Errorf("unknown GitHub %s", downloadURL)
	}

	return extractPluginBinary(pluginZipPath, WithBinarySuffix(pathInArchive), localPath)
}

func downloadFile(ctx context.Context, localPath string, downloadURL string, dops DownloaderOptions) (string, error) {
	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	urlForLog := redactURLQuery(downloadURL)

	checksum := ""
	options := []retry.Option{
		retry.RetryIf(isRetryableDownloadError),
		retry.Context(ctx),
		retry.Attempts(RetryAttempts),
		retry.Delay(RetryWaitTime),
		retry.MaxDelay(MaxRetryWaitTime),
	}
	retrier := retry.New(options...)
	err = retrier.Do(func() error {
		checksum = ""
		// Each attempt rewrites the file from the start, so a body that was cut off
		// mid-copy cannot leave its bytes in front of the next attempt's download.
		if err := truncateFile(out); err != nil {
			return err
		}

		// Get the data
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
		if err != nil {
			return fmt.Errorf("failed create request %s: %w", urlForLog, err)
		}

		// Do http request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get url %s: %w", urlForLog, err)
		}
		defer resp.Body.Close()
		// Check server response
		if resp.StatusCode == http.StatusNotFound {
			return errNotFound
		} else if resp.StatusCode != http.StatusOK {
			if isRetryableStatusCode(resp.StatusCode) {
				fmt.Printf("Failed downloading %s with status code %d. Retrying\n", urlForLog, resp.StatusCode)
			} else {
				fmt.Printf("Failed downloading %s with status code %d\n", urlForLog, resp.StatusCode)
			}
			return &httpStatusError{statusCode: resp.StatusCode}
		}

		fmt.Printf("Downloading %s\n", urlForLog)

		s := sha256.New()
		writers := []io.Writer{out, s}

		if !dops.NoProgress {
			bar := downloadProgressBar(resp.ContentLength, "Downloading")
			writers = append(writers, bar)
		}

		// Write the body to file
		written, err := io.Copy(io.MultiWriter(writers...), resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy body to file %s: %w", out.Name(), err)
		}
		if resp.ContentLength >= 0 && written != resp.ContentLength {
			return fmt.Errorf("%w: %s got %d bytes, want %d", errShortRead, out.Name(), written, resp.ContentLength)
		}
		checksum = fmt.Sprintf("%x", s.Sum(nil))
		return nil
	})
	if err != nil {
		if errors.Is(err, errNotFound) {
			return "", errNotFound
		}
		return "", fmt.Errorf("failed downloading URL %q. Error %w", urlForLog, err)
	}
	return checksum, nil
}

func truncateFile(f *os.File) error {
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file %s: %w", f.Name(), err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to rewind file %s: %w", f.Name(), err)
	}
	return nil
}

func downloadProgressBar(maxBytes int64, description ...string) *progressbar.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	return progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
}

func WithBinarySuffix(filePath string) string {
	if runtime.GOOS == "windows" {
		return filePath + ".exe"
	}
	return filePath
}
