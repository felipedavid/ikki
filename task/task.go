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
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Memory        int64
	Cpu           float64
	Disk          int
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
	Env           []string
}

type TaskResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

func (t *Task) Run(dc *client.Client) TaskResult {
	ctx := context.Background()
	_, err := dc.ImagePull(ctx, t.Image, image.PullOptions{})
	if err != nil {
		slog.Error("Error pulling image", "err", err, "name", t.Image)
		return TaskResult{
			Error: err,
		}
	}

	restartPolicy := container.RestartPolicy{
		Name: container.RestartPolicyMode(t.RestartPolicy),
	}

	resources := container.Resources{
		Memory:   t.Memory,
		NanoCPUs: int64(t.Cpu * math.Pow(10, 9)),
	}

	config := container.Config{
		Image:        t.Image,
		Tty:          false,
		Env:          t.Env,
		ExposedPorts: t.ExposedPorts,
	}

	hostConfig := container.HostConfig{
		RestartPolicy:   restartPolicy,
		Resources:       resources,
		PublishAllPorts: true,
	}

	resp, err := dc.ContainerCreate(ctx, &config, &hostConfig, nil, nil, t.Name)
	if err != nil {
		slog.Error("Error creating the container", "err", err, "image", t.Image)
		return TaskResult{Error: err}
	}

	err = dc.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return TaskResult{Error: err}
	}

	t.ContainerID = resp.ID

	containerLogs, err := dc.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		slog.Error("Unable to get container logs", "err", err, "container_id", resp.ID)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, containerLogs)

	return TaskResult{
		ContainerID: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (t *Task) Stop(cl *client.Client) TaskResult {
	ctx := context.Background()
	stopOptions := container.StopOptions{}
	err := cl.ContainerStop(ctx, t.ContainerID, stopOptions)
	if err != nil {
		slog.Error("Unable to stop the container", "container_id", t.ContainerID, "err", err)
		return TaskResult{
			ContainerID: t.ContainerID,
			Action:      "stop",
			Error:       err,
		}
	}

	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}
	err = cl.ContainerRemove(ctx, t.ContainerID, removeOptions)
	if err != nil {
		slog.Error("Unable to remove the container", "container_id", t.ContainerID, "err", err)
		return TaskResult{
			ContainerID: t.ContainerID,
			Action:      "stop",
			Error:       err,
		}
	}

	return TaskResult{
		ContainerID: t.ContainerID,
		Action:      "stop",
		Result:      "success",
	}
}
