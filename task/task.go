package task

import (
	"context"
	"log/slog"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID uuid.UUID
	Name string
	State State
	Image string
	Memory int
	Disk int
	ExposedPorts nat.PortSet
	PortBindings map[string]string
	RestartPolicy string
	StartTime time.Time
	FinishTime time.Time
	Docker
}

func New() (*Task, error) {
	t := &Task{}
	client, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	t.Docker.Client = client


	return t, nil
}

type TaskEvent struct {
	ID uuid.UUID
	State State
	Timestamp time.Time
	Task Task
}

type Config struct {
	Name string
	AttachStdin bool
	AttachStdout bool
	AttachStderr bool
	ExposedPorts nat.PortSet
	Cmd []string
	Image string
	Cpu float64
	Memory int64
	Disk int64
	Env []string
	RestartPolicy string
}

type Docker struct {
	Client *client.Client
	Config Config
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	_, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		slog.Error("Error pulling image", "err", err, "name", d.Config.Image)
		return DockerResult{
			Error: err,
		}
	}

	restartPolicy := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}

	resources := container.Resources{
		Memory: d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	config := container.Config{
		Image: d.Config.Image,
		Tty: false,
		Env: d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hostConfig := container.HostConfig{
		RestartPolicy: restartPolicy,
		Resources: resources,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &config, &hostConfig, nil, nil, d.Config.Name)
	if err != nil {
		slog.Error("Error creating the container", "err", err, "image", d.Config.Name)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return DockerResult{Error: err}
	}

	//d.Config.Runtime.ContainerID = resp.ID

	containerLogs, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		slog.Error("Unable to get container logs", "err", err, "container_id", resp.ID)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, containerLogs)

	return DockerResult{
		ContainerID: resp.ID,
		Action: "start",
		Result: "success",
	}
}

type DockerResult struct {
	Error error
	Action string
	ContainerID string
	Result string
}