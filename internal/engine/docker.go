package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type DockerEngine struct {
	client *client.Client
}

func NewDockerEngine(
	dockerClient *client.Client,
) *DockerEngine {
	return &DockerEngine{
		client: dockerClient,
	}
}

func (e *DockerEngine) Find(ctx context.Context, serviceName string) (*Container, error) {
	listOpts := container.ListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("fickle.service.name=%s", serviceName),
		}),
	}

	containers, err := e.client.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	if len(containers) > 0 {
		cont := containers[0]

		status := ContainerStatusUnknown
		if cont.State == container.StateRunning {
			status = ContainerStatusRunning
		} else if cont.State == container.StateRestarting {
			status = ContainerStatusRestarting
		}

		return &Container{
			ID:     cont.ID,
			Status: status,
		}, nil
	}

	return nil, errors.New("container not found")
}

func (e *DockerEngine) Stop(ctx context.Context, id string) error {
	return e.client.ContainerStop(ctx, id, container.StopOptions{})
}

func (e *DockerEngine) Start(ctx context.Context, id string) error {
	return e.client.ContainerStart(ctx, id, container.StartOptions{})
}
