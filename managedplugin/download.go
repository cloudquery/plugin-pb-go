package managedplugin

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"
)

const (
	DefaultDownloadDir = ".cq"
	RetryAttempts      = 5
	RetryWaitTime      = 1 * time.Second
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

	for _, downloadURL := range urls {
		err := retry.Do(func() error {
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
		}, retry.RetryIf(func(err error) bool {
			return err == err401 || err == err429
		}),
			retry.Context(ctx),
			retry.Attempts(RetryAttempts),
			retry.Delay(RetryWaitTime),
			retry.LastErrorOnly(true),
		)
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

func DownloadPluginFromHub(ctx context.Context, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions, dops DownloaderOptions) (DownloadSource, error) {
	if _, err := os.Stat(ops.LocalPath); err == nil {
		return DownloadSourceCached, nil
	}
	return DownloadSourceRemote, doDownloadPluginFromHub(ctx, c, ops, dops)
}

func doDownloadPluginFromHub(ctx context.Context, c *cloudquery_api.ClientWithResponses, ops HubDownloadOptions, dops DownloaderOptions) error {
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
		return fmt.Errorf("unauthorized. Try logging in via `cloudquery login`")
	case http.StatusNotFound:
		return fmt.Errorf("failed to download plugin %v %v/%v@%v: plugin version not found. If you're trying to use a private plugin you'll need to run `cloudquery login` first", ops.PluginKind, ops.PluginTeam, ops.PluginName, ops.PluginVersion)
	case http.StatusTooManyRequests:
		return fmt.Errorf("too many download requests. Try logging in via `cloudquery login` to increase rate limits")
	default:
		return fmt.Errorf("failed to download plugin %v %v/%v@%v: unexpected status code %v", ops.PluginKind, ops.PluginTeam, ops.PluginName, ops.PluginVersion, statusCode)
	}
	if pluginAsset == nil {
		return fmt.Errorf("failed to get plugin metadata from hub for %v %v/%v@%v: missing json response", ops.PluginKind, ops.PluginTeam, ops.PluginName, ops.PluginVersion)
	}
	location := pluginAsset.Location
	if len(location) == 0 {
		return fmt.Errorf("failed to get plugin metadata from hub: empty location from response")
	}
	pluginZipPath := ops.LocalPath + ".zip"
	writtenChecksum, err := downloadFile(ctx, pluginZipPath, location, dops)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	if pluginAsset.Checksum == "" {
		fmt.Printf("Warning - checksum not verified: %s\n", writtenChecksum)
	} else if writtenChecksum != pluginAsset.Checksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", pluginAsset.Checksum, writtenChecksum)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

	fileInArchive, err := archive.Open(fmt.Sprintf("plugin-%s-%s-%s-%s", ops.PluginName, ops.PluginVersion, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}

	out, err := os.OpenFile(ops.LocalPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", ops.LocalPath, err)
	}
	_, err = io.Copy(out, fileInArchive)
	if err != nil {
		return fmt.Errorf("failed to copy body to file: %w", err)
	}
	err = out.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
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

func DownloadPluginFromGithub(ctx context.Context, logger zerolog.Logger, localPath string, org string, name string, version string, typ PluginType, dops DownloaderOptions) (DownloadSource, error) {
	if _, err := os.Stat(localPath); err == nil {
		return DownloadSourceCached, nil
	}
	return DownloadSourceRemote, doDownloadPluginFromGithub(ctx, logger, localPath, org, name, version, typ, dops)
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
	logger.Debug().Msg(fmt.Sprintf("Downloading %s", downloadURL))
	if _, err := downloadFile(ctx, pluginZipPath, downloadURL, dops); err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

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

	pathInArchive = WithBinarySuffix(pathInArchive)
	fileInArchive, err := archive.Open(pathInArchive)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive plugins/source/%s: %w", name, err)
	}
	out, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	_, err = io.Copy(out, fileInArchive)
	if err != nil {
		return fmt.Errorf("failed to copy body to file: %w", err)
	}
	err = out.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

func downloadFile(ctx context.Context, localPath string, downloadURL string, dops DownloaderOptions) (string, error) {
	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	checksum := ""
	err = retry.Do(func() error {
		checksum = ""
		// Get the data
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
		if err != nil {
			return fmt.Errorf("failed create request %s: %w", downloadURL, err)
		}

		// Do http request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get url %s: %w", downloadURL, err)
		}
		defer resp.Body.Close()
		// Check server response
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("not found")
		} else if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed downloading %s with status code %d. Retrying\n", downloadURL, resp.StatusCode)
			return errors.New("statusCode != 200")
		}

		urlForLog := downloadURL
		parsedURL, err := url.Parse(downloadURL)
		if err == nil {
			parsedURL.RawQuery = ""
			parsedURL.Fragment = ""
			urlForLog = parsedURL.String()
		}
		fmt.Printf("Downloading %s\n", urlForLog)

		s := sha256.New()
		writers := []io.Writer{out, s}

		if !dops.NoProgress {
			bar := downloadProgressBar(resp.ContentLength, "Downloading")
			writers = append(writers, bar)
		}

		// Write the body to file
		_, err = io.Copy(io.MultiWriter(writers...), resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy body to file %s: %w", out.Name(), err)
		}
		checksum = fmt.Sprintf("%x", s.Sum(nil))
		return nil
	}, retry.RetryIf(func(err error) bool {
		return err.Error() == "statusCode != 200"
	}),
		retry.Context(ctx),
		retry.Attempts(RetryAttempts),
		retry.Delay(RetryWaitTime),
	)
	if err != nil {
		for _, e := range err.(retry.Error) {
			if e.Error() == "not found" {
				return "", e
			}
		}
		return "", fmt.Errorf("failed downloading URL %q. Error %w", downloadURL, err)
	}
	return checksum, nil
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
