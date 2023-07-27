package managedplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func isDockerImageAvailable(ctx context.Context, imageName string) (bool, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// List the images matching the given name
	images, err := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return false, fmt.Errorf("failed to list Docker images: %w", err)
	}

	// Check if the image with the specified name was found
	return len(images) > 0, nil
}

func pullDockerImage(ctx context.Context, imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Pull the image
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull Docker image: %v", err)
	}
	defer out.Close()

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
		pr.bar.Finish()
		pr.bar.Close()
	}

	return nil
}
