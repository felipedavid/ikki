package main

import (
	"fmt"
	"ikki/task"
)

func main() {
	t, err := task.New()
	if err != nil {
		panic(err)
	}
	t.Config.Image = "postgres"

	result := t.Run()
	fmt.Printf("Result = %v\n", result)
}