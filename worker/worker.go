package worker

import (
	"fmt"
	"ikki/task"
	"ikki/utils"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     utils.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}

func (w *Worker) RunTask() {
	fmt.Println("I will start or stop the task")
}

func (w *Worker) StartTask() {
	fmt.Println("I will start a task")
}

func (w *Worker) StartTask(t *task.Task) task.DockerResult {
	d, err := task.NewDocker(nil)
	if err != nil {
		return task.DockerResult{Error: err}
	}

	t.Docker = d
	result := t.Run()

	return result
}

func (w *Worker) StopTask(t *task.Task) task.DockerResult {
	d, err := task.NewDocker(nil)
	if err != nil {
		return task.DockerResult{Error: err}
	}

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		slog.Error("Unable to stop container", "container_id", t.ContainerID)
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed

	return result
}
