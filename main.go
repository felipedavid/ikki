package main

import (
	"fmt"
	"ikki/task"
)

func main() {
	taskConfig := &task.Config{
		Name: "test-container-1",
		Image: "postgres:latest",
		Env: []string{
			"POSTGRES_USER=me",
			"POSTGRES_PASSWORD=you",
		},
	}

	task, err := task.New(taskConfig)
	if err != nil {
		panic(err)
	}

	result := task.Run()
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("Test container started successfully. (container_id: %s)", result.ContainerID)
}