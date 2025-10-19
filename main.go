package main

import (
	"fmt"
	"ikki/task"
	"ikki/worker"
	"time"
)

func main() {
	t := &task.Task{
		Name:  "test-container-1",
		Image: "postgres:latest",
		Env: []string{
			"POSTGRES_USER=me",
			"POSTGRES_PASSWORD=you",
		},
	}

	w, err := worker.New()
	if err != nil {
		panic(err)
	}
	result := w.StartTask(t)
	fmt.Printf("Test container started successfully. (container_id: %s)\n", result.ContainerID)

	fmt.Println("Waiting for 10 seconds")
	time.Sleep(time.Second * 10)

	fmt.Println("Getting rid of the task")
	result = w.StopTask(t)

	if result.Error != nil {
		panic(result.Error)
	}

}
