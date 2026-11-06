package main

import (
	"fmt"
	"log"
	"taskflow/internal/task"
)

func main() {
	newTask, err := task.NewTask(
		"Learn Go",
		"Continue developing TaskFlow",
		task.PriorityHigh,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task created succesfully!")
	fmt.Printf("Status: %s\n", newTask.Status())
	fmt.Println(newTask.Title())
	fmt.Println(newTask.Description())
	fmt.Println(newTask.Priority())
}
