package managedplugin

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/schollz/progressbar/v3"
)

const (
	DefaultDownloadDir = ".cq"
	RetryAttempts      = 5
	RetryWaitTime      = 1 * time.Second
)

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

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
		if err != nil {
			return "", fmt.Errorf("failed create request %s: %w", url, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to get url %s: %w", url, err)
		}
		// Check server response
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			continue
		} else if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Printf("Failed downloading %s with status code %d. Retrying\n", url, resp.StatusCode)
			return "", errors.New("statusCode != 200")
		}
		resp.Body.Close()
		return url, nil
	}

	return "", fmt.Errorf("failed to find plugin %s/%s version %s", org, name, version)
}

func DownloadPluginFromGithub(ctx context.Context, localPath string, org string, name string, version string, typ PluginType) error {
	downloadDir := filepath.Dir(localPath)
	pluginZipPath := localPath + ".zip"

	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	url, err := getURLLocation(ctx, org, name, version, typ)
	if err != nil {
		return fmt.Errorf("failed to get plugin url: %w", err)
	}
	if err := downloadFile(ctx, pluginZipPath, url); err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

	var pathInArchive string
	switch {
	case strings.HasPrefix(url, "https://github.com/cloudquery/cloudquery/releases/download/plugins-plugin"):
		pathInArchive = fmt.Sprintf("plugins/plugin/%s", name)
	case strings.HasPrefix(url, "https://github.com/cloudquery/cloudquery/releases/download/plugins-source"):
		pathInArchive = fmt.Sprintf("plugins/source/%s", name)
	case strings.HasPrefix(url, "https://github.com/cloudquery/cloudquery/releases/download/plugins-destination"):
		pathInArchive = fmt.Sprintf("plugins/destination/%s", name)
	case strings.HasPrefix(url, fmt.Sprintf("https://github.com/%s/cq-plugin", org)):
		pathInArchive = fmt.Sprintf("cq-plugin-%s", name)
	case strings.HasPrefix(url, fmt.Sprintf("https://github.com/%s/cq-source", org)):
		pathInArchive = fmt.Sprintf("cq-source-%s", name)
	case strings.HasPrefix(url, fmt.Sprintf("https://github.com/%s/cq-destination", org)):
		pathInArchive = fmt.Sprintf("cq-destination-%s", name)
	default:
		return fmt.Errorf("unknown GitHub %s", url)
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

func downloadFile(ctx context.Context, localPath string, url string) error {
	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	err = downloadFileFromURL(ctx, out, url)
	if err != nil {
		return err
	}
	return nil
}

func downloadFileFromURL(ctx context.Context, out *os.File, url string) error {
	err := retry.Do(func() error {
		// Get the data
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("failed create request %s: %w", url, err)
		}

		// Do http request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get url %s: %w", url, err)
		}
		defer resp.Body.Close()
		// Check server response
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("not found")
		} else if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed downloading %s with status code %d. Retrying\n", url, resp.StatusCode)
			return errors.New("statusCode != 200")
		}

		fmt.Printf("Downloading %s\n", url)
		bar := downloadProgressBar(resp.ContentLength, "Downloading")

		// Writer the body to file
		_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy body to file %s: %w", out.Name(), err)
		}
		return nil
	}, retry.RetryIf(func(err error) bool {
		return err.Error() == "statusCode != 200"
	}),
		retry.Attempts(RetryAttempts),
		retry.Delay(RetryWaitTime),
	)
	if err != nil {
		for _, e := range err.(retry.Error) {
			if e.Error() == "not found" {
				return e
			}
		}
		return fmt.Errorf("failed downloading URL %q. Error %w", url, err)
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

type dockerProgressReader struct {
	decoder        *json.Decoder
	bar            *progressbar.ProgressBar
	downloadedByID map[string]int64
	totalBytes     int64
}

type dockerProgressInfo struct {
	Status       string `json:"status"`
	Progress     string `json:"progress"`
	ProgressData struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	ID string `json:"id"`
}

func (pr *dockerProgressReader) Read(p []byte) (n int, err error) {
	var progress dockerProgressInfo
	err = pr.decoder.Decode(&progress)
	if err != nil {
		if err == io.EOF {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("failed to decode JSON: %v", err)
	}
	if progress.Status == "Downloading" {
		if pr.bar == nil {
			pr.bar = downloadProgressBar(1, "Downloading")
			pr.bar.RenderBlank()
		}
		if _, seen := pr.downloadedByID[progress.ID]; !seen {
			pr.downloadedByID[progress.ID] = 0
			pr.totalBytes += progress.ProgressData.Total
			pr.bar.ChangeMax64(pr.totalBytes)
		}
		pr.downloadedByID[progress.ID] = progress.ProgressData.Current
		total := int64(0)
		for _, v := range pr.downloadedByID {
			total += v
		}
		if total < pr.totalBytes {
			// progressbar stops responding if it reaches 100%, so as a workaround we don't update
			// the bar if we're at 100%, because there may be more layers of the image
			// coming that we don't know about.
			pr.bar.Set64(total)
		}
	}
	return
}
