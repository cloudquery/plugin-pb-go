package managedplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func isDockerImageAvailable(ctx context.Context, imageName string) (bool, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, fmt.Errorf("failed to create Docker client: %w", err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return false, fmt.Errorf("failed to list Docker images: %w", err)
	}

	return len(images) > 0, nil
}

func pullDockerImage(ctx context.Context, imageName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull Docker image: %v", err)
	}
	defer out.Close()

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
func getContainerConnectionString(ctx context.Context, containerName string, port string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	// Get a list of running containers, filtering by given name
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("status", "running"),
			filters.Arg("name", fmt.Sprintf("^/%s$", containerName)),
		),
	})
	if err != nil {
		return "", err
	}

	// Find the container with the specified name
	var targetContainer types.Container
	for _, c := range containers {
		for _, name := range c.Names {
			if name == containerName {
				targetContainer = c
				break
			}
		}
	}

	// Check if the container with the specified name was found
	if targetContainer.ID == "" {
		return "", fmt.Errorf("container with name '%s' not found", containerName)
	}

	// Get container details including the published ports
	containerDetails, err := cli.ContainerInspect(ctx, targetContainer.ID)
	if err != nil {
		return "", err
	}

	// Extract the published ports
	ports := containerDetails.NetworkSettings.Ports
	for _, binding := range ports[nat.Port(port)] {
		return binding.HostIP + ":" + binding.HostPort, nil
	}
	return "", fmt.Errorf("port %s not found", port)
}
