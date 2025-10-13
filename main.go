package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"ikki/task"
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

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	result := t.Run(cli)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("Test container started successfully. (container_id: %s)", result.ContainerID)
}
