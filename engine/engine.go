package engine

import (
	"context"
	"time"

	"github.com/bodaay/mosalamaagent/logging"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type EngineManager struct {
	dockerClient *client.Client
}

func NewEngineManager() (*EngineManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &EngineManager{
		dockerClient: cli,
	}, nil
}

type ContainerResources struct {
	CPUQuota int64
	Memory   int64
	// Add GPU resource fields if necessary
}

func (e *EngineManager) StartEngine(ctx context.Context, imageName string, containerName string, cmd []string, ports map[string]string, resources ContainerResources) error {
	// Pull the image
	reader, err := e.dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	// Optionally, read the output from reader to monitor progress

	// Configure resource constraints
	resourcesConfig := container.Resources{
		CPUQuota: resources.CPUQuota,
		Memory:   resources.Memory,
		// Add GPU resource configurations if necessary
	}

	// Configure container
	containerConfig := &container.Config{
		Image: imageName,
		Cmd:   cmd,
	}

	hostConfig := &container.HostConfig{
		PortBindings: natPortBindings(ports),
		Resources:    resourcesConfig,
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Binds: []string{
			"/var/mosalamaagent/models:/models", // Adjust the paths as necessary
		},
	}

	resp, err := e.dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	// Start the container
	if err := e.dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	logging.Log.Infof("Engine started with container ID: %s", resp.ID)
	return nil
}

func (e *EngineManager) StopEngine(ctx context.Context, containerName string) error {
	timeout := time.Second * 10
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}
	if err := e.dockerClient.ContainerStop(ctx, containerName, stopOptions); err != nil {
		return err
	}
	logging.Log.Infof("Engine stopped: %s", containerName)
	return nil
}
func (e *EngineManager) ListEngines(ctx context.Context) ([]types.Container, error) {
	containers, err := e.dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// Helper function to convert port mappings
func natPortBindings(ports map[string]string) nat.PortMap {
	portMap := nat.PortMap{}
	for containerPort, hostPort := range ports {
		portMap[nat.Port(containerPort)] = []nat.PortBinding{
			{
				HostPort: hostPort,
			},
		}
	}
	return portMap
}
