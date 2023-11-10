package managedplugin

import (
	"archive/zip"
	"context"
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

	for _, downloadURL := range urls {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, downloadURL, nil)
		if err != nil {
			return "", fmt.Errorf("failed create request %s: %w", downloadURL, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to get url %s: %w", downloadURL, err)
		}
		// Check server response
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			continue
		} else if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Printf("Failed downloading %s with status code %d. Retrying\n", downloadURL, resp.StatusCode)
			return "", errors.New("statusCode != 200")
		}
		resp.Body.Close()
		return downloadURL, nil
	}

	return "", fmt.Errorf("failed to find plugin %s/%s version %s", org, name, version)
}

func DownloadPluginFromHub(ctx context.Context, authToken, localPath, team, name, version string, typ PluginType) error {
	downloadDir := filepath.Dir(localPath)
	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	target := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	// We don't want to follow redirects because we want to get the download URL and show progress bar while downloading
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	c, err := cloudquery_api.NewClient(APIBaseURL(), cloudquery_api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if authToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		}
		return nil
	}))
	c.Client = client
	if err != nil {
		return fmt.Errorf("failed to create Hub API client: %w", err)
	}

	downloadURL, err := c.DownloadPluginAsset(ctx, team, cloudquery_api.PluginKind(typ.String()), name, version, target)
	if err != nil {
		return fmt.Errorf("failed to get plugin url: %w", err)
	}
	defer downloadURL.Body.Close()
	switch downloadURL.StatusCode {
	case http.StatusOK, http.StatusNoContent, http.StatusFound:
		// we allow these status codes, but typically expect a redirect (302)
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized. Try logging in via `cloudquery login`")
	case http.StatusNotFound:
		return fmt.Errorf("failed to download plugin %v %v/%v@%v: plugin version not found. If you're trying to use a private plugin you'll need to run `cloudquery login` first", typ, team, name, version)
	case http.StatusTooManyRequests:
		return fmt.Errorf("too many download requests. Try logging in via `cloudquery login` to increase rate limits")
	default:
		return fmt.Errorf("failed to download plugin %v %v/%v@%v: unexpected status code %v", typ, team, name, version, downloadURL.StatusCode)
	}
	location, ok := downloadURL.Header["Location"]
	if !ok {
		return fmt.Errorf("failed to get plugin url for %v %v/%v@%v: missing location header from response", typ, team, name, version)
	}
	if len(location) == 0 {
		return fmt.Errorf("failed to get plugin url: empty location header from response")
	}
	pluginZipPath := localPath + ".zip"
	err = downloadFile(ctx, pluginZipPath, location[0])
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	archive, err := zip.OpenReader(pluginZipPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer archive.Close()

	fileInArchive, err := archive.Open(fmt.Sprintf("plugin-%s-%s-%s-%s", name, version, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
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

func DownloadPluginFromGithub(ctx context.Context, localPath string, org string, name string, version string, typ PluginType) error {
	downloadDir := filepath.Dir(localPath)
	pluginZipPath := localPath + ".zip"

	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory %s: %w", downloadDir, err)
	}

	downloadURL, err := getURLLocation(ctx, org, name, version, typ)
	if err != nil {
		return fmt.Errorf("failed to get plugin url: %w", err)
	}
	if err := downloadFile(ctx, pluginZipPath, downloadURL); err != nil {
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

func downloadFile(ctx context.Context, localPath string, downloadURL string) error {
	// Create the file
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer out.Close()

	err = downloadFileFromURL(ctx, out, downloadURL)
	if err != nil {
		return err
	}
	return nil
}

func downloadFileFromURL(ctx context.Context, out *os.File, downloadURL string) error {
	err := retry.Do(func() error {
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
		return fmt.Errorf("failed downloading URL %q. Error %w", downloadURL, err)
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
