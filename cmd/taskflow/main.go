package main

import (
	"context"
	"fmt"
	"log"
	"taskflow/internal/idgen"
	"taskflow/internal/repository/sqlite"
	"taskflow/internal/service"
	"taskflow/internal/task"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	repository, err := sqlite.NewRepository(ctx, "data/taskflow.db")
	if err != nil {
		return fmt.Errorf("initialize repository: %w", err)
	}
	defer func() {
		if err := repository.Close(); err != nil {
			log.Printf("close repository: %v", err)
		}
	}()

	generator := idgen.UUIDGenerator{}
	taskService := service.NewTaskService(
		generator,
		repository,
	)

	newTask, err := taskService.CreateTask(
		ctx,
		"Learn Go",
		"Continue developing TaskFlow",
		task.PriorityHigh,
	)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	fmt.Println("Task created successfully!")
	fmt.Printf("ID: %s\n", newTask.ID())
	fmt.Printf("Title: %s\n", newTask.Title())
	fmt.Printf("Description: %s\n", newTask.Description())
	fmt.Printf("Status: %s\n", newTask.Status())
	fmt.Printf("Priority: %s\n", newTask.Priority())
	return nil
}
