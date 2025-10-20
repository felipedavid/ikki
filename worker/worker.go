package worker

import (
	"github.com/docker/docker/client"
	"ikki/task"
	"ikki/utils"
	"log/slog"

	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     utils.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
	dClient   *client.Client
}

func New() (*Worker, error) {
	clt, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &Worker{
		dClient: clt,
	}, nil
}

func (w *Worker) AddTask(t *task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) StartTask(t *task.Task) task.TaskResult {
	result := t.Run(w.dClient)
	if result.Error != nil {
		slog.Error("Error running task", "task_id", t.ID, "error", result.Error)
		t.State = task.Failed
		return result
	}

	t.State = task.Running

	return result
}

func (w *Worker) StopTask(t *task.Task) task.TaskResult {
	result := t.Stop(w.dClient)
	if result.Error != nil {
		slog.Error("Unable to stop container", "container_id", t.ContainerID)
	}

	t.State = task.Completed

	return result
}
