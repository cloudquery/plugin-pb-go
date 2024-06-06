package managedplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/distribution/reference"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/schollz/progressbar/v3"
)

var (
	ErrLoginRequired = fmt.Errorf("login required")
	ErrTeamRequired  = fmt.Errorf("team required")
)

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

func (pr *dockerProgressReader) Read(_ []byte) (n int, err error) {
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
			_ = pr.bar.RenderBlank()
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
			_ = pr.bar.Set64(total)
		}
	}
	return 0, nil
}

func newDockerClient(options ...client.Opt) (*client.Client, error) {
	defaultOptions := []client.Opt{
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	}
	options = append(defaultOptions, options...)
	return client.NewClientWithOpts(options...)
}

func isDockerImageAvailable(ctx context.Context, imageName string) (bool, error) {
	cli, err := newDockerClient()
	if err != nil {
		return false, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// List the images matching the given name
	images, err := cli.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return false, fmt.Errorf("failed to list Docker images: %w", err)
	}

	// Check if the image with the specified name was found
	return len(images) > 0, nil
}

func pullDockerImage(ctx context.Context, imageName string, authToken string, teamName string, dockerHubAuth string, dops DownloaderOptions) error {
	// Pull the image
	additionalHeaders := make(map[string]string)
	opts := image.PullOptions{}
	if strings.HasPrefix(imageName, "docker.cloudquery.io") {
		if authToken == "" {
			return ErrLoginRequired
		}
		if teamName == "" {
			return ErrTeamRequired
		}
		namedRef, err := reference.ParseNormalizedNamed(imageName)
		if err != nil {
			return fmt.Errorf("failed to parse Docker image tag: %v", err)
		}
		nameWithTag, ok := namedRef.(reference.NamedTagged)
		if !ok {
			return fmt.Errorf("failed to parse Docker image tag: %v", err)
		}
		additionalHeaders["X-Meta-Plugin-Version"] = nameWithTag.Tag()
		additionalHeaders["X-Meta-User-Team-Name"] = teamName
		authConfig := registry.AuthConfig{
			Username: "managedplugin",
			Password: authToken,
		}
		encodedAuth, err := registry.EncodeAuthConfig(authConfig)
		if err != nil {
			return fmt.Errorf("failed to encode Docker auth config: %v", err)
		}
		opts.RegistryAuth = encodedAuth
	} else if dockerHubAuth != "" {
		opts.RegistryAuth = dockerHubAuth
	}

	cli, err := newDockerClient(client.WithHTTPHeaders(additionalHeaders))
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	out, err := cli.ImagePull(ctx, imageName, opts)
	if err != nil {
		return fmt.Errorf("failed to pull Docker image: %v", err)
	}
	defer out.Close()

	if dops.NoProgress {
		_, err = io.Copy(io.Discard, out)
		if err != nil {
			return fmt.Errorf("failed to copy image pull output: %v", err)
		}
		return nil
	}

	// Create a progress reader to display the download progress
	pr := &dockerProgressReader{
		decoder:        json.NewDecoder(out),
		downloadedByID: map[string]int64{},
		bar:            nil,
	}
	_, err = io.Copy(io.Discard, pr)
	if err != nil {
		return fmt.Errorf("failed to copy image pull output: %v", err)
	}
	if pr.bar != nil {
		_ = pr.bar.Finish()
		pr.bar.Close()
	}

	return nil
}
